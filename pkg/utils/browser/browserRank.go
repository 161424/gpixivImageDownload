package browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/dao/sql"
	"gpixivImageDownload/internal/addr"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Br struct {
	Client *http.Client
}

var Brs = &Br{}
var l *log.Logs
var refere = "https://www.pixiv.net"

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
	Timestamp string
	Illust    map[string]map[string]interface{}
	User      map[string]map[string]interface{}
}

func init() {
	l = log.NewSlogGroup("Browser")
	sockets, _ := url.Parse(addr.Proxy.Ip + ":" + addr.Proxy.Port)
	l.Send(slog.LevelInfo, fmt.Sprintf("Proxy:%s", sockets), log.LogFiles|log.LogStdouts)
	//fmt.Println(Brs)
	Brs.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			Proxy:           http.ProxyURL(sockets),
		},
		Timeout: 10 * time.Second,
	}

}

func hashString(s string) uint64 {
	var h uint64 = 14695981039346656037 // offset
	for i := 0; i < len(s); i++ {
		h = h ^ uint64(s[i])
		h = h * 1099511628211 // prime
	}
	return h
}

func GetPixivPage(urls string) []byte {
	client := Brs.Client
	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", urls, nil)
	if err != nil {
		return nil
	}
	//fmt.Printf("conf.Header.UserAgent, %T", conf.Header.UserAgent)
	req.Header.Set("User-Agent", addr.Header.UserAgent)
	req.Header.Set("Referer", refere)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for _, i := range strings.Split(sql.A.Cookies, ";") {
		a := strings.Split(i, "=")
		req.AddCookie(&http.Cookie{Name: a[0], Value: a[1]})
	}

	//ck.reps.Store(hashString(sql.DefaultAuth().Cookies),req)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	bts := bytes.Buffer{}

	bts.ReadFrom(resp.Body)
	//fmt.Println(a, b, bts.Bytes(), resp.Body)
	defer resp.Body.Close()
	return bts.Bytes()
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

func (br *Br) GetPixivPage(url string) (RJ *RpJs, s string) {

	//url = br.fixUrl(url, true)
	rb := GetPixivPage(url)
	l.Send(slog.LevelDebug, url, log.LogFiles|log.LogStdouts)
	if strings.Contains(string(rb), "Just a moment...") {
		s = "受到了CloudFlare 5s盾的管控， 请检查cookie和梯子是否有效"
		return
	}
	//fmt.Printf()
	if strings.Contains(string(rb), "不在排行榜统计范围内") {
		s = "不在排行榜统计范围内"
		return
	}

	if err := json.Unmarshal(rb, &RJ); err != nil {
		s = "解析错误，请检查参数是否匹配"
		return
	}

	return
}

func (br *Br) GetPixivRanking(mode, date string, comm model.Common, page int) (*RpJs, string) {
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

	df, err := br.GetPixivPage(url)

	if err != "" {
		return nil, err
	}

	return df, ""
}

func (br *Br) AnalysisUrl(url string) (*goquery.Document, error) {

	resp := GetPixivPage(url)
	df := bytes.NewReader(resp)
	doc, err := goquery.NewDocumentFromReader(df)
	if err != nil {
		return nil, err
	}
	br.CheckWeb(doc)
	return doc, nil
}

func (br *Br) GetImageInfo(imageId string) (imageInfo *sql.ImageInfo, err error) {

	var doc *goquery.Document
	l.Send(slog.LevelDebug, fmt.Sprintf("Getting image page: %s", imageId), log.LogStdouts)
	// https://www.pixiv.net/artworks/115906527
	url := fmt.Sprintf("https://www.pixiv.net%s/artworks/%s/", "", imageId)
	if doc, err = br.AnalysisUrl(url); err != nil {
		return nil, err
	}

	imageInfo, err = ParseRankImage(doc, imageId)

	if err != nil {
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

func (br *Br) DownloadImage(imgInfo *sql.ImageInfo, path string) (state [2]int) {

	var wg = sync.WaitGroup{}
	var _p string
	for index, imgUrl := range imgInfo.ImageUrls {
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
			result := br.DownLoadImage(imgurl, filepath, refere, &imginfo)
			if result {
				l.Send(slog.LevelInfo, fmt.Sprintf("Image %s DownLoad Success And Save In: %s", imginfo.ImageId, filepath), log.LogFiles|log.LogStdouts)
				state[0]++
			} else {
				l.Send(slog.LevelWarn, fmt.Sprintf("Image %s DownLoad Fail", imginfo.ImageId), log.LogFiles|log.LogStdouts)
				state[1]++
			}
			wg.Done()
		}(imgUrl, filePath, *imgInfo)
	}
	wg.Wait()
	imgInfo.Status = true
	return
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

	resp := GetPixivPage(imgUrl)

	f, err := os.OpenFile(fileNameSave, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		l.Send(slog.LevelError, "图片文件创建失败", log.LogFiles|log.LogStdouts)
		return false
	}
	df := bytes.NewReader(resp)
	//resp.Read(buf.Bytes())
	//buf.Write(resp)
	f.Write(resp)
	defer f.Close()

	st, _ := os.Stat(fileNameSave)
	if st.Size() == int64(df.Len()) {
		return true
	}
	return false

}
