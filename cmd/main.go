package main

import (
	"github.com/go-well/spider/spy"
	"log"
	"time"
)

func main() {

	cli, err := spy.Connect(spy.Options{
		Addr: "127.0.0.1:8888",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("连接成功", cli)

	for {
		time.Sleep(time.Minute)
	}

}
