package err

import "fmt"

var MsgFlags = map[int]string{
	UnknownErr:         "未知错误",
	ConfigFileNotFound: "config.yaml未发现",
	ConfigFileReadErr:  "配置文件读取错误",
	ConfigReadErr:      "配置文件发现未知错误",
	ConfigReadSuccess:  "%s 读取成功",
	UIDERR:             "抱歉，您当前所寻找的个用户已经离开了pixiv, 或者这ID不存在。",
	NameErr:            "用户名错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	s := msg
	if !ok {
		s = MsgFlags[UnknownErr]
	}
	return s
}

func ErrWrap(code int) error {
	msg, ok := MsgFlags[code]
	s := msg
	if !ok {
		s = MsgFlags[UnknownErr]
	}

	return fmt.Errorf(s)
}
