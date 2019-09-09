package main

import (
	"strconv"
	_ "test/test1/model"
	"test/test1/server"
)

func main() {
	server.Ps = strconv.Itoa(*server.P)
	if 	server.Load != nil{
		go server.Control(*server.P,*server.Load)
	}
	server.Run(server.P)
}




