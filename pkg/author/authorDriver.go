package author

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/dao/sql"
	errs "gpixivImageDownload/internal/err"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"gpixivImageDownload/pkg/core"
	"gpixivImageDownload/pkg/utils/browser"
	"log/slog"
	"slices"
	"strings"
)

var l *log.Logs

//var ipA = &inputAuthor{
//	User:  "author",
//	DwTop: 24,
//}
//
//type inputAuthor struct {
//	User       string
//	AuthorId   string
//	AuthorName string
//	DwTop      int
//	Tags       string
//	BookMark   bool
//	OffSet     string
//	Profile    bool
//}

//func init() {
//	cmdRoot.AddCommand(cmdAuthor)
//	if cmdAuthor.MarkFlagRequired("authorid") != nil && cmdAuthor.MarkFlagRequired("authorname") != nil {
//		l.Send(slog.LevelError, "author信息未输入", log.LogFiles|log.LogStdouts)
//		panic("")
//	}
//	f := cmdAuthor.Flags()
//	f.StringVarP(&ipA.AuthorId, "authorid", "i", ipA.AuthorId, "作者id")
//	f.StringVarP(&ipA.AuthorName, "authorname", "n", ipA.AuthorName, "作者名字")
//
//	f.IntVarP(&ipA.DwTop, "count", "c", ipA.DwTop, "默认下载数量")
//	f.IntVarP(&ipA.DwType, "type", "t", ipA.DwType, "1代表下载插画，2代表收藏，3代表全都要")
//
//}

func DownLoadAuth(ctx context.Context, cmpts *model.Common, autpts *model.Author, pip chan float64, result *string) {
	DwSucc := 0
	DwErr := 0
	olduid := []string{}
	oldname := []string{}
	ctxChild, _ := context.WithCancel(ctx)
	var url []string
	// 确认uid或name的有效性
	if len(autpts.AuthorId) != 0 {
		for _, auid := range autpts.AuthorId {
			if slices.Contains(olduid, auid) {
				continue
			}
			// https://www.pixiv.net/ajax/user/1122006
			_url := "https://www.pixiv.net/user/" + auid
			rb := browser.GetPixivPage(_url)
			if strings.Contains(string(rb), errs.GetMsg(errs.UIDERR)) {
				url = append(url, errs.GetMsg(errs.UIDERR))
			} else {
				url = append(url, _url)
			}
			olduid = append(olduid, auid)
			r := map[string]interface{}{}
			body := r["body"].(map[string]interface{})
			oldname = append(oldname, body["name"].(string))
		}
	}
	if len(autpts.AuthorName) != 0 {
		for _, aname := range autpts.AuthorName {
			auid := ""
			_url := fmt.Sprintf("https://www.pixiv.net/search_user.php?s_mode=s_usr&i=0&nick=%s", aname)
			rb := browser.GetPixivPage(_url)
			rep := bytes.NewReader(rb)
			doc, _ := goquery.NewDocumentFromReader(rep)
			sec, exist := doc.Find(".user-recommendation-item > a").Attr("href")
			if exist {
				auid = strings.Split(sec, "/")[2]
				if slices.Contains(olduid, auid) {
					continue
				}
				url = append(url, _url)
			} else {
				url = append(url, errs.GetMsg(errs.NameErr))
			}
			olduid = append(olduid, auid)
			oldname = append(oldname, aname)
		}

	}

	artistWork := &sql.AuthorWorks{}
	var rootPath string
	if autpts.DownloadPath != "" {
		rootPath = autpts.DownloadPath
	} else {
		rootPath = cmpts.DownloadPath
	}

	utils.QuoteOrCreateFile(rootPath)
	for idx, urls := range url {
		pageurl := ""
		if urls[:4] != "http" {
			continue
		}
		l.Send(slog.LevelInfo, fmt.Sprintf("Processing Member Id: %s", olduid[idx]), log.LogFiles|log.LogStdouts)
		// 空间，头像等 https://www.pixiv.net/ajax/user/%s/profile/all
		if autpts.BookMark {
			// 表示收藏 https://www.pixiv.net/ajax/user/6558698/illusts/bookmarks?tag=FGO&offset=0&limit=24&rest=show
			pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illusts/bookmarks?tag=%s&offset=%d&limit=%d&rest=show", olduid[idx], autpts.Tags, autpts.OffSet, autpts.DwTop)
			autpts.Tags = append(autpts.Tags, "BookMark")
		} else {
			if len(autpts.Tags) > 0 {
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=%s&offset=%d&limit=%d", olduid[idx], autpts.Tags, autpts.OffSet, autpts.DwTop)
			} else if cmpts.R18 == false {
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=R-18&offset=%s&limit=%s", olduid[idx], autpts.OffSet, autpts.DwTop)
			} else {
				// 表示作品 https://www.pixiv.net/ajax/user/6558698/illustmanga/tag?tag=&offset=0&limit=24
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=&offset=%d&limit=%d", olduid[idx], autpts.OffSet, autpts.DwTop)
			}
		}

		//if !autpts.Profile {
		//
		//}
		artistWork = core.GetAuthorWork(pageurl)
		for _, imgwork := range artistWork.Works {
			select {
			case <-ctxChild.Done():
				return
			default:
			}
			re, err := core.ProcessAuthImage(rootPath, imgwork, cmpts.R18, autpts.Tags)
			fmt.Println(re, err)
			DwSucc += re[0]
			DwErr += re[1]

		}
		pip <- float64(idx / len(url))
		//l.Send(slog.LevelInfo, fmt.Sprintf("Member Url %s\n", pageurl), log.LogFiles|log.LogStdouts)

	}
	*result = fmt.Sprintf("识别到%d个用户，成功下载%d个图片，下载失败%d个图片", len(url), DwSucc, DwErr)

}

