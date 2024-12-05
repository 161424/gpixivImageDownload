package browser

import (
	"encoding/json"
	"fmt"
	"gpixivImageDownload/dao/sql"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model/utils"
	"log/slog"
	"strings"
	"sync"
)

func (br *Br) GetArtistPage(url string) (body *sql.AuthorWorks, s string) {

	url = br.fixUrl(url, true)
	l.Send(slog.LevelDebug, url, log.LogFiles|log.LogStdouts)
	rb, err := GetPixivPage(url, 0)
	if err != nil {
		l.Send(slog.LevelError, err.Error(), log.LogFiles|log.LogStdouts)
		return
	}

	s = ""
	if strings.Contains(string(rb), "Just a moment...") {
		s = "受到了CloudFlare 5s盾的管控， 请检查cookie和梯子是否有效"
		return
	}

	if err := json.Unmarshal(rb, &body); err != nil {
		s = "解析错误，请检查参数是否匹配"
		return
	}

	return
}

func (br *Br) GetAuthorProfile(artistId string) *sql.AuthorWorks {
	return getAuthorProfile(artistId)
}

// author db
//func (br *Br) GetAuthorImageInfo(url string) (imageInfo *sql.Artist, err error) {
//	//todo 调用数据库
//
//	var doc *goquery.Document
//	if doc, err = br.AnalysisUrl(url); err != nil {
//		return nil, err
//	}
//
//	imgInfo, err := ParseAuthorInfo(doc)
//	if err != nil {
//		return nil, err
//	}
//	if imgInfo.ImageMode == "ugoira_view" {
//		ugoiraMetaUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/ugoira_meta", imageId)
//		// ugoira 代表是动图，暂不进行过多考虑
//		fmt.Println(ugoiraMetaUrl)
//	}
//	imageInfo.ImageId = imageId
//	return
//}

func (br *Br) DownloadAuthorImage(imgInfo *sql.ImageInfo, path string) {
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

			} else {
				l.Send(slog.LevelWarn, fmt.Sprintf("Image %s DownLoad Fail", imginfo.ImageId), log.LogFiles|log.LogStdouts)
			}
			wg.Done()
		}(imgUrl, filePath, *imgInfo)
	}
	wg.Wait()
	imgInfo.Status = true

}
