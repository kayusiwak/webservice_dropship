package main

import "fmt"
import "html/template"
import "net/http"

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var t, err = template.ParseFiles("index2.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		t.Execute(w, nil)
	})

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
