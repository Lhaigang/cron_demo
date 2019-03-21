package main

import (
	"cron_demo/http_task"
	"cron_demo/task"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {

	timeStr:=time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(timeStr)

	task := task.LoopTask{
		C: 50, //并发数
	}
	manager := http_task.NewManager(task)
	fmt.Println("开始压测请等待")
	c := make(chan os.Signal, 1)
	go func() {
		task.Run(manager)
	}()
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c
	fmt.Println("准备停止")
	task.Stop()
	task.Wait()
	fmt.Println("压测完成")
}
