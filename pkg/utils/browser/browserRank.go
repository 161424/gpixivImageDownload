package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/dao/sql"
	"gpixivImageDownload/internal/addr"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Br struct {
	Client *http.Client
}

var Brs = &Br{}
var l *log.Logs
var refere = "https://www.pixiv.net"
var mutilhttps = []string{}

// rank image
type RpJs struct {
	Mode      string    `json:"mode"`
	Content   string    `json:"content"`
	Contents  []content `json:"contents"`
	CurrPage  any       `json:"page"`
	NextPage  any       `json:"next"`
	PrevPage  any       `json:"prev"`
	CurrDate  string    `json:"date"`
	NextDate  any       `json:"next_date"`
	PrevDate  any       `json:"prev_date"`
	RankTotal int       `json:"rank_total"`
	Url       string    `json:"url"`
}

// 对于跨包，子特征名大小写不重要，但是变量名一定要大写才能被解析到

type content struct {
	Title             string         `json:"title"`
	Tags              []string       `json:"tags"`
	IllustType        string         `json:"illust_type"`
	IllustBookStyle   string         `json:"illust_book_style"`
	IllustPageCount   string         `json:"illust_page_count"`
	UserName          string         `json:"user_name"`
	IllustContentType map[string]any `json:"illust_content_type"`
	IllustId          float64        `json:"illust_id"`
	UserId            float64        `json:"user_id"`
}

type Cont struct {
	Timestamp string                            `json:"timestamp"`
	Illust    map[string]Illust                 `json:"illust"`
	User      map[string]map[string]interface{} `json:"user"`
}

type Illust struct {
	UserId        string                 `json:"userId"`
	UserName      string                 `json:"userName"`
	PageCount     float64                `json:"pageCount"`
	Urls          map[string]string      `json:"urls"`
	IllustTitle   string                 `json:"illustTitle"`
	SeriesNavData map[string]interface{} `json:"seriesNavData"`
	ViewCount     float64                `json:"viewCount"`
	LikeCount     float64                `json:"likeCount"`
	Tags          Tags                   `json:"tags"`
	CreateDate    string                 `json:"createDate"`
	Width         float64                `json:"width"`
	Height        float64                `json:"height"`
	BookmarkCount float64                `json:"bookmarkCount"`
	ResponseCount float64                `json:"responseCount"`
	AiType        float64                `json:"aiType"`
}

type Tags struct {
	AuthorId string                   `json:"authorId"`
	Tags     []map[string]interface{} `json:"tags"`
}

func init() {
	l = log.NewSlogGroup("Browser")
	sockets, _ := url.Parse("http://127.0.0.1:10809")
	l.Send(slog.LevelInfo, fmt.Sprintf("Proxy:%s", sockets), log.LogFiles|log.LogStdouts)
	Brs.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          50,
			IdleConnTimeout:       60 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			Proxy:                 http.ProxyURL(sockets),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

}

func SetMutliHttps(cks []string) {
	mutilhttps = cks
}

func GetPixivPage(urls string, idx int) ([]byte, error) {
	client := Brs.Client
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", urls, nil)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("url的req创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return nil, err
	}
	//req.Header.Set("if-none-match", "lns0njwq6e5otu")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", refere)
	req.Header.Set("User-Agent", addr.Header.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	for _, i := range strings.Split(mutilhttps[idx], "; ") {
		a := strings.Split(i, "=")
		req.AddCookie(&http.Cookie{Name: a[0], Value: a[1]})
	}

	resp, err := client.Do(req)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("url 请求信息失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return nil, err
	}

	bts, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bts, nil
}

func (br *Br) fixUrl(url string, useHttps bool) string {
	if !strings.HasPrefix(url, "http") {
		if !strings.HasPrefix(url, "/") {
			url = "/" + url
		}
		if useHttps {
			return "https://www.pixiv.net" + url
		} else {
			return "http://www.pixiv.net" + url
		}
	}
	return url
}

func CheckResp(resp []byte) (s string) {
	if strings.Contains(string(resp), "Just a moment...") {
		s = "受到了CloudFlare 5s盾的管控， 请检查cookie和梯子是否有效"
		return
	}

	if strings.Contains(string(resp), "不在排行榜统计范围内") {
		s = "不在排行榜统计范围内"
		return
	}

	return
}

