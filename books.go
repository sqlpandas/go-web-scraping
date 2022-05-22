package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type Book struct {
	Title string
	Price string
}

func crawl() {
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{"Title", "Price"}
	err = writer.Write(headers)
	if err != nil {
		return
	}

	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML(".next > a", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		err := c.Visit(nextPage)
		if err != nil {
			return
		}
	})

	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		book := Book{}
		book.Title = e.ChildAttr(".image_container img", "alt")
		book.Price = e.ChildText(".price_color")
		row := []string{book.Title, book.Price}
		err := writer.Write(row)
		if err != nil {
			return
		}
	})

	startUrl := fmt.Sprintf("https://books.toscrape.com/")
	err = c.Visit(startUrl)
	if err != nil {
		return
	}
}

func main() {
	crawl()
}
