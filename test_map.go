package main

import (
	"fmt"
	"time"
)

func main() {
	//myApp := app.New()
	//myWindow := myApp.NewWindow("Button Widget")
	//
	//content := widget.NewButton("click me", func() {
	//	log.Println("tapped")
	//})
	//
	////content := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
	////	log.Println("tapped home")
	////})
	//
	//myWindow.SetContent(content)
	//myWindow.ShowAndRun()
	w := make(chan int, 0)
	<-w
	e := time.Date(2023, 3, 1-29, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	fmt.Println(e)
}
