package main

import (
	"log"
	"net/http"
)

func main() {
	// create file server handler
	fs := http.FileServer(http.Dir("./webapp"))

	// start HTTP server with `fs` as the default handler
	log.Fatal(http.ListenAndServe(":9001", fs))
}

//func main() {
//	router := httprouter.New()
//	//router.GET("/", Index)
//	router.ServeFiles("/webapp/*filepath", http.Dir("./webapp"))
//
//	log.Fatal(http.ListenAndServe(":9001", router))
//}
