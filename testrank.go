package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gpixivImageDownload/model"
	themes "gpixivImageDownload/theme"
	"log"
	"strconv"
	"strings"
	"time"
)

var RanksVar = &model.Ranks{}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&themes.MyTheme{})
	myWindow := myApp.NewWindow("Form Layout")
	myWindow.SetContent(canvasRanks(myWindow))
	myWindow.ShowAndRun()
}

func canvasRanks(win fyne.Window) fyne.CanvasObject {
	//
	oldCheck := map[string]*[]int{"2024-5": &[]int{1, 2, 3, 4, 5, -19}, "2024-7": &[]int{1, 2, 3, 4, 5}}
	newCheck := map[string]*[]int{}

	caleButton := widget.NewButton("请选择下载日期", func() {
		Calendars(fyne.CurrentApp().NewWindow("请选择下载日期"), oldCheck, newCheck)
	})

	//var downloadprogress = map[string]float64{"data": 0, "week": 0, "month": 0}
	dataprogress := widget.NewProgressBar()
	weekprogress := widget.NewProgressBar()
	monthprogress := widget.NewProgressBar()
	dataprogress.DisableColor = true
	weekprogress.DisableColor = true
	monthprogress.DisableColor = true

	fmt.Println("4")
	wf := 0.0

	go func() {
		for i := 0; i < 10; i++ {
			wf += 0.1
			dataprogress.SetValue(wf)
			time.Sleep(time.Second)
			if i == 9 {
				i = 0
				wf = 0
			}
		}

	}()

	datas := ""
	binddatas := binding.BindString(&datas)
	dataEntry := widget.NewEntryWithData(binddatas)

	dataLabel1 := widget.NewLabel("日排行榜")
	dataCheck1 := widget.NewCheck("下载", func(b bool) {
		if b {
			dataEntry.Enable()
			datas = ""
			for key, _ := range newCheck {
				for _, j := range *newCheck[key] {
					if j <= 0 {
						continue
					}
					s := fmt.Sprintf("%s-%d", key, j) // 2024-5-10
					RanksVar.Day = append(RanksVar.Day, s)
					datas += s
					datas += ","
				}

			}
			if len(datas) == 0 {
				return
			}
			datas = datas[:len(datas)-1]
			binddatas.Set(datas)
		} else {
			dataEntry.Disable()
		}
	})
	datarank := container.NewBorder(nil, nil, container.NewHBox(dataLabel1, dataCheck1), nil, container.NewGridWithColumns(2, dataEntry, dataprogress))

	weeks := ""
	bindweek := binding.BindString(&weeks)
	weekEntry := widget.NewEntryWithData(bindweek)

	weekLabel1 := widget.NewLabel("周排行榜")
	weekCheck1 := widget.NewCheck("下载", func(b bool) {
		if b {
			weekEntry.Enable()
			weeks = ""
			for key, _ := range newCheck {
				for _, j := range *newCheck[key] {
					if j > 0 || j < -100 {
						continue
					}
					y := strings.Split(key, "-")[0]
					s := fmt.Sprintf("%s第%d周", y, -j)
					ys := fmt.Sprintf("%s%d", y, j)
					RanksVar.Week = append(RanksVar.Week, ys) // 2024-19
					weeks += s
					weeks += ","
				}
			}
			if len(weeks) == 0 {
				return
			}
			weeks = weeks[:len(weeks)-1]
			bindweek.Set(weeks)
		} else {
			weekEntry.Disable()
		}
	})
	weekrank := container.NewBorder(nil, nil, container.NewHBox(weekLabel1, weekCheck1), nil, container.NewGridWithColumns(2, weekEntry, weekprogress))

	months := ""
	bindmonths := binding.BindString(&months)
	monthEntry := widget.NewEntryWithData(bindmonths)

	monthLabel1 := widget.NewLabel("月排行榜")
	monthCheck1 := widget.NewCheck("下载", func(b bool) {
		if b {
			monthEntry.Enable()
			months = ""
			for key, _ := range newCheck {
				for _, j := range *newCheck[key] {
					if j > -100 {
						continue
					}
					ym := strings.Split(key, "-")
					s := fmt.Sprintf("%s第%s月", ym[0], ym[1])
					RanksVar.Month = append(RanksVar.Month, key) // 2024-5
					months += s
					months += ","
				}
			}
			if len(months) == 0 {
				return
			}
			months = months[:len(months)-1]
			bindmonths.Set(months)
		} else {
			monthEntry.Disable()
		}

	})
	monthrank := container.NewBorder(nil, nil, container.NewHBox(monthLabel1, monthCheck1), nil, container.NewGridWithColumns(2, monthEntry, monthprogress))

	f := 50
	data := binding.BindInt(&f)
	label := widget.NewLabel("下载个数top=")
	entry := widget.NewEntryWithData(binding.IntToString(data))
	entry.Validator = func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if (n < 0) || (n > 100) {
			return errors.New("超出范围")
		}
		return nil
	}
	text := canvas.NewText("0~100", theme.PlaceHolderColor())
	text.TextSize = theme.CaptionTextSize()
	downloadtop := container.NewHBox(label, container.NewVBox(entry, text))

	rankdownloadpath := binding.BindString(nil)
	rankdownloadPath := widget.NewLabelWithData(rankdownloadpath)
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
			rankdownloadpath.Set(_rankdownloadpath)
			RanksVar.DownloadPath = _rankdownloadpath
			//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			//dialog.ShowInformation("Folder Open", out, win)
		}, win)
	})

	path := container.New(layout.NewHBoxLayout(), openFolder, rankdownloadPath)

	downloadbutton := widget.NewButton("下载", func() {
		if dataCheck1.Checked {
			dataprogress.DisableColor = false
		}
		if weekCheck1.Checked {
			weekprogress.DisableColor = false
		}
		if monthCheck1.Checked {
			monthprogress.DisableColor = false
		}

	})
	downloadbutton.Importance = widget.HighImportance
	cancelbutton := widget.NewButton("取消", func() {

	})

	bt := container.NewHBox(cancelbutton, downloadbutton)
	btx := container.NewBorder(nil, nil, nil, bt, bt)

	return container.NewVBox(caleButton, datarank, weekrank, monthrank, downloadtop, path, btx)

}