func (br *Br) GetPixivPage(url string, idx int) (RJ *RpJs, s string) {
	l.Send(slog.LevelInfo, url, log.LogFiles|log.LogStdouts)
	rb, err := GetPixivPage(url, idx)
	if err != nil {
		l.Send(slog.LevelError, err.Error(), log.LogFiles|log.LogStdouts)
		return
	}

	s = CheckResp(rb)
	if err = json.Unmarshal(rb, &RJ); err != nil {
		s = "解析错误，请检查参数是否匹配"
		return
	}

	return
}

func (br *Br) GetPixivRanking(mode, date string, comm *model.Common, page int) (df *RpJs, error string) {
	url := fmt.Sprintf("https://www.pixiv.net/ranking.php?mode=%s", mode)
	// https://www.pixiv.net/ranking.php?mode=daily
	// https://www.pixiv.net/ranking.php?mode=weekly
	// https://www.pixiv.net/ranking.php?mode=monthly
	// https://www.pixiv.net/ranking.php?mode=rookie  新人
	// 原创
	// ai
	// 受男
	// 受女
	if len(date) > 0 {
		if comm.R18 {
			url = fmt.Sprintf("%s&date=%s_r18", url, date)
		} else {
			url = fmt.Sprintf("%s&date=%s", url, date)
		}

		//daily https://www.pixiv.net/ranking.php?mode=daily&date=20240512
		//week https://www.pixiv.net/ranking.php?mode=weekly&date=20240509  2024-5-3~2024-5-9
		//monthly https://www.pixiv.net/ranking.php?mode=monthly&date=20240506  2024-4-7~2024-5-6
		//rookie https://www.pixiv.net/ranking.php?mode=rookie&date=20240505  2024-5-3~2024-5-9
	}
	// content illust  只下载插画
	// https://www.pixiv.net/ranking.php?mode=daily&content=illust&date=20240509
	// content all  综合

	if !comm.SkipIllus {
		url = fmt.Sprintf("%s&content=%s", url, "illust")
	} else if !comm.SkipManga {
		url = fmt.Sprintf("%s&content=%s", url, "manga")
	} else if !comm.SkipUgoira {
		url = fmt.Sprintf("%s&content=%s", url, "ugoira")
	}

	url = fmt.Sprintf("%s&p=%d&format=json", url, page)
	for idx := range mutilhttps {
		df, error = br.GetPixivPage(url, idx)
		if error == "" {
			break
		}
	}

	return
}

func (br *Br) AnalysisUrl(url string) (*goquery.Document, error) {

	resp, err := GetPixivPage(url, 0)
	fmt.Println(string(resp))
	s := CheckResp(resp)
	if err != nil {
		return nil, err
	}
	if s != "" {
		return nil, errors.New(s)
	}

	df := bytes.NewReader(resp)
	doc, err := goquery.NewDocumentFromReader(df)
	if err != nil {
		return nil, err
	}
	br.CheckWeb(doc)
	return doc, nil
}

func (br *Br) CheckWeb(doc *goquery.Document) {
	if IsNotLoggedIn(doc) {
		panic("Not Logged In!")
	}
	if IsNeedPermission(doc) {
		panic("Not in MyPick List, Need Permission!")
	}
	if IsNeedAppropriateLevel(doc) {
		panic("Public works can not be viewed by the appropriate level!")
	}
	if IsDeleted(doc) {
		panic("Image not found/already deleted!")
	}
	if IsGuroDisabled(doc) {
		panic("Image is disabled for under 18, check your setting page (R-18/R-18G)!")
	}

}

func (br *Br) GetImageInfo(imageId string) (imageInfo *sql.ImageInfo, err error) {

	var doc *goquery.Document
	l.Send(slog.LevelDebug, fmt.Sprintf("Getting image page: %s", imageId), log.LogStdouts)
	// https://www.pixiv.net/artworks/115906527
	url := fmt.Sprintf("https://www.pixiv.net/artworks/%s", imageId)
	//fmt.Println("url", url)
	if doc, err = br.AnalysisUrl(url); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("url 解析失败: %s.%s", imageId, err), log.LogStdouts)
		return nil, err
	}

	imageInfo, err = ParseRankImage(doc, imageId)

	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("resp.body 解析json失败: %s", imageId), log.LogStdouts)
		return nil, err
	}
	if imageInfo.ImageMode == "ugoira_view" {
		ugoiraMetaUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/ugoira_meta", imageId)
		// ugoira 代表是动图，暂不进行过多考虑
		fmt.Println(ugoiraMetaUrl)
	}
	imageInfo.ImageId = imageId
	return
}

