package core

import (
	"fmt"
	"gpixivImageDownload/dao/sql"
	log "gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/utils/browser"
	"log/slog"
	"os"
)

var Br *browser.Br
var l log.Logs

func GetPixivRanking(mode, date string, cmpts model.Common, page int) (*browser.RpJs, string) {
	return Br.GetPixivRanking(mode, date, cmpts, page)
}

func GetAuthorWork(url string) (imageInfo *sql.AuthorWorks) {
	return Br.GetAuthorProfile(url)
}

func ProcessRankImage(pgcount int, imageId, rootpath string, r18 bool, model, date, content string) (state [2]int, err error) {

	var fileName string

	// 本地判重代码
	//var inDb = false
	//// GetImageInfo
	imgInfo, err := Br.GetImageInfo(imageId)
	//if err != nil {
	//	l.Send(slog.LevelError, fmt.Sprintf("获取图片信息失败,err=%s", err), log.LogFiles|log.LogStdouts)
	//	return err
	//}
	imgInfo.ImageCount = pgcount
	//// CreateDB or pass
	//
	//err, path := globalOptions.DB.SelectImageByImageId(imageId)
	//
	//// 只有检测到真实存在的文件，num才会+1
	//if err == nil && path != "" {
	//	inDb, err = utils.CheckImageStatus("", pgcount)
	//}
	//
	//if inDb == true {
	//	l.Send(slog.LevelInfo, fmt.Sprintf("Already downloaded in DB: %s", imageId), log.LogFiles|log.LogStdouts)
	//	return nil
	//}

	if fileName, err = browser.MakeFilename(imgInfo, rootpath, r18, "rank", model, date, content); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("image文件夹创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return
	} else {
		fileName += string(os.PathSeparator)
		imgInfo.SavePath = fileName
	}

	// 暂不设置下载图像大小，默认默认
	browser.PrintInfo(imgInfo)
	state = Br.DownloadImage(imgInfo, fileName)

	//if ok := globalOptions.DB.SaveImageId(imgInfo); ok != nil {
	//	l.Send(slog.LevelError, fmt.Sprintf("图片信息存储失败，err=%s", err), log.LogFiles|log.LogStdouts)
	//}

	return
}

func ProcessAuthImage(rootpath string, imgwork sql.Work, r18 bool, tags []string) (state [2]int, err error) {

	var fileName string

	imgInfo, err := Br.GetImageInfo(imgwork.Id)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("获取图片信息失败,err=%s", err), log.LogFiles|log.LogStdouts)
		return
	}
	content := ""
	for _, tag := range tags {
		content += tag
		content += "-"
	}
	content = content[:len(content)-1]
	if fileName, err = browser.MakeFilename(imgInfo, rootpath, r18, "author", "", "", content); err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("image文件夹创建失败，err=%s", err), log.LogFiles|log.LogStdouts)
		return
	} else {
		fileName += string(os.PathSeparator)
		imgInfo.SavePath = fileName
	}

	// 下载

	browser.PrintInfo(imgInfo)
	state = Br.DownloadImage(imgInfo, fileName)
	return
	//
	//df := bytes.NewReader(resp)
	//
	//resp := Br.GetAuthorImageInfo(url)
	//var r map[string]interface{}
	//json.Unmarshal(resp, &r)
	//da.GetPA(memberId, r, false, offset, limit)
	//ArtistSomeInfo(da)
	//if len(da.ImageList) == 0 {
	//	fmt.Printf("No images for Member Id: %s, from Bookmark: %s", memberId, bookmark)
	//}
	////da.ReferenceImageId = da.ImageList[0]
	//da.GetMemberInfoWhitecube(memberId, bookmark)
	//
	//fmt.Printf("Member Name %s\n", da.ArtistName)
	//fmt.Printf("Member Avatar %s\n", da.ArtistAvatar)
	//fmt.Printf("Member Backgrd %s\n", da.ArtistBackground)
	//var printOffsetStop int
	//if offsetStop < da.TotalImages {
	//	printOffsetStop = offsetStop
	//} else {
	//	printOffsetStop = da.TotalImages
	//}
	//fmt.Printf("Processing images from %d to %d of %d\n", offset+1, printOffsetStop, da.TotalImages)
	//
	//if !da.HaveImages {
	//	log.Printf("No image found for: %d\n", memberId)
	//}
	//
	//var dnlist []string
	//var id int
	//var ns = 0
	//for {
	//	rootpath := GetRootPath(Rootpath, da.ArtistName, bookmark, tags, profile)
	//	fmt.Println(rootpath)
	//	for id, imgid = range da.ImageList {
	//		up := fmt.Sprintf("[ %d of %d ]", id+1, printOffsetStop)
	//		if da.TotalImages > 0 {
	//			printOffsetStop -= offset
	//		} else {
	//			printOffsetStop = (sp-1)*20 + len(da.ImageList)
	//		}
	//		for retryCount := 0; retryCount < browser.NewWork.Retry; retryCount++ {
	//			fmt.Printf("MemberId: %s Page: %d Post %d of %s\n", memberId, sp, up, da.TotalImages)
	//
	//			result := br.ProcessImage(imgid, rootpath, "mb")
	//			if result == "YES" {
	//				ns += 1
	//				dnlist = append(dnlist, imgid)
	//				break
	//			}
	//		}
	//	}
	//	sp += 1
	//	if sp > ep {
	//		break
	//	}
	//}
	//
	//// 更新artist数据库数据
	//handle.DB.UpdateArtist(da)
	//da.LocalImagesList = dnlist
	//da.LocalImages = ns
	//fmt.Printf("last image_id: %s\n", imgid)
	//fmt.Printf("Member_id: %s  completed: %s\n")
}

//func ArtistSomeInfo(da *artist.PixivArtist) {
//	// https://www.pixiv.net/ajax/user/83739
//	// https://www.pixiv.net/ajax/user/6558698
//
//	//da.ArtistBackground = ((r["body"]).(map[string]interface{})["background"]).(string)
//}
