package author

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/dao/sql"
	errs "gpixivImageDownload/internal/err"
	log2 "gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"gpixivImageDownload/pkg/core"
	"gpixivImageDownload/pkg/utils/browser"
	"log/slog"
	url2 "net/url"
	"strconv"
	"strings"
)

var l = log2.Logger
var defaultUserUrl = "https://www.pixiv.net/users/"

func DownLoadAuth(ctx context.Context, cmpts *model.Common, autpts *model.Author, pip chan float64, result *string) {
	stopG := false
	go func() {
		<-ctx.Done()
		stopG = true
	}()
	autpts.OffSet = 0
	var DwSucc int64 = 0
	var DwErr int64 = 0
	olduid := []string{}
	//oldname := []string{}

	var url []string
	// 确认uid或name的有效性
	fmt.Println(autpts)
	if len(autpts.AuthorId) != 0 {
		for _, auid := range autpts.AuthorId {
			if _, err := strconv.Atoi(auid); err != nil {
				continue
			}
			// https://www.pixiv.net/ajax/users/1122006
			_url := defaultUserUrl + auid
			// 初次判断链接是否正确
			rb, err := browser.GetPixivPage(_url, 0)
			if err != nil {
				continue
			}
			if strings.Contains(string(rb), errs.GetMsg(errs.UIDERR)) {
				url = append(url, errs.GetMsg(errs.UIDERR))
			} else {
				url = append(url, _url)
			}
			olduid = append(olduid, auid)
		}
	}
	// https://www.pixiv.net/search/users?nick=ひづるめ%28Hidzzz%29&s_mode=s_usr
	if len(autpts.AuthorName) != 0 {
		for _, aname := range autpts.AuthorName {
			if strings.TrimSpace(aname) == "" {
				continue
			}
			fmt.Println("?", aname)
			//_url := fmt.Sprintf("https://www.pixiv.net/search_user.php?s_mode=s_usr&i=0&nick=%s", aname)
			//nick=ひづるめ(Hidzzz)&s_mode=s_usr

			// 若不对特殊字符进行编码，会导致错误
			aname = url2.QueryEscape(aname)
			_url := fmt.Sprintf("https://www.pixiv.net/search/users?nick=%s&s_mode=s_usr&", aname)

			// 获取到用户的uid
			rb, err := browser.GetPixivPage(_url, 0)

			if err != nil {
				continue
			}

			rep := bytes.NewReader(rb)
			doc, _ := goquery.NewDocumentFromReader(rep)
			hu := map[string]struct{}{}
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				// 检查 class 属性是否存在且不为空
				if s.AttrOr("data-ga4-entity-id", "") != "" {
					// 处理没有 class 属性或 class 属性为空的 <a> 标签
					link, exists := s.Attr("href")
					if exists {
						auid := strings.Split(link, "/")[2]
						hu[auid] = struct{}{}
					}
				}
			})
			for key, _ := range hu {
				olduid = append(olduid, key)
				_url = defaultUserUrl + key
				url = append(url, _url)
			}

		}

	}

	artistWork := &sql.AuthorWorks{}
	var rootPath string
	if autpts.DownloadPath != "" {
		rootPath = autpts.DownloadPath
	} else {
		rootPath = cmpts.DownloadPath
	}

	err := utils.QuoteOrCreateFile(rootPath)
	if err != nil {
		*result = err.Error()
		return
	}
	for idx, urls := range url {

		fmt.Println("url", url, urls, autpts.Tags)

		pageurl := ""
		if urls[:4] != "http" {
			continue
		}
		l.Send(slog.LevelInfo, fmt.Sprintf("Processing Member Id: %s", olduid[idx]), 3)
		// 空间，头像等 https://www.pixiv.net/ajax/user/%s/profile/all
		if autpts.BookMark {
			// 表示收藏 https://www.pixiv.net/ajax/user/6558698/illusts/bookmarks?tag=FGO&offset=0&limit=24&rest=show
			pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illusts/bookmarks?tag=%s&offset=%d&limit=%d&rest=show", olduid[idx], autpts.Tags, autpts.OffSet, autpts.DwTop)
			autpts.Tags = "BookMark"
		} else {
			if len(autpts.Tags) > 0 {
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=%s&offset=%d&limit=%d", olduid[idx], autpts.Tags, autpts.OffSet, autpts.DwTop)
			} else if cmpts.R18 == false {
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=R-18&offset=%d&limit=%s", olduid[idx], autpts.OffSet, autpts.DwTop)
			} else {
				// 表示作品 https://www.pixiv.net/ajax/user/6558698/illustmanga/tag?tag=&offset=0&limit=24
				pageurl = fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/illustmanga/tag?tag=&offset=%d&limit=%d", olduid[idx], autpts.OffSet, autpts.DwTop)
			}
		}
		fmt.Println("pageurl", pageurl)
		artistWork = core.GetAuthorWork(pageurl)
		if artistWork.Error != "" {
			fmt.Println(artistWork.Error)
			continue
		}
		for _, imgwork := range artistWork.Works {
			if stopG {
				*result += "下载已停止"
				return
			}

			re, err := core.ProcessAuthImage(ctx, rootPath, imgwork, cmpts.R18, autpts.Tags)
			if err != nil {
				fmt.Println(re, err)
			}

			DwSucc += re[0]
			DwErr += re[1]
			pip <- float64(idx / len(url))
		}

		//l.Send(slog.LevelInfo, fmt.Sprintf("Member Url %s\n", pageurl), log.LogFiles|log.LogStdouts)

	}
	pip <- 1.0
	*result += fmt.Sprintf("识别到%d个用户，成功下载%d个图片，下载失败%d个图片", len(url), DwSucc, DwErr)

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
