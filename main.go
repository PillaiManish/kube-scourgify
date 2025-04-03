package main

import (
	"kube-scourgify/cmd"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			logger.Println(r)
			return
		}
	}()
	cmd.Execute()
}
