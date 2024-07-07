package rank

import (
	"context"
	"fmt"
	"gpixivImageDownload/log"
	"gpixivImageDownload/model"
	"gpixivImageDownload/model/utils"
	"gpixivImageDownload/pkg/core"
	"log/slog"
	"strconv"
	"sync"
)

//var Rootpath = (conf.ConfigData["DownloadControl"]["Path"]).(string)

//const Spt = string(os.PathSeparator)

var l *log.Logs

func DownloadRank(ctx context.Context, dwtype string, cmpts *model.Common, rkpts *model.Ranks, pip chan float64) {
	ctxChild, _ := context.WithCancel(ctx)
	var rootPath string
	if rkpts.DownloadPath != "" {
		rootPath = rkpts.DownloadPath
	} else {
		rootPath = cmpts.DownloadPath
	}

	//validModes := []string{"daily", "weekly", "monthly", "rookie", "original", "male", "female"}
	//validContents := []string{"all", "illust", "ugoira", "manga"}

	utils.QuoteOrCreateFile(rootPath)

	l.Send(slog.LevelInfo, fmt.Sprintf(
		"Downloading Setting View:R18 %v,Content %s,TopPage %d, SkipIllus %v, SkipUgoira %v, SkipManga %v",
		cmpts.R18, rkpts.Content(), rkpts.Tops, cmpts.SkipIllus, cmpts.SkipUgoira, cmpts.SkipManga), log.LogStdouts|log.LogFiles)

	//var InSCount = 0
	var DwCount = 0
	var ErrCount = 0
	//var DwCounts = 0
	var AllCount = 0
	var SkipCount = 0

	var page = 1
	var total = rkpts.Tops
	wd := sync.Once{}

	l.Send(slog.LevelInfo, fmt.Sprintf("---------- DownLoadType %s ----------", dwtype), 3)
	var content = "all"
	var dayList []string
	switch dwtype {
	case "day":
		dayList = rkpts.Day
	case "week":
		dayList = rkpts.Week
	case "month":
		dayList = rkpts.Month
	}

	for _, day := range dayList {
		for AllCount <= total {
			page = AllCount/50 + 1
			l.Send(slog.LevelInfo, fmt.Sprintf("** Page %d **", page), 3)
			ranks, output := core.GetPixivRanking(dwtype, day, *cmpts, page)
			l.Send(slog.LevelInfo, output, log.LogFiles|log.LogStdouts)

			if ranks == nil {
				break
			}

			wd.Do(func() {
				l.Send(slog.LevelInfo, fmt.Sprintf("*Mode :%s", ranks.Mode), log.LogStdouts)
				l.Send(slog.LevelInfo, fmt.Sprintf("*Content :%s", ranks.Content), log.LogStdouts)
				l.Send(slog.LevelInfo, fmt.Sprintf("*Total :%d", ranks.RankTotal), log.LogStdouts)
			})

			for k := 0; k < len(ranks.Contents); k++ {
				select {
				case <-ctxChild.Done():
					return
				default:

				}
				//if IpI.IncludeSkipTime {
				//	InSCount = AllCount
				//} else {
				//	InSCount = DwCount
				//}
				//
				//if InSCount == IpI.RankTop {
				//	break bk
				//}

				post := ranks.Contents[k]
				//AllCount += 1
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
				l.Send(slog.LevelInfo, fmt.Sprintf("{page:%d}.Image id %s", page, imageId), log.LogStdouts)
				pgCount, _ := strconv.Atoi(post.IllustPageCount)

				result, err := core.ProcessRankImage(pgCount, imageId, rootPath, cmpts.R18, dwtype, day, content)

				if err != nil {
					l.Send(slog.LevelDebug, fmt.Sprintf("Dowdload err:%v", err), log.LogStdouts)
				}
				DwCount += result[0]
				ErrCount += result[1]

				AllCount += 1
				l.Send(slog.LevelDebug, fmt.Sprintf("result:%v", result), log.LogStdouts)
				pip <- float64(AllCount) / float64(total)
			}

			// 判断是否是最后一页
			if ranks.NextPage == false {
				break
			}
		}

	}

	feedback := fmt.Sprintf("预计下载前top%d，需下载%d张图片，成功下载%d张图片，下载失败%d张图片，跳过%d张图片", rkpts.Tops, AllCount, DwCount, ErrCount, SkipCount)
	l.Send(slog.LevelInfo, feedback, log.LogFiles|log.LogStdouts)
	//globalOptions.Auth.Output = feedback
	//globalOptions.DB.SaveAuthDownloadRecord(globalOptions.Auth)

}
