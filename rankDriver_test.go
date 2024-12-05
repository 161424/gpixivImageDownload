package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/core"
	"gpixivImageDownload/pkg/utils/browser"
	"testing"
	"time"
)

func TestDownloadRank(t *testing.T) {
	defer TrackTime(time.Now())
	var cmpts *model.Common = &model.Common{
		Ck:           []string{"first_visit_datetime_pc=2024-03-10%2021%3A51%3A20; p_ab_id=6; p_ab_id_2=7; p_ab_d_id=2108241874; yuid_b=NkI2c3c; privacy_policy_notification=0; a_type=0; b_type=1; login_ever=yes; privacy_policy_agreement=7; c_type=26; cf_clearance=_QG0FR5htU3hqSvQ.kfUfI6TMAKGe6y.sUNBFuEC424-1727945401-1.2.1.1-dZy7_ND1PgV2KGQFUpxrPdJMDh7LQ9G4hgHJWI7ngX_m_klBBa7FRRKnC1cjSKzmD42Xwt4qVXxn2ZTVYax36.7M6IL6a1eEPVCD7w8t4eXBWS5rwVZ6FaWTe9GYmE8cOFgYXP1MS5TgmxfXS4aetoPRETX9bCeER6ig6ZDMRBilpPcpwJqPKDYoC81X23G.kdCJQPGQK4XxrRQIfMwqlBa68knFG6kJ5v5lcpEkymrGuXfyxMaDEG7X1fxLyj7lmXdQxMyl.R7ABN4ni49gvEXlq48bH0PHU7XsYI2VLBd2L.3tUvChO21GA4weLcJJiovLqtjOwdcXvvMMKifQ9c0Puj8KwqR4Dm167N7CBZZ7mJvTGAmP3rEYvdC7vPNdtSDeCq6DBbpC5m3Pp..vgQ; PHPSESSID=9g3m7i534l30iql4unala8hlu1e84sln; cc1=2024-11-30%2002%3A41%3A50; __cf_bm=qM9gDz_c9LDrXY716Ot_8rNYAgC0UjBVzNQYmpgd29E-1732904260-1.0.1.1-2CasonHx5pr7MCbXWb8Dm_HnVdU_cwNDni05_3ZD547_Dp0X1dtrGgLqCiJDaZgclUdqAIammo_g8Sdlbjmr7yaNWFuptC3lzkuX.aKdtx0"},
		MThread:      true,
		DownloadPath: "D:\\编程\\新建文件夹\\rank\\daily",
		SkipManga:    true,
		SkipUgoira:   true,
	}
	browser.SetMutliHttps(cmpts.Ck)
	result, err := core.ProcessRankImage(context.TODO(), 1, 0, "124466110", cmpts.DownloadPath, cmpts.R18, "daily", "20241118", "all", cmpts.MThread == true)
	fmt.Println("end time:", result, err)
	result1, err1 := browser.GetPixivPage("https://www.pixiv.net/artworks/124466110", 0)
	fmt.Println("end time:", string(result1), err1)

}

func TrackTime(start time.Time) time.Duration {
	elapsed := time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	return elapsed
}

func TestKu(t *testing.T) {
	//var auid string
	//_url := fmt.Sprintf("https://www.pixiv.net/search_user.php?s_mode=s_usr&i=0&nick=%s", aname)
	//nick=ひづるめ%28Hidzzz%29&s_mode=s_usr
	ck := []string{"first_visit_datetime_pc=2024-03-10%2021%3A51%3A20; p_ab_id=6; p_ab_id_2=7; p_ab_d_id=2108241874; yuid_b=NkI2c3c; privacy_policy_notification=0; a_type=0; b_type=1; login_ever=yes; privacy_policy_agreement=7; c_type=26; PHPSESSID=23316552_zAdZn4Fbz6WTMDuWtJKoytWC7OgLl7Fl; device_token=a808143587bfae9f63a445c30d2128db; cf_clearance=AL8Vq6xdXqZmvff_dTH9tz1r0qQZs2IrpMGiHBweT2A-1733288571-1.2.1.1-baPWUYJ86Tz3tak3y9EtUFoTctm5zmpku0YaWh7Wb.3ppgCgn3EJrVXqlz4W3NoYBjKi3g_wFdRYDLnRSVg9F18v6meotcdJtTYxoKy.Xh5ejwwyLnvnlaCIqFAIHIwmq4KziieQTy8EI4TjLx.stmf.Ja3sEN1aMCC0TWm1Y.Oxi8q8Na.0907fjV0XQy4p_FR8vZ4KLTK5rxG6w0W6uZBGRWNPxrUsYzUKv.IcuYCnftmE9TsqHAqYL_EEYXBlORrMxvW0zw7FWoHWEr.Fcus5McBDiCNZ35W8NJd2mbhZSNhqvdYLFQyz948xtUmeBbfiFOD21ZCo2bJuCAkW7zfwIuufL_zJIVLD6AsjDOXi_5KRRyzJFNbFgD7920UH1iaufwSlLsxtFHq9Y3gf4Q; __cf_bm=2HaaCOA2Kong6zuIvFvr9vM3IvMQvHbewoTqaBydC8k-1733292521-1.0.1.1-HByTHpJuViKx8ULqJFnegUGltiwgnCOiZS3EM6QRVPZ5vrqupAOGlSrHePMe7tzbBgqXi3258xgqPzMUmAMkYMfH09eIhTKYGmFbUeWfhSA"}
	browser.SetMutliHttps(ck)
	aname := "ひづるめ(Hidzzz)"
	_url := fmt.Sprintf("https://www.pixiv.net/search/users?s_mode=s_usr&nick=%s", aname)
	fmt.Println(_url)
	rb, err := browser.GetPixivPage(_url, 0)
	fmt.Println(string(rb), err)
	if err != nil {
		fmt.Println(err)
	}

	rep := bytes.NewReader(rb)
	doc, _ := goquery.NewDocumentFromReader(rep)
	sec, exist := doc.Find("a.class=").Attr("href")
	//_, exist := doc.Find(".user-recommendation-item > a").Attr("href")
	fmt.Println("?", sec, exist)
}
