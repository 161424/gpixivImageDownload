package browser

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/dao/sql"
	"gpixivImageDownload/internal/addr"
	"gpixivImageDownload/model/utils"
	"path/filepath"
	//"github.com/chen/download_pixiv_pic/pkg/Browser"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var NewWork = addr.NewNetWork()

func IsNotLoggedIn(page *goquery.Document) bool {
	check := page.Find(".signup_button")
	if len(check.Text()) > 0 {
		return true
	}
	check = page.Find(".ui-button _signup")
	if len(check.Text()) > 0 {
		return true
	}
	return false
}

func IsNeedPermission(page *goquery.Document) bool {
	errorMessages := []string{"この作品は.+さんのマイピクにのみ公開されています|この作品は、.+さんのマイピクにのみ公開されています",
		"This work is viewable only for users who are in .+\\'s My pixiv list",
		"Only .+\\'s My pixiv list can view this.",
		"<section class=\"restricted-content\">"}
	return HaveString(page, errorMessages)
}

func HaveString(page *goquery.Document, s []string) bool {
	for _, str := range s {
		reg := regexp.MustCompile(str)
		result := reg.FindAllStringSubmatch(page.Text(), -1)
		for _, r := range result {
			if len(r[1]) > 0 {
				return true
			}
		}
	}
	return false
}

func IsNeedAppropriateLevel(page *goquery.Document) bool {
	errorMessages := []string{"該当作品の公開レベルにより閲覧できません。"}
	return HaveString(page, errorMessages)
}

func IsDeleted(page *goquery.Document) bool {
	errorMessages := []string{"該当イラストは削除されたか、存在しないイラストIDです。|該当作品は削除されたか、存在しない作品IDです。",
		"この作品は削除されました。",
		"The following work is either deleted, or the ID does not exist.",
		"This work was deleted.",
		"Work has been deleted or the ID does not exist."}
	return HaveString(page, errorMessages)
}

func IsGuroDisabled(page *goquery.Document) bool {
	errorMessages := []string{"表示されるページには、18歳未満の方には不適切な表現内容が含まれています。",
		"The page you are trying to access contains content that may be unsuitable for minors"}
	return HaveString(page, errorMessages)
}

//func IsErrorExist(page *goquery.Document) {
//	//check := page.Find("span.error")
//}

func PrintInfo(img *sql.ImageInfo) {
	fmt.Printf("  User Name: %s\n", img.UserName)
	fmt.Printf("  Image Title: %s\n", img.ImageTitle)
	fmt.Printf("  Rank: %d\n", img.ImageRank)
	fmt.Printf("  Tags Title: %s\n", img.ImageTags)
	fmt.Printf("  Translated Tags: %s\n", img.TranslationTag)
	fmt.Printf("  Date: %s\n", img.WorksDateDateTime)
	fmt.Printf("  Mode: %s\n", img.ImageMode)
	fmt.Printf("  Urls: %s\n", img.ImageUrls)

	if img.ImageMode == "manga" {
		fmt.Printf("  Pages: %d\n", img.ImageCount)
	}
	fmt.Printf("  Bookmarks %d\n", img.BookmarkCount)
}

func MakeFilename(img *sql.ImageInfo, path string, r18 bool, user, mode, date, content string) (string, error) {
	fn := utils.DelSpeChar(img.UserName)
	var r string
	r = fmt.Sprintf("%s(%s)", fn, img.UserID)
	_user := user
	if r18 {
		_user = _user + "_r18"
	}
	if user == "rank" {
		r = filepath.Join(_user, mode, fmt.Sprintf("%s-%s", date, content), r)
	} else if user == "author" {
		// 作者的收藏，作品等等
		r = filepath.Join(_user, r, content)
	} else if user == "tag" {
		r = filepath.Join(_user, fmt.Sprintf("%s-%s", date, content), r)
	}

	path = filepath.Join(path, r)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
			if err != nil {
				return "", err
			}
		}

	}
	fmt.Printf("  Creating directory: %s. Successful!\n", path)

	return path, nil
}

