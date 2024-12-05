package core

import (
	"context"
	"fmt"
	"gpixivImageDownload/dao/sql"
	log "gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/utils/browser"
	"log/slog"
	"os"
)

var Br = *browser.Brs
var l log.Logs

func GetPixivRanking(mode, date string, cmpts *model.Common, page int) (*browser.RpJs, string) {
	return Br.GetPixivRanking(mode, date, cmpts, page)
}

func GetAuthorWork(url string) (imageInfo *sql.AuthorWorks) {
	return Br.GetAuthorProfile(url)
}

func ProcessRankImage(ctx context.Context, pgcount, imageRank int, imageId, rootpath string, r18 bool, model, date, content string, thread_ bool) (state [2]int64, err error) {

	var fileName string
	// url预处理
	imgInfo, err := Br.GetImageInfo(imageId)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("获取image信息失败，err=%s", err), log.LogFiles|log.LogStdouts)
		state[1]++
		return
	}
	imgInfo.ImageRank = imageRank
	imgInfo.ImageCount = pgcount

	if fileName, err = browser.MakeFilename(imgInfo, rootpath, r18, "rank", model, date, content); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("image文件夹创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		state[1]++
		return
	} else {
		fileName += string(os.PathSeparator)
		imgInfo.SavePath = fileName
	}

	//
	browser.PrintInfo(imgInfo)
	// 实际下载

	state = Br.DownloadImage(ctx, imgInfo, fileName, thread_)

	//if ok := globalOptions.DB.SaveImageId(imgInfo); ok != nil {
	//	l.Send(slog.LevelError, fmt.Sprintf("图片信息存储失败，err=%s", err), log.LogFiles|log.LogStdouts)
	//}

	return
}

func ProcessAuthImage(ctx context.Context, rootpath string, imgwork sql.Work, r18 bool, tags string) (state [2]int64, err error) {

	var fileName string

	imgInfo, err := Br.GetImageInfo(imgwork.Id)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("获取图片信息失败,err=%s", err), log.LogFiles|log.LogStdouts)
		return
	}
	content := tags
	if fileName, err = browser.MakeFilename(imgInfo, rootpath, r18, "author", "", "", content); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("image文件夹创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return
	} else {
		fileName += string(os.PathSeparator)
		imgInfo.SavePath = fileName
	}

	// 下载

	browser.PrintInfo(imgInfo)
	state = Br.DownloadImage(ctx, imgInfo, fileName, true)
	return
}
