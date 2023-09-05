package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url, image, name, price string
}

func main() {

	c := colly.NewCollector()
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed")
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	writer.Write(headers)

	var pokemonProducts []PokemonProduct

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnHTML("li.product", func(h *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = h.ChildAttr("a", "href")
		pokemonProduct.image = h.ChildAttr("img", "src")
		pokemonProduct.name = h.ChildText("h2")
		pokemonProduct.price = h.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(pokemonProducts[0].url)

		for _, pokemonProduct := range pokemonProducts {
			record := []string{
				pokemonProduct.url,
				pokemonProduct.image,
				pokemonProduct.name,
				pokemonProduct.price,
			}

			writer.Write(record)
		}
	})

	defer writer.Flush()

	// downloading the target HTML page
	err1 := c.Visit("https://scrapeme.live/shop/page/1/")

	if err1 != nil {
		log.Printf("failed to visit url: %v\n", err)
		return
	}
}
