package main

import (
	"os"
	"time"
	"yuanlimm-worker/cache"
	"yuanlimm-worker/cron"
	"yuanlimm-worker/model"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	cache.Redis()
	model.Database(os.Getenv("MYSQL_DSN"))
}

func main() {
	cron.Start()

	for true {
		time.Sleep(2 * time.Second)
	}
	// woker.Start()
}
