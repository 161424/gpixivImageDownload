package model

type Common struct {
	Ck           []string `json:"ck"`
	R18          bool     `json:"r18"`
	MThread      bool     `json:"m_thread"`
	SkipIllus    bool     `json:"skipIllus"`    // 跳过插图
	SkipUgoira   bool     `json:"skipUgoira"`   // 跳过动图
	SkipManga    bool     `json:"skipManga"`    // 跳过漫画
	DownloadPath string   `json:"downloadPath"` //  下载路径

}
