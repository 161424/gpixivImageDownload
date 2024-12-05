package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	themes "gpixivImageDownload/theme"
	"gpixivImageDownload/tutorials/author"
	"gpixivImageDownload/tutorials/home"
	"gpixivImageDownload/tutorials/rank"
)

const preferenceCurrentTutorial = "currentTutorial"

func main() {

	menu := []string{"首页", "排行", "作者"}
	defaultTheme := &themes.MyTheme{}
	a := app.NewWithID("io.fyne.demo")
	a.SetIcon(data.FyneLogo)
	a.Settings().SetTheme(defaultTheme)
	mainWin := a.NewWindow("go-pixiv")

	menuWin := widget.NewList(
		func() int {
			return len(menu)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.HomeIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id == 0 {
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.HomeIcon())
			} else if id == 1 {
				rankIcon, _ := fyne.LoadResourceFromPath("./theme/icon/排行.svg")
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(rankIcon)
			} else if id == 2 {
				authorIcon, _ := fyne.LoadResourceFromPath("./theme/icon/作者.svg")
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(authorIcon)
			}
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(menu[id])
		},
	)

	content := container.NewStack()

	menuWin.OnSelected = func(id widget.ListItemID) {

		if id == 0 {
			content.Objects = []fyne.CanvasObject{home.CanvasCommon(mainWin)}
		} else if id == 1 {
			content.Objects = []fyne.CanvasObject{rank.CanvasRanks(mainWin)}
		} else if id == 2 {
			content.Objects = []fyne.CanvasObject{author.CanvasAuthor(mainWin)}
		}
		content.Refresh()
	}
	menuWin.Select(0)
	//themes := container.NewGridWithColumns(2,
	//	widget.NewButton("Dark", func() {
	//		a.Settings().SetTheme(defaultTheme)
	//	}),
	//	widget.NewButton("Light", func() {
	//		a.Settings().SetTheme(defaultTheme)
	//		a.Settings().Theme()
	//	}),
	//)

	minBorder := container.NewCenter()
	minBorder.Resize(fyne.NewSize(1, 1))

	newMenu := container.NewBorder(minBorder, container.NewVBox(minBorder), minBorder, minBorder, menuWin)

	contentView := container.NewHSplit(newMenu, container.NewBorder(minBorder, minBorder, minBorder, minBorder, content))
	contentView.Offset = 0.2
	mainWin.SetContent(contentView)

	mainWin.Resize(fyne.NewSize(640, 460))
	mainWin.ShowAndRun()
}
