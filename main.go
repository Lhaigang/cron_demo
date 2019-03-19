package main

import (
	"cron_demo/http_task"
	"cron_demo/task"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	//i := 0
	//c := cron.New()
	//spec := "0 */1 * * * *"
	//c.AddFunc(spec, func() {
	//	i++
	//	log.Println("execute per second", i)
	//})
	//c.Start()
	//select {}

	task := task.LoopTask{
		C: 2, //并发数
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
