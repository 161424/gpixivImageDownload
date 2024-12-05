package rank

import (
	"context"
	"fmt"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"gpixivImageDownload/pkg/core"
	"gpixivImageDownload/pkg/utils/browser"
	"log/slog"
	"strconv"
	"strings"
)

//var Rootpath = (conf.ConfigData["DownloadControl"]["Path"]).(string)

//const Spt = string(os.PathSeparator)

var l = *log.Logger

func DownloadRank(ctx context.Context, dwtype string, cmpts *model.Common, rkpts *model.Ranks, pip chan float64) {

	stopDw := false
	go func() {
		<-ctx.Done()
		stopDw = true
	}()

	var rootPath string
	if rkpts.DownloadPath != "" {
		rootPath = rkpts.DownloadPath
	} else {
		rootPath = cmpts.DownloadPath
		err := utils.QuoteOrCreateFile(rootPath)

		if err != nil {
			l.Send(slog.LevelWarn, err.Error(), 2)
			rootPath = rkpts.DownloadPath
		}
	}

	l.Send(slog.LevelInfo, fmt.Sprintf(
		"---------- DownLoadType %s ----------"+"\n"+
			"Downloading Setting View:R18 %v,Content %s,TopPage %d, SkipIllus %v, SkipUgoira %v, SkipManga %v",
		dwtype, cmpts.R18, rkpts.Content(), rkpts.Tops, cmpts.SkipIllus, cmpts.SkipUgoira, cmpts.SkipManga), log.LogStdouts)

	var DwCount int64 = 0
	var ErrCount int64 = 0
	var SkipCount = 0

	var page = 1
	var total = rkpts.Tops
	//wd := sync.Once{}

	var content = "all"
	var dateList []string
	switch dwtype {
	case "daily":
		dateList = rkpts.Day
	case "weekly":
		dateList = rkpts.Week
	case "monthly":
		dateList = rkpts.Month
	}
	browser.SetMutliHttps(cmpts.Ck)
	fmt.Println(dateList, total)

	pip <- 0.0
	for dayidx, day := range dateList {
		day = strings.Join(strings.Split(day, "-"), "")
		var DwTop = 0
		for DwTop < total {
			page = DwTop/50 + 1
			// 分析排行榜页面
			ranks, err := core.GetPixivRanking(dwtype, day, cmpts, page)
			//fmt.Println(ranks, output)
			l.Send(slog.LevelInfo, fmt.Sprintf("** (%s) Page %d **", dwtype, page), 3)
			fmt.Println("?1")
			if ranks == nil || ranks.Contents == nil {
				l.Send(slog.LevelInfo, fmt.Sprintf("rank is nil or err:%s", err), 3)
				break
			}

			for k := 0; k < len(ranks.Contents); k++ {
				// 下载截至条件
				if stopDw {
					goto td
				}
				if DwTop >= total {
					break
				}
				l.Send(slog.LevelInfo, fmt.Sprintf("正在下载第%d个图集", (page-1)*50+k+1), 3)
				post := ranks.Contents[k]
				switch post.IllustType {
				case "0":
					if cmpts.SkipIllus {
						content = "noIllus"
						SkipCount++
						continue
					}
				case "1":
					if cmpts.SkipManga {
						content = "noManga"
						SkipCount++
						continue
					}
				case "2":
					if cmpts.SkipUgoira {
						content = "noUgoira"
						SkipCount++
						continue
					}
				}

				imageId := fmt.Sprintf(strconv.FormatFloat(post.IllustId, 'f', 0, 64))
				pgCount, _ := strconv.Atoi(post.IllustPageCount)
				result, err := core.ProcessRankImage(ctx, pgCount, k+1, imageId, rootPath, cmpts.R18, dwtype, day, content, cmpts.MThread == true)

				if err != nil {
					l.Send(slog.LevelDebug, fmt.Sprintf("Dowdload err:%v.{page:%d}.Image id %s", err, page, imageId), log.LogStdouts)
				}

				DwCount += result[0]
				ErrCount += result[1]
				DwTop += 1
				dwPer := float64(DwTop+total*dayidx) / float64(total*len(dateList))
				pip <- dwPer
				if dwPer == 1.0 {
					break
				}
			}

			// 判断是否是最后一页
			if ranks.NextPage == false {
				break
			}
		}

	}
	fmt.Println("?2")
td:
	fmt.Println("?3")
	pip <- 1.0
	fmt.Println("?4")
	feedback := fmt.Sprintf("预计下载%v前top%d，成功下载%d张图片，下载失败%d张图片，跳过%d张图片", dateList, rkpts.Tops, DwCount, ErrCount, SkipCount)
	l.Send(slog.LevelInfo, feedback, log.LogFiles|log.LogStdouts)

	return
}
