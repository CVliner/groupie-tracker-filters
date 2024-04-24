package main

import (
	"fmt"
	server "groupie-tracker/back"
	"log"
	"net/http"
)

func main() {
	
	fs := http.FileServer(http.Dir("./front/css_img"))
	http.Handle("/css_img/", http.StripPrefix("/css_img", fs))
	http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/artists/", server.InfoAboutArtist)
	http.HandleFunc("/filter/", server.FilterhHandler)
	fmt.Println("http://localhost:4747/")
	err := http.ListenAndServe(":4747", nil)
	if err != nil {
		log.Fatal(err)
	}
}
