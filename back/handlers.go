package server

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var tmp *template.Template

func init() {
	tmp = template.Must(template.ParseGlob("front/html/*.html"))
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err := "404 Page not found"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 Method is not allowed"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusMethodNotAllowed)
		return
	}

	artists, err := GetAllArtists()
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	location, err := GetAllLocations()
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	allInfo := Everything{artists, location}

	err = tmp.ExecuteTemplate(w, "index.html", allInfo)
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
}

func InfoAboutArtist(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artists/" {
		err := "404 Page not found"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 method not allowed"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusMethodNotAllowed)
		return
	}
	artists, err := GetAllArtists()
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 || id > len(artists) || err != nil {
		err := "404 Page not found"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusNotFound)
		return
	}
	infoAboutOne, err := OneArtist(id)
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	rel, err := Relations(id)
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	artist := ArtistInfo{infoAboutOne, rel}
	err = tmp.ExecuteTemplate(w, "artist.html", artist)
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
}

func FilterhHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.URL.Path != "/filter/" {
		err := "404 Page not found"
		ErrorPage(w, err, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		err := "405 Method is not allowed"
		ErrorPage(w, err, http.StatusMethodNotAllowed)
		return
	}

	artists, err := GetAllArtists()
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
	location, err := GetAllLocations()
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}

	range1 := r.FormValue("range1")
	range2 := r.FormValue("range2")
	searchedArtists := Everything{artists, location}
	locations := r.FormValue("Locations")

	if locations != "" {
		searchedArtists, err = LocationSearch(searchedArtists, locations)
		if err != nil {
			err := "500 Internal Server Error"
			ErrorPage(w, err, http.StatusInternalServerError)
			return
		}
	}

	if range1 != "" && range2 != "" {
		searchedArtists, err = CreationDate(searchedArtists, range1, range2)
		if err != nil {
			err := "500 Internal Server Error"
			ErrorPage(w, err, http.StatusInternalServerError)
			return
		}
	}

	firstalbum1 := r.FormValue("range3")
	firstalbum2 := r.FormValue("range4")
	if firstalbum1 != "" && firstalbum2 != "" {
		searchedArtists, err = FirstAlbumDates(searchedArtists, firstalbum1, firstalbum2)
		if err != nil {
			err := "500 Internal Server Error"
			ErrorPage(w, err, http.StatusInternalServerError)
			return
		}
	}
	member := r.Form["interest"]
	if len(member) != 0 {
		searchedArtists, err = MembersNumber(searchedArtists, member)
		if err != nil {
			err := "500 Internal Server Error"
			ErrorPage(w, err, http.StatusInternalServerError)
			return
		}
	}

	err = tmp.ExecuteTemplate(w, "search.html", searchedArtists)
	if err != nil {
		err := "500 Internal Server Error"
		ErrorPage(w, err, http.StatusInternalServerError)
		return
	}
}

func ErrorPage(w http.ResponseWriter, errors string, code int) {
	w.WriteHeader(code)
	err:=tmp.ExecuteTemplate(w, "error.html", errors)
	if err != nil {
		err := "500 Internal Server Error"
		fmt.Println(err)
		http.Error(w, err, 500)
		return
	}
}
