package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	themes "gpixivImageDownload/theme"
	"image/color"
)

func main() { //定义了一个名为main的函数，这是Go程序的入口函数。程序从这里开始执行。

	myApp := app.New()
	myApp.Settings().SetTheme(&themes.MyTheme{})
	myWindow := myApp.NewWindow("Form Layout")
	tx := canvas.NewText("123", color.RGBA{0xff, 0xa5, 0x00, 0x7f})
	//ce := false
	//openFolder := widget.NewColorButton("        ", func() {
	//	ce = !ce
	//	if ce {
	//		tx.Color = color.RGBA{0xff, 0x0, 0x0, 0xff}
	//	} else {
	//		tx.Color = color.RGBA{0xff, 0xa5, 0x00, 0x7f}
	//	}
	//	tx.Refresh()
	//
	//})

	//openFolder.TextStyle = &widget.TextSegment{
	//	Style: widget.RichTextStyle{
	//		Alignment: fyne.TextAlignLeading,
	//		ColorName: theme.ColorNameHyperlinkssss,
	//		Inline:    true,
	//		TextStyle: fyne.TextStyle{},
	//	},
	//	Text: "AAABBB",
	//}

	txs := widget.NewRichText([]widget.RichTextSegment{&widget.TextSegment{
		Style: widget.RichTextStyle{
			Alignment: fyne.TextAlignLeading,
			//ColorName: theme.ColorNameHyperlink,
			Inline:     true,
			TextStyle:  fyne.TextStyle{},
			ColorNRGBA: themes.ColorBluAA,
		},
		Text: "AAABBB",
	}}...)

	txss := widget.NewRichTextWithText("AAABBB")
	txss.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameHyperlink

	cks := false
	txsss := widget.NewColorButton("", func() {

	})
	//ww := binding.NewString()
	txsss.TextStyle = &widget.TextSegment{
		Style: widget.RichTextStyle{
			Alignment: fyne.TextAlignLeading,
			//ColorName: theme.ColorNameHyperlink,
			ColorNRGBA: themes.ColorBluAA,
			Inline:     true,
			TextStyle:  fyne.TextStyle{},
		},
		Text: "AAABBB",
	}

	txsss.OnTapped = func() {
		cks = !cks
		if cks {
			fmt.Println(1)
			txsss.TextStyle.Style.ColorName = theme.ColorNameHover
			txsss.BorderColor = color.RGBA{0x6f, 0x0, 0x0, 0x00}
		} else {
			fmt.Println(2)
			txsss.TextStyle.Style.ColorName = theme.ColorNameHyperlink
			txsss.BorderColor = color.RGBA{0xff, 0x0e, 0xa0, 0xff}
		}
		fmt.Println(txsss.TextStyle.Style.ColorName)
		txsss.Refresh()
	}
	//text := widget.NewRichTextWithText("123123123")

	myWindow.SetContent(container.NewVBox(container.NewCenter(tx), txs, txsss))
	myWindow.ShowAndRun()
}
