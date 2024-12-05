package model

import (
	"sync"
	"sync/atomic"
)

//var Rankser Ranks

type Ranks struct {
	sync.Once
	User         string
	Day          []string
	Week         []string
	Month        []string
	OldCheck     map[string][]int
	Tops         int
	DownloadPath string
	State        atomic.Int64 // 表示有多少个下载进程
	List         [3]int       // 0,1  0：停止，1：下载中
}

func (r *Ranks) Content() string {
	s := ""
	if len(r.Day) != 0 {
		s += "Day=["
		for _, i := range r.Day {
			s += i
			s += ","
		}
		s = s[:len(s)-1] + "];"
	}

	if len(r.Week) != 0 {
		s += "Week=["
		for _, i := range r.Week {
			s += i
			s += ","
		}
		s = s[:len(s)-1] + "];"
	}

	if len(r.Month) != 0 {
		s += "Month=["
		for _, i := range r.Month {
			s += i
			s += ","
		}
		s = s[:len(s)-1] + "];"
	}
	if len(s) > 0 {
		s = s[:len(s)-1]
	}

	return s
}
