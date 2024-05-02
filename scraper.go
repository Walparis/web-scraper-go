package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Pokemon struct {
	Name  string
	Price string
	Image string
}

func main() {
	// Create a new Colly collector
	c := colly.NewCollector()

	// Slice to hold scraped Pokemon data
	var pokemons []Pokemon

	// Visit the target webpage
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemon := Pokemon{
			Name:  e.ChildText("h2"),
			Price: e.ChildText("span.amount"),
		}
		// Extract image URL
		pokemon.Image = e.ChildAttr("img.wp-post-image", "src")
		pokemons = append(pokemons, pokemon)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Scraping finished")
	})

	// Start scraping
	err := c.Visit("https://scrapeme.live/shop/")
	if err != nil {
		log.Fatalln("Error visiting website:", err)
	}

	// Create CSV file to store scraped data
	file, err := os.Create("pokemons.csv")
	if err != nil {
		log.Fatalln("Error creating CSV file:", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	err = writer.Write([]string{"Name", "Price", "Image"})
	if err != nil {
		log.Fatalln("Error writing CSV header:", err)
	}

	// Write scraped data to CSV
	for _, pokemon := range pokemons {
		err := writer.Write([]string{pokemon.Name, pokemon.Price, pokemon.Image})
		if err != nil {
			log.Println("Error writing CSV record:", err)
		}
	}

	log.Println("Scraped data has been saved to 'pokemons.csv'")
}
