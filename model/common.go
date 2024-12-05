package model

import (
	"fmt"
	"gpixivImageDownload/conf"
	"strings"
	"sync"
)

var CommonVar = &Common{}

type Common struct {
	sync.Once
	Ck           []string `json:"ck"`
	R18          bool     `json:"r18"`
	MThread      bool     `json:"m_thread"`
	SkipIllus    bool     `json:"skipIllus"`    // 跳过插图
	SkipUgoira   bool     `json:"skipUgoira"`   // 跳过动图
	SkipManga    bool     `json:"skipManga"`    // 跳过漫画
	DownloadPath string   `json:"downloadPath"` //  下载路径
}

func (c *Common) Save() string {

	r := conf.Conf
	v := conf.Conf.Sub("CmtOption")
	fmt.Println(c.Ck)
	if len(c.Ck) != 0 {
		v.Set("Ck", strings.Join(c.Ck, "*"))
	}
	v.Set("R18", c.R18)
	v.Set("MThread", c.MThread)
	v.Set("SkipIllus", c.SkipIllus)
	v.Set("SkipUgoira", c.SkipUgoira)
	v.Set("SkipManga", c.SkipManga)
	fmt.Println("save功能待完善")
	e := r.WriteConfig()
	if e != nil {
		fmt.Println(e)
		return "common 保存失败"
	}
	return "common 保存成功"

	//return "save"
}
