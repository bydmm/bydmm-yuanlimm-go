package woker

import (
	"encoding/json"
	"fmt"
	"os"

	workers "github.com/jrallison/go-workers"
)

// Worker 工人通用接口
type Worker interface {
	Perform() error
}

// JobMessage 异步任务
type JobMessage struct {
	Class      string   `json:"class"`
	Args       []string `json:"args"`
	Retry      bool     `json:"retry"`
	Queue      string   `json:"queue"`
	Jid        string   `json:"jid"`
	CreatedAt  float64  `json:"created_at"`
	EnqueuedAt float64  `json:"enqueued_at"`
}

// DecodeJob 解析异步任务
func DecodeJob(message *workers.Msg) (JobMessage, error) {
	var msg JobMessage
	data := message.OriginalJson()
	err := json.Unmarshal([]byte(data), &msg)
	return msg, err
}

// Perform 通用执行
func Perform(msg JobMessage, w Worker) {
	err := w.Perform()
	if err != nil {
		fmt.Printf("Error Worker %s: %s\n", msg.Class, err.Error())
	}
}

// QueueHandle 分发任务
func QueueHandle(message *workers.Msg) {
	msg, err := DecodeJob(message)
	if err != nil {
		return
	}
	switch msg.Class {
	case "BuyPriceRankWorker":
		w := BuyPriceRankWorker{}
		Perform(msg, w)
	}
}

// Start 启动队列
func Start() {
	workers.Configure(map[string]string{
		"server":   os.Getenv("WORKER_REDIS_ADDR"),
		"password": os.Getenv("WORKER_REDIS_PW"),
		"database": os.Getenv("WORKER_REDIS_DB"),
		"pool":     "30",
		"process":  "1",
	})

	// pull messages from "myqueue" with concurrency of 10
	workers.Process("default", QueueHandle, 1)

	// stats will be available at http://localhost:8080/stats
	go workers.StatsServer(8080)

	// Blocks until process is told to exit via unix signal
	workers.Run()
}
