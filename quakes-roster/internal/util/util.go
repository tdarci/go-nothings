package util

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchDoc(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func DownloadFile(filepath, url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer out.Close()

	_, _ = out.ReadFrom(resp.Body)
}

func SaveJSON(path string, data any) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(data)
}

func LoadJSON(path string, v any) {
	f, _ := os.Open(path)
	defer f.Close()
	json.NewDecoder(f).Decode(v)
}

func Sanitize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
