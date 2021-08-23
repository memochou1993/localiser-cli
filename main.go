package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	Endpoint  = "http://127.0.0.1:8000/api" // FIXME: use env file
	ProjectId = "1"                         // FIXME: use env file
)

var (
	Dir    string
	Locale string
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type Language struct {
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

type Values map[string]string

func main() {
	parseFlags()
	if Locale != "" {
		if err := download(Locale); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := downloadAll(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&Dir, "o", ".", "output directory")
	flag.StringVar(&Locale, "l", "", "locale")
	flag.Parse()
}

func downloadAll() (err error) {
	languages, err := fetchLanguages()
	if err != nil {
		return
	}
	for _, language := range languages {
		if err = download(language.Locale); err != nil {
			return
		}
	}
	return
}

func download(locale string) (err error) {
	values, err := fetchValues(locale)
	if err != nil {
		return
	}
	data, err := json.MarshalIndent(values, "", "\t")
	if err != nil {
		return
	}
	filename := fmt.Sprintf("%s/%s.json", Dir, locale)
	if err = ioutil.WriteFile(filename, data, 0644); err != nil {
		return
	}
	return
}

func fetchLanguages() (languages []Language, err error) {
	b, err := fetch(fmt.Sprintf("projects/%s/cache/languages", ProjectId))
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &languages); err != nil {
		return
	}
	return
}

func fetchValues(locale string) (values Values, err error) {
	b, err := fetch(fmt.Sprintf("projects/%s/cache/values?locale=%s", ProjectId, locale))
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &values); err != nil {
		return
	}
	return
}

func fetch(target string) (b []byte, err error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s", Endpoint, target))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Unexpected response code: %v", resp.StatusCode))
		return
	}
	defer closeBody(resp.Body)
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func closeBody(reader io.ReadCloser) {
	if err := reader.Close(); err != nil {
		log.Fatal(err)
	}
}
