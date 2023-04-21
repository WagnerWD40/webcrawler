package website

import (
	"crawler/db"
	"crawler/model"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	Links []model.VisitedLink
}

func Run() {
	tmpl, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links, err := db.ListLinks()
		if err != nil {
			fmt.Println(err)
		}

		data := Data{Links: links}

		tmpl.Execute(w, data)
	})

	fmt.Println("Serving in port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
