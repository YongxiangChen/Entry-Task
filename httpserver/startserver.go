package main

import (
	"entrytask1/httpserver/server"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	mux := &server.MyMux{}
	err := http.ListenAndServe(":8001", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}