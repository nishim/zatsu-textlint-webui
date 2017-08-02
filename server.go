package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func getHtml(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
			fmt.Fprint(w, "error.")
		}
		t.Execute(w, nil)
	})

	http.HandleFunc("/lint/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		r.ParseForm()
		html, err := getHtml(r.Form["url"][0])
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("npm", "run", "--silent", "textlint")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, html)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		fmt.Fprint(w, string(out))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
