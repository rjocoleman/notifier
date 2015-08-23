package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type payload struct {
	App     string `json:"app"`
	Release string `json:"release"`
	URL     string `json:"url"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		var p payload
		err = json.Unmarshal(body, &p)
		if err != nil {
			log.Fatal(err)
		}

		notifyName := strings.TrimPrefix(req.URL.Path, "/")
		log.Print(fmt.Sprintf("Received Webhook for %s\n", notifyName))
		notifyURL := os.Getenv(strings.ToUpper(notifyName) + "_URL")
		notifyTemplate := notifyName + ".tmpl"
		notifyTemplatePath := "templates/" + notifyTemplate

		if notifyURL == "" {
			log.Fatal(fmt.Sprintf("Notification URL not configured for %s. Please set ENV %s\n", notifyName, strings.ToUpper(notifyName)+"_URL"))
		}
		if _, err := os.Stat(notifyTemplatePath); err != nil {
			log.Fatal(fmt.Sprintf("Template (%s) doesn't exist or can't be accessed", notifyTemplatePath))
		}
		pr, pw := io.Pipe()

		go func() {
			defer pw.Close()

			t := template.Must(template.New(notifyTemplate).ParseFiles(notifyTemplatePath))
			err = t.ExecuteTemplate(pw, notifyTemplate, p)

			if err != nil {
				log.Fatal(err)
			}
		}()

		log.Print(fmt.Sprintf("Sending Notification to %s\n", notifyURL))
		notifyReq, err := http.NewRequest("POST", notifyURL, pr)

		if err != nil {
			log.Fatal(err)
		}
		_, err = http.DefaultClient.Do(notifyReq)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":3000", mux)
}