func (br *Br) DownloadImage(ctx context.Context, imgInfo *sql.ImageInfo, path string, thread_ bool) (state [2]int64) {
	stopG := false
	go func() {
		<-ctx.Done()
		stopG = true

	}()
	var wg = sync.WaitGroup{}
	var thread = make(chan [0]int)
	var _p string
	stda1 := atomic.Int64{}
	stda2 := atomic.Int64{}
	stda1.Store(0)
	stda2.Store(0)
	for index, imgUrl := range imgInfo.ImageUrls {
		if stopG {
			return [2]int64{}
		}
		if strings.TrimSpace(imgUrl) == "" {
			continue
		}
		l.Send(slog.LevelInfo, fmt.Sprintf(" - [%d/%d]Image URL : %s", index+1, imgInfo.ImageCount, imgUrl), log.LogFiles|log.LogStdouts)
		rp := strings.Split(imgUrl, ".") // 图片类型
		It := utils.DelSpeChar(imgInfo.ImageTitle)
		if len(imgInfo.ImageUrls) > 1 {
			_p = fmt.Sprintf("[%s_p%d]%s.%s", imgInfo.ImageId, index, It, rp[len(rp)-1])
		} else {
			_p = fmt.Sprintf("[%s]%s.%s", imgInfo.ImageId, It, rp[len(rp)-1])
		}

		filePath := path + _p
		wg.Add(1)
		go func(imgurl string, filepath string, imginfo sql.ImageInfo) {
			if !thread_ {
				thread <- [0]int{}
			}

			result := br.DownLoadImage(imgurl, filepath, refere, &imginfo)
			if result {
				l.Send(slog.LevelInfo, fmt.Sprintf("Image %s DownLoad Success And Save In: %s", imginfo.ImageId, filepath), log.LogFiles|log.LogStdouts)
				stda1.Add(1)
			} else {
				l.Send(slog.LevelWarn, fmt.Sprintf("Image %s DownLoad Fail", imginfo.ImageId), log.LogFiles|log.LogStdouts)
				stda2.Add(1)
			}
			wg.Done()
			state[0], state[1] = stda1.Load(), stda2.Load()
			if !thread_ {
				<-thread
			}
		}(imgUrl, filePath, *imgInfo)

	}
	wg.Wait()
	imgInfo.Status = true
	return
}

func (br *Br) DownLoadImage(imgurl, fileName, referer string, image *sql.ImageInfo) bool {

	for retryCount := 0; retryCount < NewWork.Retry; retryCount++ {
		if br.PerformDownload(imgurl, fileName) {
			if NewWork.DownloadDelay != 0 {
				time.Sleep(time.Duration(NewWork.DownloadDelay) * 100 * time.Millisecond)
			}
			return true
		}
		l.Send(slog.LevelWarn, fmt.Sprintf("Retry time %d", retryCount), log.LogStdouts)
		time.Sleep(time.Duration(NewWork.RetryWait) * 100 * time.Millisecond)
	}
	return false
}

func (br *Br) PerformDownload(imgUrl, fileNameSave string) bool {

	resp, err := GetPixivPage(imgUrl, 0)
	if err != nil {
		l.Send(slog.LevelError, "图片链接解析失败", log.LogFiles|log.LogStdouts)
		return false
	}

	f, err := os.OpenFile(fileNameSave, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		l.Send(slog.LevelError, "图片文件创建失败", log.LogFiles|log.LogStdouts)
		return false
	}
	df := bytes.NewReader(resp)
	_, err = f.Write(resp)
	if err != nil {
		l.Send(slog.LevelError, "图片写入失败", log.LogFiles|log.LogStdouts)
		return false
	}
	defer f.Close()

	st, _ := os.Stat(fileNameSave)
	if st.Size() == int64(df.Len()) {
		return true
	}
	fmt.Println("文件大小不对", imgUrl, st.Size(), int64(df.Len()))
	return false

}
