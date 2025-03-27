package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", listManuals)
	http.HandleFunc("/view", viewManual)
	http.HandleFunc("/edit", editManual)
	http.HandleFunc("/save", saveManual)

	http.ListenAndServe(":8080", nil)
}

func listManuals(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./manuals")
	if err != nil {
		http.Error(w, "Unable to list manuals", http.StatusInternalServerError)
		return
	}

	var manuals []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			manuals = append(manuals, file.Name())
		}
	}

	templates.ExecuteTemplate(w, "list.html", manuals)
}

func viewManual(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Manual name is required", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, "./manuals/"+name)
}

func editManual(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	content := ""

	if name != "" {
		data, err := ioutil.ReadFile("./manuals/" + name)
		if err == nil {
			content = string(data)
		}
	}

	templates.ExecuteTemplate(w, "edit.html", map[string]interface{}{
		"Name":    name,
		"Content": content,
	})
}

func saveManual(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	content := r.FormValue("content")

	if name == "" {
		http.Error(w, "Manual name is required", http.StatusBadRequest)
		return
	}

	err := os.WriteFile("./manuals/"+name+".html", []byte(content), 0644)
	if err != nil {
		http.Error(w, "Unable to save manual", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
