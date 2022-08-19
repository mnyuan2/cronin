package main

import "cron/internal/server"

func main(){
	r := server.Init()
	r.Run(":8081")
}
