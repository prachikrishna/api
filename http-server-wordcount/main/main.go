package main

import (
	"http-server-wordcount/handler"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/wordcount/", handler.WordCountHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
