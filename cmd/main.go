package main

import (
	"fmt"
	"log"
	"net/http"
)

// my id
var whiteListID = []int{123}

func allMightHandler(writer http.ResponseWriter, request *http.Request) {
	// намапить пришедший джейсон на структуру (вероятно анонимную)
	// извлечь аудишник и сравнить с моим захардкоженным
	// если все ок - выполнять логику
}

func main() {
	fmt.Println("hello world")

	http.HandleFunc("/", allMightHandler)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal(err)
	}
}