func Calendars(_ fyne.Window, oldCheck, newCheck map[string]*[]int) {
	var ranks = map[string]bool{"date": false, "week": false, "month": false}
	myWindow := fyne.CurrentApp().NewWindow("请选择下载日期")
	years := time.Now().Year()
	months := int(time.Now().Month())
	*(oldCheck["2024-7"]) = append(*(oldCheck["2024-7"]), 23)

	// 当月
	binderq := binding.BindString(nil)
	s := fmt.Sprintf("%d-%d", years, months)
	binderq.Set(s)
	newCheck[s] = new([]int)
	var bts []fyne.CanvasObject
	bts = printCalendarss(years, months, oldCheck[s], newCheck[s], &ranks)
	btss := container.NewGridWithColumns(8, bts...)

	contentCache := map[string]*fyne.Container{}
	content := &fyne.Container{}
	*content = *btss
	contentCache[s] = btss

	bt1 := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		years -= 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)

		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = new([]int)
			if oldCheck[s] == nil {
				oldCheck[s] = new([]int)
			}
			bts = printCalendarss(years, months, oldCheck[s], newCheck[s], &ranks)
			btss = container.NewGridWithColumns(8, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})

	bt2 := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		months -= 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)

		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = new([]int)
			if oldCheck[s] == nil {
				oldCheck[s] = new([]int)
			}
			bts = printCalendarss(years, months, oldCheck[s], newCheck[s], &ranks)
			btss = container.NewGridWithColumns(8, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})

	rq := widget.NewLabelWithData(binderq)
	rq.Alignment = fyne.TextAlignCenter

	bt3 := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		months += 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)

		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = new([]int)
			if oldCheck[s] == nil {
				oldCheck[s] = new([]int)
			}
			bts = printCalendarss(years, months, oldCheck[s], newCheck[s], &ranks)
			btss = container.NewGridWithColumns(8, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})
	bt4 := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		years += 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)

		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = new([]int)
			if oldCheck[s] == nil {
				oldCheck[s] = new([]int)
			}
			bts = printCalendarss(years, months, oldCheck[s], newCheck[s], &ranks)
			btss = container.NewGridWithColumns(8, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}

	})

	H := container.NewGridWithColumns(5, bt1, bt2, rq, bt3, bt4)

	daysWeek := container.NewGridWithColumns(8,
		widget.NewLabel(" 周 数 "),
		widget.NewLabel(" 星期一 "),
		widget.NewLabel(" 星期二 "),
		widget.NewLabel(" 星期三 "),
		widget.NewLabel(" 星期四 "),
		widget.NewLabel(" 星期五 "),
		widget.NewLabel(" 星期六 "),
		widget.NewLabel(" 星期日 "),
	)

	daysRank := widget.NewCheck("日排行", func(b bool) {
		if b {
			ranks["date"] = true
		} else {
			ranks["date"] = false
		}
	})

	weeksRank := widget.NewCheck("周排行", func(b bool) {
		if b {
			ranks["week"] = true
		} else {
			ranks["week"] = false
		}
	})

	monthRank := widget.NewCheck("月排行", func(b bool) {
		if b {
			ranks["month"] = true
		} else {
			ranks["month"] = false
		}
	})

	cale := container.NewVBox(daysWeek, content)
	ranksCheck := container.NewHBox(daysRank, weeksRank, monthRank)
	cales := container.NewVBox(cale, ranksCheck)
	confirm := widget.NewButton("Confirm", func() {
		myWindow.Close()
	})
	myWindow.SetContent(container.New(layout.NewVBoxLayout(), H, cales, confirm))
	myWindow.Show()

}