// 解析图片页面
func ParseRankImage(doc *goquery.Document, imgId string) (*sql.ImageInfo, error) {
	var err error
	imageInfo := &sql.ImageInfo{}

	r, _ := doc.Find("meta#meta-preload-data").Attr("content")
	var con = &Cont{}
	var root Illust
	var ok bool
	err = json.Unmarshal([]byte(r), con)
	if root, ok = con.Illust[imgId]; ok == false || err != nil {
		fmt.Println("abstractRank err:", err, r)
		return nil, err
	}

	imageCount := int(root.PageCount)

	tempUrl := root.Urls["original"]
	tempResizedUrl := root.Urls["regular"]
	imageInfo.ImageCount = imageCount

	//  不够齐全
	if imageCount == 1 {
		if strings.Contains(tempUrl, "ugoira") { // 动图
			imageInfo.ImageMode = "ugoira_view"
			tempUrl = strings.Replace(tempUrl, "/img-original/", "/img-zip-ugoira/", -1)
			tempUrl = strings.Split(tempUrl, "_ugoira0")[0]
			tempUrl = tempUrl + "_ugoira1920x1080.zip"
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, tempUrl)

			tempResizedUrl = strings.Replace(tempResizedUrl, "/img-original/", "/img-zip-ugoira/", -1)
			tempResizedUrl = strings.Split(tempResizedUrl, "_ugoira0")[0]
			tempResizedUrl = tempResizedUrl + "_ugoira600x600.zip"
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, tempResizedUrl)
		} else {
			imageInfo.ImageMode = "big"
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, tempUrl)
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, tempResizedUrl)
		}
	} else {
		imageInfo.ImageMode = "manga"
		for i := 0; i < imageCount; i++ {
			_tempUrl := strings.Replace(tempUrl, "_p0", "_p"+strconv.Itoa(i), -1)
			imageInfo.ImageUrls = append(imageInfo.ImageUrls, _tempUrl)
			_tempResizedUrl := strings.Replace(tempResizedUrl, "_p0", "_p"+strconv.Itoa(i), -1)
			imageInfo.ImageResizedUrls = append(imageInfo.ImageResizedUrls, _tempResizedUrl)
		}

	}
	imageInfo.ImageTitle = root.IllustTitle
	//imageInfo.ImageCaption = (root["illustComment"]).(string)
	if root.SeriesNavData == nil { // 是否是系列
		imageInfo.SeriesNavData = map[string]string{"nil": "nil"}
	} else {
		mss := make(map[string]string, 6)
		mss["seriesType"] = root.SeriesNavData["seriesType"].(string)
		mss["seriesId"] = root.SeriesNavData["seriesId"].(string)
		mss["title"] = root.SeriesNavData["title"].(string)
		mss["order"] = strconv.Itoa(root.SeriesNavData["order"].(int))
		if root.SeriesNavData["prev"] == nil {
			mss["prev"] = ""
		} else {
			mss["prev"] = root.SeriesNavData["prev"].(map[string]interface{})["id"].(string)
		}

		if root.SeriesNavData["next"] == nil {
			mss["next"] = ""
		} else {
			mss["next"] = root.SeriesNavData["next"].(map[string]interface{})["id"].(string)
		}

		imageInfo.SeriesNavData = mss
	}

	imageInfo.JdRtv = root.ViewCount // 查看人数
	imageInfo.JdRtc = root.LikeCount // 喜欢人数

	if root.Tags.Tags != nil {
		p := root.Tags.Tags
		for _, tagp := range p {
			//tagp := tag
			imageInfo.ImageTags = append(imageInfo.ImageTags, (tagp["tag"]).(string))
			if v, isok := tagp["translation"]; isok == true {
				if v != nil {
					for _, valu := range v.(map[string]interface{}) {
						imageInfo.TranslationTag = append(imageInfo.TranslationTag, (valu).(string))
					}
				}
			}
		}
	}

	imageInfo.WorksDateDateTime, _ = time.Parse(time.RFC3339, root.CreateDate)
	imageInfo.WorksResolution = fmt.Sprintf("%.0fx%.0f", root.Width, root.Height)

	imageInfo.BookmarkCount = int64(root.BookmarkCount)
	imageInfo.ImageResponseCount = root.ResponseCount

	imageInfo.AiType = int(root.AiType) //1 == non-AI, 2 == AI-generated
	if imageInfo.AiType == 2 {
		imageInfo.ImageTags = append(imageInfo.ImageTags, "AI-generated")
	}

	imageInfo.UserID = root.UserId
	imageInfo.UserName = root.UserName
	//imageInfo.UserAccount = (root["userAccount"]).(string)

	return imageInfo, nil

}
