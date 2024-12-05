package model

import "sync"

type Author struct {
	sync.Once
	AuthorName   []string
	AuthorId     []string
	DwTop        int
	OnlyTag      bool
	Tags         string
	BookMark     bool
	OffSet       int
	Profile      bool
	DownloadPath string
}
