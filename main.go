package main

import (
	"github.com/makalexs/godr/httpserver"
	"log"
)

func main() {

	if _, err := httpserver.StartServer(); err != nil {
		log.Fatalln("failed to start jsonrpc server")
		return
	}

}
