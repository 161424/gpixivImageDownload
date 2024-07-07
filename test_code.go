package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	themes "gpixivImageDownload/theme"
	"strconv"
	"time"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&themes.MyTheme{})
	myWindow := myApp.NewWindow("Form Layout")
	years := time.Now().Year()
	months := int(time.Now().Month())

	oldCheck := map[string][]int{}
	newCheck := map[string][]int{}

	oldCheck["2024-7"] = append(oldCheck["2024-7"], 23)

	binderq := binding.BindString(nil)
	s := fmt.Sprintf("%d-%d", years, months)
	binderq.Set(s)
	newCheck[s] = make([]int, 0)
	bts := printCalendars(years, months, oldCheck[s], newCheck[s])
	btss := container.NewGridWithColumns(7, bts...)
	fmt.Println(btss)
	contentCache := map[string]*fyne.Container{}
	content := &fyne.Container{}
	*content = *btss
	contentCache[s] = btss
	bt1 := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		years -= 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)
		fmt.Println(s)
		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = make([]int, 0)
			bts = printCalendars(years, months, oldCheck[s], newCheck[s])
			btss = container.NewGridWithColumns(7, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}

	})
	bt2 := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		months -= 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)
		fmt.Println(s)
		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = make([]int, 0)
			bts = printCalendars(years, months, oldCheck[s], newCheck[s])
			btss = container.NewGridWithColumns(7, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})

	rq1 := widget.NewLabelWithData(binderq)
	rq1.Alignment = fyne.TextAlignCenter

	bt3 := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		months += 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)
		fmt.Println(s)
		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects
		} else {
			newCheck[s] = make([]int, 0)
			bts = printCalendars(years, months, oldCheck[s], newCheck[s])
			btss = container.NewGridWithColumns(7, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})
	bt4 := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		years += 1
		s = fmt.Sprintf("%d-%d", years, months)
		binderq.Set(s)
		fmt.Println(s)
		if oldBtss, ok := contentCache[s]; ok {
			content.Layout = oldBtss.Layout
			content.Objects = oldBtss.Objects

		} else {
			newCheck[s] = make([]int, 0)
			bts = printCalendars(years, months, oldCheck[s], newCheck[s])
			btss = container.NewGridWithColumns(7, bts...)
			content.Layout = btss.Layout
			content.Objects = btss.Objects
			contentCache[s] = btss
		}
	})

	H := container.NewGridWithColumns(5, bt1, bt2, rq1, bt3, bt4)

	//day := widget.NewLabel("      Sun Mon Tue Wed Thu Fri Sat")

	daysWeek := container.NewGridWithColumns(7,
		widget.NewLabel(" Sun "),
		widget.NewLabel(" Mon  "),
		widget.NewLabel(" Tue  "),
		widget.NewLabel(" Wed  "),
		widget.NewLabel(" Thu  "),
		widget.NewLabel(" Fri  "),
		widget.NewLabel(" Sat  "),
	)
	fmt.Println("ccc")
	cale := container.NewVBox(daysWeek, content)
	confirm := widget.NewButton("Confirm", func() {
		cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", confirmCallback, myWindow)
		cnf.SetDismissText("Nah")
		cnf.SetConfirmText("Oh Yes!")
		cnf.Show()
	})
	myWindow.SetContent(container.New(layout.NewVBoxLayout(), H, cale, confirm))
	fmt.Println("aaa")
	myWindow.ShowAndRun()
	fmt.Println("bbb")
}

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func printCalendars(year, month int, old, new []int) []fyne.CanvasObject { //定义了一个名为printCalendar的函数，它接受两个整数参数，分别表示年份和月份。

	oldDays := []int{}
	if len(old) != 0 {
		oldDays = old
	}

	re := []int{}
	days := make([]int, 42)                    // 创建一个长度为42的整数切片days。切片是Go语言中的动态数组，可以用于存储一系列整数。
	day := 1                                   //初始化一个变量day，将其赋值为 1。`day`将用于迭代日期。
	dayOfWeek := getDayOfWeeks(year, month, 1) //调用名为getDayOfWeek的函数，传递年份、月份和日期（初始化为 1）作为参数。这个函数将返回对应的星期几的值。
	maxDay := getMaxDayOfMonths(year, month)   //调用名为getMaxDayOfMonth的函数，传递年份和月份作为参数。这个函数将返回该月的最大天数。
	for i := 0; i < maxDay; i++ {              //使用一个循环，从0开始迭代到`maxDay`减 1 的范围。循环体将执行多次，每次迭代对应一天。
		days[i+dayOfWeek] = day // 在循环体内，将当前日期`day`存储到切片`days`的指定索引位置。索引是通过计算星期几和循环迭代次数得到的。
		day++                   //在每次循环迭代后，将日期变量`day`加 1，以便在下一次迭代中处理下一天的日期。
	}
	//fmt.Println("      Sun Mon Tue Wed Thu Fri Sat") //输出一周七天的头部标题，用于对齐日期。

	index := 0                // 初始化一个变量index，将其赋值为 0。index将用于迭代切片中的日期。
	for i := 1; i <= 6; i++ { //使用另一个循环，从1开始迭代到6的范围。循环体将执行多次，每次迭代对应一周中的一天。
		//fmt.Print(i, "   ") //- `fmt.Print(i, "   ")`：在循环体内，输出当前迭代的索引（表示星期几），并在后面添加三个空格。
		re = append(re, -1)
		for j := 0; j < 7; j++ { // 使用嵌套的循环，从0开始迭代到6的范围。循环体将执行多次，每次迭代对应一周中的一个位置。
			if index < len(days) && days[index] != 0 { //在循环体内，检查当前索引是否小于切片的长度，并且对应的日期值不为0。这是为了确保在输出日期时不越界。
				//fmt.Printf("%4d", days[index]) //如果条件满足，输出当前日期的值，使用格式化字符串`%4d`来指定输出格式为四位数的日期。
				re = append(re, days[index])
			} else {
				//fmt.Print("    ")
				re = append(re, 0)
			}
			index++ //在每次循环迭代后，将索引变量`index`加 1，以便在下一次迭代中处理下一个日期。
		}
		//fmt.Println() //循环结束后，输出一个换行符，以便将不同星期的日期分隔开。
	}

	bts := []fyne.CanvasObject{}
	dayCheck := 1
	for _, i := range re {

		btt := &widget.ColorButton{}
		checked := false
		orangeButton := false

		if i == 0 {
			btt = widget.NewColorButton(" ", func() {})
		} else {
			for _, k := range oldDays {
				if dayCheck == k {
					orangeButton = true
					break
				}
			}

			btt = widget.NewColorButton(strconv.Itoa(i), func() {
				//btt.Importance = widget.HighImportance
				checked = !checked
				//fmt.Println("前", btt.Importance, btt.Color)
				if checked {
					//colors := color.RGBA{0x0, 0xae, 0xec, 0x7f}
					//btt.ButtonColor(colors)
					btt.Importance = widget.HighImportance
					new = append(new, dayCheck)
				} else {
					//theme.ButtonColor()
					fmt.Printf("???,%p", &orangeButton)
					//if orangeButton {
					//	colors := color.RGBA{0xff, 0xa5, 0x00, 0x7f}
					//	btt.ButtonColor(colors)
					//} else {
					//	btt.ButtonColor(theme.ButtonColor())
					//}
					new = DeleteSlice3(new, dayCheck)
				}
				//fmt.Println("后", btt.Importance, btt.Color)
				fmt.Println("button")
			})
			if orangeButton {
				//colorsr := color.RGBA{0xff, 0xa5, 0x00, 0x7f}
				//colors := color.RGBA{0x0, 0xae, 0xec, 0x7f}
				//btt.Color = colorsr
			}
			dayCheck += 1

		}
		bts = append(bts, btt)

	}

	return bts
}

