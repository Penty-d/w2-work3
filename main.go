package main

import (
	//"context"
	"log"
	"w2work3/internal/config"
	"w2work3/internal/repository"
	//"github.com/cloudwego/hertz/pkg/app"
	//"github.com/cloudwego/hertz/pkg/app/server"
	//"github.com/cloudwego/hertz/pkg/common/utils"
	//"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	_, err = repository.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("database connected")
}
