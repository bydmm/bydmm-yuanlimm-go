package cron

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron"
)

// Cron 定时器单例
var Cron *cron.Cron

// Run 运行
func Run(job func()) {
	from := time.Now().UnixNano()
	job()
	to := time.Now().UnixNano()
	jobName := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	fmt.Printf("%s: %dms\n", jobName, (to-from)/int64(time.Millisecond))
}

// Start 启动定时任务
func Start() {
	if Cron == nil {
		Cron = cron.New()
	}

	if os.Getenv("DEBUG") == "" {
		Cron.AddFunc("@every 1m", func() { Run(MatchingTrade) })
		Cron.AddFunc("@hourly", func() { Run(TrendHour) })
		Cron.AddFunc("@daily", func() { Run(TrendHour) })
		Cron.AddFunc("@daily", func() { Run(LoveClear) })
		Cron.AddFunc("@daily", func() { Run(GreenHatChecker) })
		Cron.AddFunc("@every 10m", func() { Run(UserRank) })
		Cron.AddFunc("@every 11m", func() { Run(HotRank) })
		Cron.AddFunc("@every 12m", func() { Run(BuyPriceRank) })
		Cron.AddFunc("@every 13m", func() { Run(SalePriceRank) })
		Cron.AddFunc("@every 14m", func() { Run(MarketValueRank) })
	}
	Cron.Start()
	fmt.Println("Cron Job Start")
}