//func GetRootPath(path, an, bookmark, tags, profile string) string {
//	path = path + Spt + "Artist" + Spt + an
//	// 表示收藏的图片
//	if bookmark != "" {
//		_path := path + Spt + bookmark
//		browser.CheckPathIsExit(_path)
//		return _path
//	} else {
//		// 表示 标签的图片
//		if tags != "" {
//			if tags == "R18" {
//				_path := path + Spt + tags
//				browser.CheckPathIsExit(_path)
//				return _path
//			}
//			_path := path + Spt + tags
//			browser.CheckPathIsExit(_path)
//			return _path
//
//		} else if profile != "" {
//			_path := path + Spt + profile
//			browser.CheckPathIsExit(_path)
//			return _path
//		} else {
//			_path := path + Spt + "all"
//			browser.CheckPathIsExit(_path)
//			return _path
//		}
//	}
//
//	return path + Spt + "other"
//
//}

//func GetInputInfo() (md, bm, ts, pf string) {
//	var url string
//
//	if {
//		bm = ""
//		fmt.Print("是否要下载用户插画的 profile(y/n): ")
//		fmt.Scanln(&t_)
//
//		if t_ == "y" {
//			pf = "true"
//		} else {
//			pf = ""
//		}
//	} else if t_ == "2" {
//		bm = "true"
//	} else {
//		goto info
//	}
//
//	fmt.Print("请输入要下载用户插画的 tags: ")
//	fmt.Scanln(&t_)
//
//	if len(t_) == 0 {
//		ts = ""
//		fmt.Println("tags is all")
//	} else if strings.Contains(t_, "R18") {
//		ts = "R-18"
//		fmt.Println("tags is R-18")
//	} else {
//		ts = t_
//		//fmt.Printf("tags is %s\n", ts)
//	}
//
//	fmt.Printf("medid: %s; bookmark: %s; tags: %s; pf: %s.\n", md, bm, ts, pf)
//	return md, bm, ts, pf
//
//}
