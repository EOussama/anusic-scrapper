package scraper

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eoussama/anusic-api/src/shared/models"
	"github.com/eoussama/anusic-api/src/shared/utils"
	"github.com/gocolly/colly"
)

// AnimeList scraps the entire anime list
func AnimeList() {
	log.Println("Scraping Anime list...")
	start := time.Now()

	// Initializing the scraper
	collector := colly.NewCollector(colly.Async(true))

	// Initializing the anime list
	animeTitles := []models.Anime{}

	// Scraping the Anime list
	collector.OnHTML("#wiki_0-9 ~ p", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("a", func(_ int, element *colly.HTMLElement) bool {

			// Extracting anime title
			extract := element.Text
			idx := strings.LastIndex(extract, " (")
			year, _ := strconv.ParseInt(extract[idx+2:len(extract)-1], 10, 16)

			anime := models.Anime{
				ID:   strings.Trim(strings.Replace(extract[:idx], " ", "", -1), " "),
				Name: strings.Trim(strings.Replace(extract[:idx], "\"", "", -1), " "),
				Year: uint16(year),
			}

			// Appending extracted anime title
			animeTitles = append(animeTitles, anime)
			return false
		})
	})

	// Visiting the target page and invoking the scraper
	collector.Visit(os.Getenv("BASE") + "anime_index")
	collector.Wait()

	log.Printf("Fetched %d Anime titles in %v", len(animeTitles), time.Since(start))
	utils.Cache.Anime = animeTitles
}

// AnimeInfo scraps Anime info
func AnimeInfo() {
	log.Println("Scraping Anime Info...")

	start := time.Now()
	count := 0

	// Initializing the scraper
	collector := colly.NewCollector(colly.Async(true))
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: len(utils.Cache.Anime),
	})

	// Scraping the Anime info
	collector.OnHTML(".md.wiki > h3", func(e *colly.HTMLElement) {

		// Getting the Anime index
		targetID := strings.Trim(strings.Replace(e.Text, " ", "", -1), " ")
		index, _ := utils.Cache.GetAnimeByID(targetID)

		if index > -1 {
			// Getting the respective Anime
			anime := &utils.Cache.Anime[index]

			// Extracting the MAL ID
			mal := e.ChildAttr("a", "href")
			idx := strings.LastIndex(mal, "/anime/")
			extr := mal[idx+len("/anime/") : len(mal)-1]
			id, err := strconv.ParseInt(extr, 10, 32)

			if err == nil {
				anime.MALID = uint16(id)
			}

			// Extracting the alt name
			if e.DOM.Next().Is("p") {
				altNamesStr := strings.Replace(e.DOM.Next().Text(), "\"", "", -1)
				altNamesFrg := strings.Split(altNamesStr, ",")
				anime.AltNames = []string{}

				for i := 0; i < len(altNamesFrg); i++ {
					anime.AltNames = append(anime.AltNames, altNamesFrg[i])
				}
			}

			// log.Printf("%+v\n", anime)
			count++
		} else {
			log.Printf("Anime “%s” not found", targetID)
		}
	})

	years := []string{"2020", "2019", "2018", "2017", "2016", "2015", "2014", "2013", "2012", "2011", "2010", "2009", "2008", "2007", "2006", "2005", "2004", "2003", "2002", "2001", "2000", "90s", "80s", "70s", "60s"}

	for _, year := range years {
		// Constructing the year index page
		url := os.Getenv("BASE") + year

		// Visiting the target page and invoking the scraper
		collector.Visit(url)
	}

	// Waiting for the scraping to resolve
	collector.Wait()

	log.Printf("Fetched %d Anime info in %v", count, time.Since(start))
}
