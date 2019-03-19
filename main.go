package main

import (
	"github.com/robfig/cron"
	"log"
)

func main() {
	i := 0
	c := cron.New()
	spec := "0 */1 * * * *"
	c.AddFunc(spec, func() {
		i++
		log.Println("execute per second", i)
	})
	c.Start()
	select {}
}
