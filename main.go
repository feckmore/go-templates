package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

// TODO: because the yaml is parsed into maps, order is not retained
// figure out a way to enforce order in the yaml

type Locale map[string]interface{}

func main() {
	r := mux.NewRouter()
	r.Methods("GET").PathPrefix("/").HandlerFunc(pageHandler)

	port := ":8000"
	fmt.Printf("listening on localhost%v\n", port)
	http.ListenAndServe(port, r)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	app, lang := parseURL(r)

	l, err := locales()
	if err != nil {
		fmt.Println(err)
		return
	}

	templateFile := fmt.Sprintf("%v.html", strings.Trim(app, "/"))
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Fprint(w, "Error parsing template:", err)
		return
	}
	err = t.Execute(w, l[lang])
	if err != nil {
		fmt.Fprint(w, "Error executing template:", err)
	}
}

func parseURL(r *http.Request) (string, string) {
	app, lang := "oi", "en" // defaults

	// used to split by slash, and not add empty elements to the slice
	splitFn := func(c rune) bool {
		return c == '/'
	}
	segments := strings.FieldsFunc(r.URL.Path, splitFn)

	if len(segments) > 0 {
		app = segments[0]
	}
	if len(segments) > 1 {
		lang = segments[1]
	}

	switch lang {
	case "de", "en":
		// do nothing
	default:
		lang = "en"
	}

	return app, lang
}

// locales returns a map of terms for a YAML file corresponding to the given app & locale.
func locales() (Locale, error) {
	content, err := ioutil.ReadFile("oi.yml")
	if err != nil {
		return nil, errors.New("Error reading locales file: " + err.Error())
	}

	l := Locale{}
	err = yaml.Unmarshal(content, &l)
	if err != nil {
		return nil, errors.New("Error unmarshalling file:" + err.Error())
	}

	return l, nil
}
