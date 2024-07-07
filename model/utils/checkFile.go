package utils

import (
	"log/slog"
	"os"
	"strings"
)

func QuoteOrCreateFile(n string) {
	_, err := os.Stat(n)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(n, os.ModePerm)
			if err != nil {
				l.Send(slog.LevelError, "创建文件夹出现错误", 3)
			}
		}
	}
}

func CheckImageStatus(str string, ct int) (bool, error) {
	dirPath, err := os.ReadDir(str)
	if err == nil {
		if len(dirPath) == ct {
			return true, nil
		}
	}
	return false, nil

}

func DelSpeChar(DF string) string {
	DF = strings.Replace(DF, "\\", "", -1)
	DF = strings.Replace(DF, "/", "", -1)
	DF = strings.Replace(DF, ":", "", -1)
	DF = strings.Replace(DF, "*", "", -1)
	DF = strings.Replace(DF, "?", "", -1)
	DF = strings.Replace(DF, "\\\"", "", -1)
	DF = strings.Replace(DF, "<", "", -1)
	DF = strings.Replace(DF, ">", "", -1)
	DF = strings.Replace(DF, "|", "", -1)
	return DF

}
