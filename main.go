package main

import (
	"os"
	"yuanlimm-worker/cache"
	"yuanlimm-worker/cron"
	"yuanlimm-worker/model"
	"yuanlimm-worker/server"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	cache.Redis()
	model.Database(os.Getenv("MYSQL_DSN"))
}

func main() {
	cron.Start()

	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
