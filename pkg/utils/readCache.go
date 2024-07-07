package utils

import (
	"encoding/json"
	"os"
)

var path = "dao/cache.json"

func Readcache() (oldCheck map[string][]int) {
	//var oldCheck = map[string][]int{}

	file, _ := os.Open(path)
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&oldCheck); err != nil {
		return
	}
	return
}

func SaveCache(oldCheck map[string][]int) {
	data, _ := json.Marshal(oldCheck)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}
