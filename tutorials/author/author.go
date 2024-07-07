package author

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/author"
	"gpixivImageDownload/tutorials/home"
	"log"
	"strconv"
	"strings"
)

var AuthorVar *model.Author

func CanvasAuthor(win fyne.Window) fyne.CanvasObject {
	loadSync, cancel := context.WithCancel(context.Background())

	authorName := widget.NewMultiLineEntry()
	authorName.CursorRow = 2
	authorName.SetPlaceHolder("竜崎いち")

	authorId := widget.NewMultiLineEntry()
	authorId.CursorRow = 2
	authorId.SetPlaceHolder("563034")

	Tagsinp := widget.NewMultiLineEntry()
	Tagsinp.CursorRow = 2
	Tagsinp.SetPlaceHolder("#宵崎奏")

	authorDwProgress := widget.NewProgressBar()
	authorDwProgress.DisableColor = true

	dwpip := make(chan float64, 0)
	//dwValue := 0.0
	//authorDwProgress.Bind(binding.BindFloat(&dwValue))

	f := "50"
	//data := binding.BindInt(&f)
	topEntry := widget.NewEntryWithData(binding.BindString(&f))
	topEntry.Validator = func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if (n < 0) || (n > 1000) {
			return errors.New("超出范围")
		}
		return nil
	}

	rankdownloadPath := widget.NewLabelWithData(binding.NewString())
	openFolder := widget.NewButton("选择保存文件地址：", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				log.Println("Cancelled")
				return
			}
			_rankdownloadpath := strings.Split(list.String(), "file://")[1]
			fmt.Println(_rankdownloadpath)
			rankdownloadPath.SetText(_rankdownloadpath)
			AuthorVar.DownloadPath = _rankdownloadpath
			//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			//dialog.ShowInformation("Folder Open", out, win)
		}, win)
	})

	path := container.New(layout.NewHBoxLayout(), openFolder, rankdownloadPath)

	Tags := widget.NewRadioGroup([]string{"作者下的Tag", "仅Tag"}, func(string) {})
	Tags.Horizontal = true
	authorForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "作者名字", Widget: authorName, HintText: "使用回车间隔"},
			{Text: "作者id", Widget: authorId, HintText: "使用,间隔"},
			{Text: "Tags选择", Widget: Tags},
			{Text: "Tags输入", Widget: Tagsinp, HintText: "使用,间隔"},
			{Text: "下载个数top=", Widget: topEntry, HintText: "0~100"},
			{Text: "", Widget: path},
		},
		OnCancel: func() {
			cancel()
		},
	}
	//authorForm.Append("Tags选择", Tags)
	//authorForm.Append("Tags输入", Tagsinp)
	result := ""
	resultShow := widget.NewMultiLineEntry()
	resultShow.Bind(binding.BindString(&result))
	resultShow.CursorRow = 7

	authorForm.Append("下载进度", authorDwProgress)

	authorForm.OnSubmit = func() {
		var err error
		authorDwProgress.DisableColor = false
		if Tags.Selected == "仅Tag" {
			AuthorVar.OnlyTag = true
		}

		AuthorVar.DwTop, err = strconv.Atoi(topEntry.Text)

		if (checkDog(AuthorVar, authorName.Text, authorId.Text, Tagsinp.Text) && AuthorVar.OnlyTag == false) || err != nil {
			w := fyne.CurrentApp().NewWindow("错误")
			w.SetContent(widget.NewLabel("输入信息错误！"))
			w.Show()
			return
		}

		go author.DownLoadAuth(loadSync, home.CommonVar, AuthorVar, dwpip, &result)
		for {
			select {
			case v, _ := <-dwpip:
				authorDwProgress.SetValue(v)
			case <-loadSync.Done():
				return
			}
		}

	}

	return container.NewVBox(authorForm, resultShow)
}

func checkDog(iptMsg *model.Author, an, ai, tg string) bool {
	iptMsg.AuthorName = strings.Split(an, "\n")
	iptMsg.Tags = strings.Split(ai, "#")
	iptMsg.AuthorId = strings.Split(tg, ",")
	return true
}
