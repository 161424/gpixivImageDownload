package conf

import (
	"fmt"
	"github.com/spf13/viper"
	e "gpixivImageDownload/internal/err"
	"gpixivImageDownload/log"
	"log/slog"
	"os"
	"path/filepath"
)

var ConfigData map[string]map[string]interface{}
var Conf *viper.Viper

func init() {
	l := log.NewSlogGroup("config")

	code := e.ConfigReadSuccess

	path, _ := os.Getwd()
	//path = filepath.Dir(path)
	fmt.Println(path)
	viper.SetConfigFile(filepath.Join(path, "conf", "config.yaml"))
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			code = e.ConfigFileNotFound
		} else {
			code = e.ConfigFileReadErr
		}
		l.Send(slog.LevelError, e.GetMsg(code), log.LogFiles|log.LogStdouts)
		return
	}

	l.Send(slog.LevelInfo, fmt.Sprintf(e.GetMsg(code), "config"), log.LogFiles|log.LogStdouts)
	Conf = viper.GetViper()

}

func GetConf() *viper.Viper {
	return Conf
}

//func SaveCmtOption(cmtOption model.Common) {
//	tpname := reflect.TypeOf(cmtOption)
//	for i := 0; i < tpname.Len(); i++ {
//		key := tpname.Field(i)
//
//
//	}
//}
