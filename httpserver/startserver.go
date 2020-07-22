package main

import (
	"encoding/gob"
	"entrytask1/httpserver/server"
	"entrytask1/tcpserver/model"
	"log"
	"net/http"
)

func main() {
	gob.Register(model.User{})
	mux := &server.MyMux{}
	err := http.ListenAndServe(":8001", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}