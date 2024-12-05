package browser

import (
	"encoding/json"
	"gpixivImageDownload/dao/sql"
)

type R struct {
	Error   interface{}      `json:"error"`
	Message string           `json:"message"`
	Body    *sql.AuthorWorks `json:"body"`
}

func getAuthorProfile(url string) *sql.AuthorWorks {
	//url := fmt.Sprintf("https://www.pixiv.net/ajax/user/%s", artistId)

	//url := fmt.Sprintf("https://www.pixiv.net/ajax/user/1122006/profile/all", artistId)
	body, err := GetPixivPage(url, 0)

	//fmt.Println(string(body))
	var r = &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		r.Body.Error = err.Error()
		return r.Body
	}
	//fmt.Println(string(body))

	return r.Body

}
