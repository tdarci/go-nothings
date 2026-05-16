package scraper

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tdarci/go-nothings/quakes-roster/internal/model"
	"github.com/tdarci/go-nothings/quakes-roster/internal/util"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://www.sjearthquakes.com/club/roster"

func Run() {
	doc := util.FetchDoc(baseURL)

	var players []model.Player
	var mu sync.Mutex

	sem := make(chan struct{}, 6) // concurrency limit (tune 4–8)

	var wg sync.WaitGroup

	doc.Find("div.oc-c-promo").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2.fa-text__title").Text())
		if title == "" {
			return
		}

		parts := strings.SplitN(title, " ", 2)

		number := ""
		name := title

		if len(parts) == 2 {
			number = strings.TrimPrefix(parts[0], "#")
			name = parts[1]
		}

		position := ""
		birthplace := ""

		s.Find(".fa-text__body p").Each(func(_ int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())

			if strings.HasPrefix(text, "Position:") {
				position = strings.TrimSpace(strings.TrimPrefix(text, "Position:"))
			}

			if strings.HasPrefix(text, "Hometown:") {
				birthplace = strings.TrimSpace(strings.TrimPrefix(text, "Hometown:"))
			}
		})

		img, _ := s.Find("img").Attr("data-src")
		if img == "" {
			img, _ = s.Find("img").Attr("src")
		}

		url, _ := s.Find("a.fa-button.-cta1").Attr("href")

		imgPath := "images/" + util.Sanitize(name) + ".jpg"
		if img != "" {
			util.DownloadFile(imgPath, img)
		}

		player := model.Player{
			Name:       name,
			Number:     number,
			Position:   position,
			Birthplace: birthplace,
			ImagePath:  imgPath,
			PlayerURL:  url,
		}

		// async enrichment (AGE)
		wg.Add(1)
		go func(p model.Player) {
			defer wg.Done()

			if p.PlayerURL == "" {
				mu.Lock()
				players = append(players, p)
				mu.Unlock()
				return
			}

			sem <- struct{}{}
			age := scrapeAge(p.PlayerURL)
			<-sem

			p.Age = age

			mu.Lock()
			players = append(players, p)
			mu.Unlock()

			time.Sleep(150 * time.Millisecond)
		}(player)
	})

	wg.Wait()

	util.SaveJSON("players.json", players)

	fmt.Printf("\nSaved %d players\n", len(players))
}

func scrapeAge(url string) string {
	doc := util.FetchDoc(url)

	var age string

	doc.Find(".mls-l-module--player-status-details__info").Each(func(i int, s *goquery.Selection) {

		title := strings.TrimSpace(s.Find("h3").Text())

		if strings.EqualFold(title, "Date of Birth") {
			text := strings.TrimSpace(s.Find("span").Text())

			// text looks like: "3.24.1995 (31)"
			start := strings.LastIndex(text, "(")
			end := strings.LastIndex(text, ")")

			if start != -1 && end != -1 && end > start {
				age = text[start+1 : end]
			}
		}
	})

	return age
}
