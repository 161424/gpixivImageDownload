package home

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gpixivImageDownload/conf"
	log2 "gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"log"
	"strings"
	"time"
)

var CommonVar = &model.Common{}
var l = log2.Logger

func init() {
	v := conf.Conf.Sub("CmtOption")
	e := v.Unmarshal(CommonVar)
	if e != nil {
		fmt.Println(e)
	}
	CommonVar.DownloadPath = v.GetString("DownloadPath")
	//fmt.Println(CommonVar)
}

func CanvasCommon(win fyne.Window) fyne.CanvasObject {

	comm2update := make(chan string, 4)

	ckList := widget.NewMultiLineEntry()
	ckList.SetPlaceHolder("使用回车区分不同ck")
	ckList.Validator = validation.NewRegexp(`^.+$`, "ck无效")
	ckList.CursorColumn = 30

	R18 := widget.NewCheck("", func(bool) {})
	mThread := widget.NewCheck("", func(bool) {})
	skipIllus := widget.NewCheck("", func(bool) {})
	skipUgoira := widget.NewCheck("", func(bool) {})
	skipManga := widget.NewCheck("", func(bool) {})
	skipIllus.Checked = true
	skipUgoira.Checked = true
	skipManga.Checked = true

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Cookies", Widget: ckList, HintText: "Pixiv的cookies"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			CommonVar.Ck = strings.Split(ckList.Text, `\n`)
			CommonVar.R18 = R18.Checked
			CommonVar.MThread = mThread.Checked
			CommonVar.SkipIllus = skipIllus.Checked
			CommonVar.SkipUgoira = skipUgoira.Checked
			CommonVar.SkipManga = skipManga.Checked
			comm2update <- "修改common配置"

		},
	}
	form.SubmitText = "确认"
	form.CancelText = "撤销"

	form.Append("R18", R18)
	form.Append("多线程", mThread)
	form.Append("跳过插图", skipIllus)
	form.Append("跳过动图", skipUgoira)
	form.Append("跳过漫画", skipManga)

	downloadPath := widget.NewLabelWithData(binding.NewString())
	openFolder := widget.NewButton("设置根下载目录：", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}
			downloadPath.SetText(list.String())
			CommonVar.DownloadPath = list.String()
			downloadPath.Refresh()
			CommonVar.Save()
			comm2update <- "修改common地址"

			//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			//dialog.ShowInformation("Folder Open", out, win)
		}, win)
	})

	downloadPathWig := container.NewHBox(openFolder, downloadPath)

	CommonVar.Do(func() {
		go func() {
			t := time.NewTicker(time.Minute)
			for {
				select {
				case s := <-comm2update:
					s1 := CommonVar.Save()
					l.Send(4, s+s1, 2)
				case <-t.C:
					l.Send(4, fmt.Sprintf("home heart", &CommonVar), 1)
				}
			}
		}()
	})

	return container.NewVBox(form, downloadPathWig)

}
