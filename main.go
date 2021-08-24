package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

var (
	config Config
	locale string
)

type Config struct {
	Endpoint        string `yaml:"endpoint"`
	ProjectID       string `yaml:"project_id"`
	OutputDirectory string `yaml:"output_directory"`
}

type Language struct {
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

type Values map[string]string

func main() {
	parseFlags()
	if err := parseConfig(); err != nil {
		log.Fatal(err)
	}
	if locale != "" {
		if err := download(locale); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := downloadAll(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig() (err error) {
	file, err := ioutil.ReadFile("localiser.yaml")
	if err != nil {
		return
	}
	return yaml.Unmarshal(file, &config)
}

func parseFlags() {
	flag.StringVar(&locale, "l", "", "locale")
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
	mode := os.ModePerm
	if _, err = os.Stat(config.OutputDirectory); os.IsNotExist(err) {
		if err = os.Mkdir(config.OutputDirectory, mode); err != nil {
			return
		}
	}
	filename := fmt.Sprintf("%s/%s.json", config.OutputDirectory, locale)
	if err = ioutil.WriteFile(filename, data, mode); err != nil {
		return
	}
	return
}

func fetchLanguages() (languages []Language, err error) {
	b, err := fetch(fmt.Sprintf("projects/%s/cache/languages", config.ProjectID))
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &languages); err != nil {
		return
	}
	return
}

func fetchValues(locale string) (values Values, err error) {
	b, err := fetch(fmt.Sprintf("projects/%s/cache/values?locale=%s", config.ProjectID, locale))
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &values); err != nil {
		return
	}
	return
}

func fetch(target string) (b []byte, err error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s", config.Endpoint, target))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected response code: %v", resp.StatusCode)
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
