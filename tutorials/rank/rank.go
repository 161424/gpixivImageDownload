package rank

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/rank"
	"gpixivImageDownload/pkg/utils"
	"gpixivImageDownload/tutorials/home"
	"log"
	"strconv"
	"strings"
	"time"
)

var RanksVar = &model.Ranks{}

func CanvasRanks(win fyne.Window, oldCheck map[string][]int) fyne.CanvasObject {
	//
	//oldCheck := map[string]*[]int{"2024-5": &[]int{1, 2, 3, 4, 5, }, "2024": &[]int{1, 2, 5, 4, 7}}
	newCheck := map[string][]int{}

	caleButton := widget.NewButton("请选择下载日期", func() {
		Calendars(fyne.CurrentApp().NewWindow("请选择下载日期"), oldCheck, newCheck)
	})
	//fmt.Println(3, oldCheck, newCheck)
	dayprogress := widget.NewProgressBar()
	weekprogress := widget.NewProgressBar()
	monthprogress := widget.NewProgressBar()
	dayprogress.DisableColor = true
	weekprogress.DisableColor = true
	monthprogress.DisableColor = true
	dayChan := make(chan float64, 0)
	weekChan := make(chan float64, 0)
	monthChan := make(chan float64, 0)

	//wf := 0.0
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		wf += 0.1
	//		dayprogress.SetValue(wf)
	//		time.Sleep(time.Second)
	//		if i == 9 {
	//			i = 0
	//			wf = 0
	//		}
	//	}
	//
	//}()

	days := ""
	binddays := binding.BindString(&days)
	dayEntry := widget.NewEntryWithData(binddays)

	dayLabel1 := widget.NewLabel("日排行榜")
	dayCheck1 := widget.NewCheck("下载", func(b bool) {
		fmt.Println("day", newCheck, oldCheck, b)
		if b {
			dayEntry.Enable()
			days = ""
			for ym, _ := range newCheck {
				for ds, _ := range newCheck[ym] {
					y := strings.Split("-", ym)
					_ds := fmt.Sprintf("%s-%s-%d", y[0], y[1], ds) // 2024-5-10
					RanksVar.Day = append(RanksVar.Day, _ds)
					days += _ds
					days += ","

				}
			}
			if len(days) == 0 {
				return
			}
			days = days[:len(days)-1]
			binddays.Set(days)
		} else {
			dayEntry.Disable()
		}
	})
	dayrank := container.NewBorder(nil, nil, container.NewHBox(dayLabel1, dayCheck1), nil, container.NewGridWithColumns(2, dayEntry, dayprogress))

	weeks := ""
	bindweek := binding.BindString(&weeks)
	weekEntry := widget.NewEntryWithData(bindweek)

	weekLabel1 := widget.NewLabel("周排行榜")
	weekCheck1 := widget.NewCheck("下载", func(b bool) {
		if b {
			weekEntry.Enable()
			weeks = ""

			for ym, _ := range newCheck {
				for ds, _ := range newCheck[ym] {
					y := strings.Split("-", ym)
					if ds > 0 && ds < -100 {
						continue
					}
					ed := -ds
					fy, fm, fd := munix(y[0], y[1], ed, "week")
					ws := fmt.Sprintf("%d-%d-%d ~ %s-%s-%d", fy, fm, fd, y[0], y[1], ed) // 2024-5-10
					RanksVar.Week = append(RanksVar.Week, fmt.Sprintf("%s-%s-%d", y[0], y[1], ed))
					weeks += ws
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

			for ym, _ := range newCheck {
				for ds, _ := range newCheck[ym] {
					y := strings.Split("-", ym)
					if ds > -100 {
						continue
					}
					ed := -ds % 100
					fy, fm, fd := munix(y[0], y[1], ed, "month")
					ms := fmt.Sprintf("%d-%d-%d ~ %s-%s-%d", fy, fm, fd, y[0], y[1], ed) // 2024-5-10
					RanksVar.Month = append(RanksVar.Month, fmt.Sprintf("%s-%s-%d", y[0], y[1], ed))
					months += ms
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
	topEntry := widget.NewEntryWithData(binding.IntToString(data))
	topEntry.Validator = func(s string) error {
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
	downloadtop := container.NewHBox(label, container.NewVBox(topEntry, text))

	//rankdownloadpath := binding.BindString(nil)
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
			//fmt.Println(_rankdownloadpath)
			rankdownloadPath.SetText(_rankdownloadpath)
			RanksVar.DownloadPath = _rankdownloadpath
			//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			//dialog.ShowInformation("Folder Open", out, win)
		}, win)
	})

	path := container.New(layout.NewHBoxLayout(), openFolder, rankdownloadPath)

	loadSync, cancel := context.WithCancel(context.Background())

	downloadbutton := widget.NewButton("下载", func() {
		fmt.Println(2, home.CommonVar, RanksVar, dayChan)
		if dayCheck1.Checked {

			go rank.DownloadRank(loadSync, "day", home.CommonVar, RanksVar, dayChan)
			dayprogress.DisableColor = false
		}
		if weekCheck1.Checked {
			go rank.DownloadRank(loadSync, "week", home.CommonVar, RanksVar, weekChan)
			weekprogress.DisableColor = false
		}
		if monthCheck1.Checked {
			go rank.DownloadRank(loadSync, "month", home.CommonVar, RanksVar, monthChan)
			monthprogress.DisableColor = false
		}

		go func() {
			for {
				select {
				case v := <-dayChan:
					dayprogress.SetValue(v)
					if v == 1 {
						for _, d := range RanksVar.Day {
							_d := strings.SplitN(d, "-", 1)
							key, _value := _d[0], _d[1]
							value, _ := strconv.Atoi(_value)
							if _, ok := oldCheck[key]; ok {
								oldCheck[key] = append(oldCheck[key], value)
							} else {
								oldCheck[key] = []int{value}
							}
						}
						utils.SaveCache(oldCheck)
					}
				case v := <-weekChan:
					dayprogress.SetValue(v)
					if v == 1 {
						for _, w := range RanksVar.Week {
							_w := strings.SplitN(w, "-", 1)
							key, _value := _w[0], _w[1]
							value, _ := strconv.Atoi(_value)
							if _, ok := oldCheck[key]; ok {
								oldCheck[key] = append(oldCheck[key], value)
							} else {
								oldCheck[key] = []int{value}
							}
						}
						utils.SaveCache(oldCheck)
					}
				case v := <-monthChan:
					dayprogress.SetValue(v)
					if v == 1 {
						for _, m := range RanksVar.Month {
							_m := strings.SplitN(m, "-", 1)
							key, _value := _m[0], _m[1]
							value, _ := strconv.Atoi(_value)
							if _, ok := oldCheck[key]; ok {
								oldCheck[key] = append(oldCheck[key], value)
							} else {
								oldCheck[key] = []int{value}
							}
						}
						utils.SaveCache(oldCheck)
					}
				case <-loadSync.Done():
					return

				}
			}

		}()

	})
	downloadbutton.Importance = widget.HighImportance
	cancelbutton := widget.NewButton("取消", func() {
		cancel()
	})

	bt := container.NewHBox(cancelbutton, downloadbutton)
	btx := container.NewBorder(nil, nil, nil, bt, bt)

	return container.NewVBox(caleButton, dayrank, weekrank, monthrank, downloadtop, path, btx)

}

func munix(ys, ms string, d int, gap string) (oy int, om time.Month, od int) {
	y, _ := strconv.Atoi(ys)
	m, _ := strconv.Atoi(ms)
	if gap == "week" {
		oy, om, od = time.Date(y, time.Month(m), d-6, 0, 0, 0, 0, time.Local).Date()
	} else if gap == "month" {
		oy, om, od = time.Date(y, time.Month(m), d-29, 0, 0, 0, 0, time.Local).Date()
	}

	return

}