func DeleteSlice3(a []int, elem int) []int {
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
func getDayOfWeeks(year, month, day int) int { //定义了一个名为getDayOfWeek的函数，它接受三个整数参数，分别表示年份、月份和日期。
	if month < 3 { //检查月份是否小于 3。这是为了处理月份在 1 月或 2 月的情况。
		month += 12 //如果月份小于 3，将其加上 12，以便将日期转换为对应的星座月份。
		year--      //同时，将年份减 1，以调整到正确的年份。
	}
	century := year / 100                                                                                     //计算年份的世纪部分，即将年份除以 100。
	year %= 100                                                                                               //将年份的余数保留下来，用于后续计算。
	week := (day + 2*month + 3*(month+1)/5 + year + year/4 - year/100 + year/400 + century/4 - 2*century) % 7 //用于确定给定日期对应的星期几。它考虑了月份、日期、年份和世纪的各种组合和调整。
	return (week + 6) % 7                                                                                     //最后，将计算得到的星期几加上6，然后再对7取余数，以将结果调整为0到6的范围，对应星期日到星期六。
}

// 这是一个计算给定日期是星期几的函数。它使用了蔡勒公式（Zeller's Congruence）来计算星期几。该公式根据年份、月份和日期计算出一个数字，表示星期几（0代表星期日，1代表星期一，以此类推）。
func getMaxDayOfMonths(year, month int) int { //定义了一个名为`getMaxDayOfMonth`的函数，它接受两个整数参数，分别表示年份和月份。
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

//func main() {
//	myApp := app.New()
//	myApp.Settings().SetTheme(&theme.MyTheme{})
//	myWindow := myApp.NewWindow("Form Layout")
//
//	ckList := widget.NewMultiLineEntry()
//	ckList.SetPlaceHolder("使用回车区分不同ck")
//	ckList.Validator = validation.NewRegexp(`\n`, "ck无效")
//
//	R18 := widget.NewCheck("???", func(in bool) { fmt.Println(in) })
//	mThread := widget.NewCheck("666", func(bool) {})
//
//	form := &widget.Form{
//		Items: []*widget.FormItem{
//			{Text: "Cookies", Widget: ckList, HintText: "Pixiv的cookies"},
//		},
//		OnCancel: func() {
//			fmt.Println("Cancelled")
//		},
//		OnSubmit: func() {
//			fmt.Println("Form submitted")
//
//			fmt.Println(ckList.Text)
//			fmt.Println(R18.Checked, R18.Text)
//			fmt.Println(mThread.Checked, mThread.Text)
//		},
//	}
//
//	form.SubmitText = "确认"
//	form.CancelText = "撤销"
//
//	form.Append("R18", R18)
//	form.Append("多线程", mThread)
//
//	grid := container.NewVBox(
//		form,
//		calendar,
//		layout.NewSpacer())
//	myWindow.SetContent(grid)
//	myWindow.ShowAndRun()
//}