func printCalendarss(year, month int, old, new *[]int, v *map[string]bool) []fyne.CanvasObject { //定义了一个名为printCalendar的函数，它接受两个整数参数，分别表示年份和月份。s
	nowYear, nowMonth, nowDay := time.Now().Date()
	re := []int{}
	days := make([]int, 42)                     // 创建一个长度为42的整数切片days。切片是Go语言中的动态数组，可以用于存储一系列整数。
	day := 1                                    //初始化一个变量day，将其赋值为 1。`day`将用于迭代日期。
	dayOfWeek := getDayOfWeekss(year, month, 1) //调用名为getDayOfWeek的函数，传递年份、月份和日期（初始化为 1）作为参数。这个函数将返回对应的星期几的值。
	maxDay := getMaxDayOfMonthss(year, month)   //调用名为getMaxDayOfMonth的函数，传递年份和月份作为参数。这个函数将返回该月的最大天数。
	for i := 0; i < maxDay; i++ {               //使用一个循环，从0开始迭代到`maxDay`减 1 的范围。循环体将执行多次，每次迭代对应一天。
		days[i+dayOfWeek] = day // 在循环体内，将当前日期`day`存储到切片`days`的指定索引位置。索引是通过计算星期几和循环迭代次数得到的。
		day++                   //在每次循环迭代后，将日期变量`day`加 1，以便在下一次迭代中处理下一天的日期。
	}
	//fmt.Println(days, dayOfWeek)
	//fmt.Println("      Sun Mon Tue Wed Thu Fri Sat") //输出一周七天的头部标题，用于对齐日期。
	weekNumStart := (time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).YearDay()+dayOfWeek)/7 + 1
	weekNumNow := (time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local).YearDay()+dayOfWeek)/7 + 1
	//weekNumEnd := (time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).YearDay() + dayOfWeek) / 7
	index := 1                // 初始化一个变量index，将其赋值为 0。index将用于迭代切片中的日期。
	for i := 1; i <= 6; i++ { //使用另一个循环，从1开始迭代到6的范围。循环体将执行多次，每次迭代对应一周中的一天。
		//fmt.Print(i, "   ") //- `fmt.Print(i, "   ")`：在循环体内，输出当前迭代的索引（表示星期几），并在后面添加三个空格。
		re = append(re, -1)
		for j := 0; j < 7; j++ { // 使用嵌套的循环，从0开始迭代到6的范围。循环体将执行多次，每次迭代对应一周中的一个位置。
			if index < len(days) && days[index] != 0 { //在循环体内，检查当前索引是否小于切片的长度，并且对应的日期值不为0。这是为了确保在输出日期时不越界。
				//fmt.Printf("%4d", days[index]) //如果条件满足，输出当前日期的值，使用格式化字符串`%4d`来指定输出格式为四位数的日期。
				re = append(re, days[index])
			} else {
				re = append(re, 0)
			}
			index++ //在每次循环迭代后，将索引变量`index`加 1，以便在下一次迭代中处理下一个日期。
		}
	}

	colorSeg := &widget.TextSegment{
		Style: widget.RichTextStyleStrong,
		Text:  "",
	}

	//fmt.Println(re)
	bts := []fyne.CanvasObject{}
	dayCheck := 1
	for _, i := range re {
		btt := &widget.ColorButton{}
		checked := false
		oldDownloadDay := false
		oldDownloadWeek := false
		oldDownloadMonth := false

		we := dayCheck
		if i == 0 {
			btt = widget.NewColorButton("", func() {})
			btt.TextStyle = colorSeg
		} else {
			bs := strconv.Itoa(i)
			if i == -1 {
				bs = fmt.Sprintf("第%d周", weekNumStart)
			}
			for _, k := range *old {
				if we == k {
					oldDownloadDay = true
					break
				}
				if -we == k {
					oldDownloadWeek = true
					break
				}
				if -100-we == k {
					oldDownloadMonth = true
					break
				}
			}

			dateTextColor := &widget.TextSegment{
				Style: widget.RichTextStyleStrong,
				Text:  bs,
			}
			weekButtonColor := theme.ButtonColor()
			monthBorderColor := theme.ButtonColor()
			btt = widget.NewColorButton(bs, func() {
				checked = !checked
				fmt.Println(v, checked, oldDownloadDay, oldDownloadWeek, oldDownloadMonth, new)
				if (*v)["date"] {
					if checked {
						dateTextColor.Style.ColorNRGBA = themes.ColorBluAA
						*new = append(*new, we)
						//fmt.Printf("%p,%p,%p,%p,%p", &dayCheck, &we, new, &checked, &orangeButton)
					} else {
						if oldDownloadDay {
							dateTextColor.Style.ColorNRGBA = themes.ColorBluHA
						} else {
							dateTextColor.Style.ColorNRGBA = nil
						}
						*new = DeleteSlices(*new, we)
					}
					btt.TextStyle = dateTextColor
					fmt.Println(dateTextColor.Style.ColorNRGBA)
				}
				if (*v)["week"] {
					if checked {
						weekButtonColor = themes.ColorGreAA
						*new = append(*new, -we)
						//fmt.Println("?", checked)
					} else {
						if oldDownloadWeek {
							weekButtonColor = themes.ColorGreHA
						} else {
							weekButtonColor = theme.ButtonColor()
						}
						*new = DeleteSlices(*new, -we)
					}
					btt.ButtonColor = weekButtonColor
					fmt.Println(weekButtonColor)
				}
				if (*v)["month"] {
					if checked {
						monthBorderColor = themes.ColorRedAA
						*new = append(*new, -100-we)
					} else {
						if oldDownloadMonth {
							monthBorderColor = themes.ColorRedHA
						} else {
							monthBorderColor = theme.ButtonColor()
						}
						*new = DeleteSlices(*new, -100-we)
					}
					btt.BorderColor = monthBorderColor

				}
				btt.Refresh()
			})

			if oldDownloadDay {
				dateTextColor.Style.ColorNRGBA = themes.ColorBluHA
			} else {
				dateTextColor.Style.ColorNRGBA = nil
			}

			if oldDownloadWeek {
				weekButtonColor = themes.ColorOrgHA
			} else {
				weekButtonColor = theme.ButtonColor()
			}

			if oldDownloadMonth {
				monthBorderColor = themes.ColorRedHA
			} else {
				monthBorderColor = theme.ButtonColor()
			}

			btt.TextStyle = dateTextColor
			btt.ButtonColor = weekButtonColor
			btt.BorderColor = monthBorderColor

			if i == -1 {
				weekNumStart++
				btt.ButtonColor = themes.ColorBlack
			} else {
				dayCheck += 1
			}
		}

		if time.Date(year, time.Month(month), dayCheck, 0, 0, 0, 0, time.Local).Sub(
			time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local)) > 0 {
			//fmt.Println("btt disabled")
			//fmt.Println(time.Date(year, time.Month(month), dayCheck, 0, 0, 0, 0, time.Local))
			//fmt.Println(time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local))
			btt.Disable()
		}

		if (year == nowYear) && i == -1 {
			if weekNumNow <= weekNumStart {
				btt.Disable()
			}
		}

		//fmt.Println(time.Date(year, time.Month(month), dayCheck-1, 0, 0, 0, 0, time.Local),
		//	time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local), time.Date(year, time.Month(month), dayCheck-1, 0, 0, 0, 0, time.Local).Sub(
		//		time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local)) > 0)
		bts = append(bts, btt)

	}

	return bts
}

