package main

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

func main() {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	check(err)

	data := struct {
		Name string
	}{
		Name: "watery light beer",
	}

	//index some data
	index.Index("id", data)

	//search for some text
	query := bleve.NewMatchQuery("watered down")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	check(err)
	fmt.Println(searchResult)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
