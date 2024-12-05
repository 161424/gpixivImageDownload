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
	log2 "gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/pkg/rank"
	"gpixivImageDownload/pkg/utils"
	"gpixivImageDownload/tutorials/home"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var RanksVar = &model.Ranks{}
var l = log2.Logger

func CanvasRanks(win fyne.Window) fyne.CanvasObject {

	ranks2update := make(chan string, 10)
	oldCheck := map[string][]int{}
	if RanksVar.OldCheck == nil {
		RanksVar.OldCheck = utils.Readcache()
		if RanksVar.OldCheck != nil {
			oldCheck = RanksVar.OldCheck
		}
	}

	//oldCheck := map[string]*[]int{"2024-5": &[]int{1, 2, 3, 4, 5, }, "2024": &[]int{1, 2, 5, 4, 7}}

	newCheck := map[string]*[]int{}

	caleButton := widget.NewButton("请选择下载日期", func() {
		Calendars(fyne.CurrentApp().NewWindow("请选择下载日期"), oldCheck, newCheck)
	})

	dayprogress := widget.NewProgressBar()
	weekprogress := widget.NewProgressBar()
	monthprogress := widget.NewProgressBar()
	dayprogress.DisableColor = true
	weekprogress.DisableColor = true
	monthprogress.DisableColor = true
	dayChan := make(chan float64)
	weekChan := make(chan float64)
	monthChan := make(chan float64)

	days := ""
	binddays := binding.BindString(&days)
	dayEntry := widget.NewEntryWithData(binddays)

	dayLabel1 := widget.NewLabel("日排行榜")
	dayCheck1 := widget.NewCheck("下载", func(b bool) {

		fmt.Println("day", newCheck, oldCheck, b)
		if b {
			dayEntry.Enable()
			days = ""
			for ym, _ := range newCheck { // key
				for _, ds := range *newCheck[ym] { // value[i]
					if ds <= 0 || ds > 32 {
						continue
					}
					y := strings.Split(ym, "-")
					//fmt.Println("ds", ds, ym, y, len(y))
					_ds := fmt.Sprintf("%s-%s-%d", y[0], y[1], ds) // 2024-5-10
					RanksVar.Day = append(RanksVar.Day, _ds)
					days += _ds + ","
				}
			}
			if len(days) == 0 {
				return
			}
			days = days[:len(days)-1]
			binddays.Set(days)
			ranks2update <- "已选择下载类型-日"
		} else {
			dayEntry.Disable()
			ranks2update <- "已取消下载类型-日"
		}
	})

	weeks := ""
	bindweek := binding.BindString(&weeks)
	weekEntry := widget.NewEntryWithData(bindweek)

	weekLabel1 := widget.NewLabel("周排行榜")
	weekCheck1 := widget.NewCheck("下载", func(b bool) {

		if b {
			weekEntry.Enable()
			weeks = ""
			for ym, _ := range newCheck {
				for _, ds := range *newCheck[ym] {
					y := strings.Split(ym, "-")
					if ds > 0 || ds < -100 {
						continue
					}
					ed := -ds
					fy, fm, fd := munix(y[0], y[1], ed, "week")
					ws := fmt.Sprintf("%d-%d-%d ~ %s-%s-%d", fy, fm, fd, y[0], y[1], ed) // 2024-5-10
					RanksVar.Week = append(RanksVar.Week, fmt.Sprintf("%s-%s-%d", y[0], y[1], ed))
					weeks += ws + ","
				}
			}
			if len(weeks) == 0 {
				return
			}
			weeks = weeks[:len(weeks)-1]
			bindweek.Set(weeks)
			ranks2update <- "已选择下载类型-周"
		} else {
			weekEntry.Disable()
			ranks2update <- "已取消下载类型-周"
		}
	})

	months := ""
	bindmonths := binding.BindString(&months)
	monthEntry := widget.NewEntryWithData(bindmonths)

	monthLabel1 := widget.NewLabel("月排行榜")
	monthCheck1 := widget.NewCheck("下载", func(b bool) {

		if b {
			monthEntry.Enable()
			months = ""

			for ym, _ := range newCheck {
				for _, ds := range *newCheck[ym] {
					y := strings.Split(ym, "-")
					if ds > -100 {
						continue
					}
					ed := -ds % 100
					fy, fm, fd := munix(y[0], y[1], ed, "month")
					ms := fmt.Sprintf("%d-%d-%d ~ %s-%s-%d", fy, fm, fd, y[0], y[1], ed) // 2024-5-10
					RanksVar.Month = append(RanksVar.Month, fmt.Sprintf("%s-%s-%d", y[0], y[1], ed))
					months += ms + ","
				}
			}
			if len(months) == 0 {
				return
			}
			months = months[:len(months)-1]
			bindmonths.Set(months)
			ranks2update <- "已选择下载类型-月"
		} else {
			monthEntry.Disable()
			ranks2update <- "已取消下载类型-月"
		}

	})

	dayrank := container.NewBorder(nil, nil, container.NewHBox(dayLabel1, dayCheck1), nil, container.NewGridWithColumns(2, dayEntry, dayprogress))
	weekrank := container.NewBorder(nil, nil, container.NewHBox(weekLabel1, weekCheck1), nil, container.NewGridWithColumns(2, weekEntry, weekprogress))
	monthrank := container.NewBorder(nil, nil, container.NewHBox(monthLabel1, monthCheck1), nil, container.NewGridWithColumns(2, monthEntry, monthprogress))

	f := 50
	RanksVar.Tops = f
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
		RanksVar.Tops = n
		ranks2update <- "已修改下载量为：" + strconv.Itoa(n)
		return nil
	}

	text := canvas.NewText("0~100", theme.PlaceHolderColor())
	text.TextSize = theme.CaptionTextSize()
	downloadtop := container.NewHBox(label, container.NewVBox(topEntry, text))

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
			rankdownloadPath.SetText(_rankdownloadpath)
			RanksVar.DownloadPath = _rankdownloadpath
			ranks2update <- "已修改保存地址：" + RanksVar.DownloadPath
		}, win)
	})

	path := container.New(layout.NewHBoxLayout(), openFolder, rankdownloadPath)

	var ctx context.Context
	var cancel context.CancelFunc

	downloadbutton := widget.NewButton("下载", func() {
	Flag:
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)

		if RanksVar.State.Load() != 0 {
			// 启动个小窗判断:1.重新下载;2.取消
			restart := false
			fmt.Println(1)
			confirmCallback := func(tr bool) {
				if tr {
					// 初始化，并
					RanksVar.State.Store(0)
					RanksVar.List = [3]int{}
					ranks2update <- "重新下载!!!!"
					cancel()
					restart = true
					fmt.Println(2)
					return
				}
			}

			cnf := dialog.NewConfirm("Confirmation", "确定要重新下载？", confirmCallback, win)
			cnf.SetDismissText("是")
			cnf.SetConfirmText("否")
			cnf.Show()
			fmt.Println(3)
			if restart {
				goto Flag
			}

		}

		ranks2update <- "下载启动!!!!"

		if dayCheck1.Checked {
			RanksVar.List[0] = 1
			RanksVar.State.Add(1)
			go rank.DownloadRank(ctx, "daily", home.CommonVar, RanksVar, dayChan)
			dayprogress.DisableColor = false
		}
		if weekCheck1.Checked {
			RanksVar.List[1] = 1
			RanksVar.State.Add(1)
			go rank.DownloadRank(ctx, "weekly", home.CommonVar, RanksVar, weekChan)
			weekprogress.DisableColor = false
		}
		if monthCheck1.Checked {
			RanksVar.List[2] = 1
			RanksVar.State.Add(1)
			go rank.DownloadRank(ctx, "monthly", home.CommonVar, RanksVar, monthChan)
			monthprogress.DisableColor = false
		}

		go func() {
			q1 := sync.Once{}
			q2 := sync.Once{}
			q3 := sync.Once{}
			for {
				select {
				case v := <-dayChan:
					dayprogress.SetValue(v)
					if v == 1 {
						q1.Do(func() {
							for _, d := range RanksVar.Day {
								_d := strings.SplitN(d[5:], "-", 2)
								key, _value := _d[0], _d[1]
								key = d[:5] + key
								value, _ := strconv.Atoi(_value)
								if _, ok := oldCheck[key]; ok == false {
									oldCheck[key] = make([]int, 0)
								}
								oldCheck[key] = append(oldCheck[key], value)
							}
							utils.SaveCache(oldCheck)
							RanksVar.List[0] = 0
							RanksVar.State.Add(-1)
							ranks2update <- "day类型下载完毕"
							if RanksVar.List == [3]int{} && RanksVar.State.Load() == 0 {
								cancel()
							}
						})
					}

				case v := <-weekChan:
					weekprogress.SetValue(v)
					if v == 1 {
						q2.Do(func() {
							for _, w := range RanksVar.Week {
								_w := strings.SplitN(w[5:], "-", 2)
								key, _value := _w[0], _w[1]
								key += w[:5]
								value, _ := strconv.Atoi(_value)
								if _, ok := oldCheck[key]; ok == false {
									oldCheck[key] = make([]int, 0)
								}
								oldCheck[key] = append(oldCheck[key], value)
							}
							utils.SaveCache(oldCheck)
							RanksVar.List[1] = 0
							RanksVar.State.Add(-1)
							ranks2update <- "week类型下载完毕"
							if RanksVar.List == [3]int{} && RanksVar.State.Load() == 0 {
								cancel()
							}
						})
					}
				case v := <-monthChan:
					monthprogress.SetValue(v)
					if v == 1 {
						q3.Do(func() {
							for _, m := range RanksVar.Month {
								_m := strings.SplitN(m[5:], "-", 2)
								key, _value := _m[0], _m[1]
								key += m[:5]
								value, _ := strconv.Atoi(_value)
								if _, ok := oldCheck[key]; ok == false {
									oldCheck[key] = make([]int, 0)
								}
								oldCheck[key] = append(oldCheck[key], value)
							}
							utils.SaveCache(oldCheck)
							RanksVar.List[2] = 0
							RanksVar.State.Add(-1)
							ranks2update <- "month类型下载完毕"
							if RanksVar.List == [3]int{} && RanksVar.State.Load() == 0 {
								cancel()
							}
						})
					}
				// 超时停止或下载完成
				case <-ctx.Done():
					newCheck = make(map[string]*[]int)
					if RanksVar.State.Load() == 0 {
						ranks2update <- "下载完毕!!!!"
						RanksVar.Day = []string{}
					} else {
						ranks2update <- "下载停止!!!!"
						RanksVar.State.Store(0)
					}
					RanksVar.List = [3]int{}
					return
				}
			}

		}()

	})
	downloadbutton.Importance = widget.HighImportance
	cancelbutton := widget.NewButton("取消", func() {
		if RanksVar.State.Load() != 0 {
			confirmCallback := func(tr bool) {
				if tr {
					// 初始化
					RanksVar.State.Store(0)
					RanksVar.List = [3]int{}
					ranks2update <- "取消下载!!!!"
					cancel()
					return
				}
			}
			cnf := dialog.NewConfirm("Confirmation", "确定要取消下载？", confirmCallback, win)
			cnf.SetDismissText("是")
			cnf.SetConfirmText("否")
			cnf.Show()
		}

	})

	bt := container.NewHBox(cancelbutton, downloadbutton)
	btx := container.NewBorder(nil, nil, nil, bt, bt)
	RanksVar.Do(func() {
		go func() {
			fmt.Println("监听启动")
			t := time.NewTicker(time.Minute)
			for {
				select {
				case s := <-ranks2update:
					l.Send(4, s, 2)
				case <-t.C:
					l.Send(4, fmt.Sprintf("rank heart", &RanksVar), 1)
				}
			}
		}()

	})

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