func DeleteSlices(a []int, elem int) []int {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

// 这是一个打印日历的函数。
// 它首先创建了一个长度为42的整型切片days，用于存储每个月份中的日期。
// 然后，使用getDayOfWeek函数获取给定年份和月份的第一天是星期几，
// 并使用getMaxDayOfMonth函数获取给定年份和月份的最大天数。
// 接下来，使用两个循环填充days切片。
// 第一个循环从1到最大天数，将日期依次放入正确的位置。第二个循环打印日历的表头以及日期。
func getDayOfWeekss(year, month, day int) int { //定义了一个名为getDayOfWeek的函数，它接受三个整数参数，分别表示年份、月份和日期。
	if month < 3 { //检查月份是否小于 3。这是为了处理月份在 1 月或 2 月的情况。
		month += 12 //如果月份小于 3，将其加上 12，以便将日期转换为对应的星座月份。
		year--      //同时，将年份减 1，以调整到正确的年份。
	}
	century := year / 100                                                          //计算年份的世纪部分，即将年份除以 100。
	year %= 100                                                                    //将年份的余数保留下来，用于后续计算。
	week := (day + 13*(month+1)/5 + year + year/4 + century/4 - 2*century - 1) % 7 //用于确定给定日期对应的星期几。它考虑了月份、日期、年份和世纪的各种组合和调整。

	//return (week + 6) % 7 //最后，将计算得到的星期几加上6，然后再对7取余数，以将结果调整为0到6的范围，对应星期日到星期六。
	return week
}

// 这是一个计算给定日期是星期几的函数。它使用了蔡勒公式（Zeller's Congruence）来计算星期几。该公式根据年份、月份和日期计算出一个数字，表示星期几（0代表星期日，1代表星期一，以此类推）。
func getMaxDayOfMonthss(year, month int) int { //定义了一个名为`getMaxDayOfMonth`的函数，它接受两个整数参数，分别表示年份和月份。
	if month == 2 { //检查月份是否为 2。
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 { //条件判断，用于确定是否为闰年。如果年份能被4整除但不能被100整除，或者能被400整除，那么就是闰年。
			return 29 //在闰年的情况下，2月份的最大天数为29。
		} else { //否则，如果不是闰年。
			return 28 //2月份的最大天数为28。
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 { //检查其他月份是否为4、6、9 或11。
		return 30 //这些月份的最大天数为30。
	} else { //否则，如果不是上述月份。
		return 31 //其他月份的最大天数为31。
	}
}
