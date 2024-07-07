package tutorials

import (
	"fyne.io/fyne/v2"
	"gpixivImageDownload/tutorials/home"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"首页": {"Common",
			"...",
			home.canvasCommon,
			true,
		},
		"ranks": {"Ranks",
			"...",
			canvasRanks,
			true,
		},
		"authors": {"Authors",
			"...",
			canvasAuthors,
			true,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"Common", "rank", "author"},
		//"collections": {"list", "table", "tree", "gridwrap"},
		//"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split"},
		//"widgets":     {"accordion", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
