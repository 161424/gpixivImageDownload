package browser

import (
	"encoding/json"
	"fmt"
	"gpixivImageDownload/dao/sql"
	"log"
	"log/slog"
	"reflect"
)

func getAuthorProfile(url string) *sql.AuthorWorks {
	//url := fmt.Sprintf("https://www.pixiv.net/ajax/user/%s", artistId)
	var artist *sql.AuthorWorks
	var r map[string]interface{}
	//var err error
	//rb := GetPixivPage(url)
	//err = json.Unmarshal(rb, &r)
	//if err != nil {
	//	artist.Error = err.Error()
	//	return artist
	//}
	//rbb, _ := json.Marshal(r["body"])
	//err = json.Unmarshal(rbb, artist)
	//if err != nil {
	//	artist.Error = err.Error()
	//	l.Send(slog.LevelError, fmt.Sprintf("解析失败,err=%s", err), 3)
	//	return artist
	//}

	//url := fmt.Sprintf("https://www.pixiv.net/ajax/user/1122006/profile/all", artistId)
	err := json.Unmarshal(GetPixivPage(url), &r)
	if err != nil {
		artist.Error = err.Error()
		return artist
	}
	bodys := r["body"].(map[string]any)
	artist, err = ParseAuthorInfo(bodys)
	//for _, i := range []string{"illusts", "manga", "novels", "mangaSeries", "novelSeries"} {
	//	if value, ok := bodys[i].(map[string]any); ok == true {
	//		body(value, artist, i)
	//	}
	//}
	//pickup := bodys["pickup"]
	//if  {
	//	if value, ok := pickup.([]map[string]any); ok == true {
	//		PickUp(artist, value)
	//	}
	//
	//}

	return artist

}

func bodys(value map[string]any, art *sql.Artist, tags string) {
	r := reflect.TypeOf(art).Elem()
	for i := 0; i < r.NumField(); i++ {
		if r.Field(i).Tag.Get("json") == tags {
			w := []string{}
			for k, _ := range value {
				w = append(w, k)
			}
			reflect.ValueOf(art).Elem().Field(i).Set(reflect.ValueOf(w))

		}
	}
}

func PickUp(art *sql.Artist, r []map[string]any) {
	for _, i := range r {
		art.DwImageInfo = append(art.DwImageInfo, ReflectWork(i))
	}
}

func ReflectWork(v map[string]any) (w sql.Work) {
	vb, err := json.Marshal(v)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("err=%s", err), 3)
	}
	err = json.Unmarshal(vb, &w)
	if err != nil {
		l.Send(slog.LevelError, fmt.Sprintf("err=%s", err), 3)
	}
	return
}

func ParseAuthorInfo(r map[string]interface{}) (*sql.AuthorWorks, error) {
	var artist = &sql.AuthorWorks{}

	if (r["error"]).(bool) {
		artist.Error = fmt.Sprint((r["error"]).(bool))
		log.Panicln("err for get page")
	}
	if r["body"] == nil {
		log.Panicln("Missing body content, possible artist id doesn't exists.")
		return artist, nil
	}
	body := (r["body"]).(map[string]interface{})
	artist.Total = body["total"].(int)
	ParseImages(artist, body)
	//ParseMangaList(artist, body)
	//ParseNovelList(artist, body)
	//artist.ArtistSomeInfo(mid)

	return artist, nil

}

func ParseImages(artist *sql.AuthorWorks, body map[string]interface{}) {

	if body["works"] != nil {
		work := (body["works"]).([]interface{})
		for _, img := range work {
			artist.Works = append(artist.Works, img.(sql.Work))
		}

	}
}

//func ParseMangaList(artist *sql.Artist, body map[string]interface{}) {
//	if len(body) != 0 && body["mangaSeries"] != nil {
//		ms := (body["mangaSeries"]).([]map[string]interface{})
//		for _, i := range ms {
//			artist.MangaSeries = append(artist.MangaSeries, (i["id"]).(string))
//		}
//	}
//}
//
//func ParseNovelList(artist *sql.Artist, body map[string]interface{}) {
//	if len(body) != 0 && body["novelSeries"] != nil {
//		ms := (body["novelSeries"]).([]map[string]interface{})
//		for _, i := range ms {
//			artist.MangaSeries = append(artist.MangaSeries, (i["id"]).(string))
//		}
//	}
//}
