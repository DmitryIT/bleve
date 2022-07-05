package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
)

func main() {
	objects, err := ReadObjectsFromConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range objects {
		fmt.Printf("%+v\n", object)
		index.Index(object.Id, object)
	}

	searchString := "ed"
	//query := bleve.NewQueryStringQuery(searchString)
	query := bleve.NewFuzzyQuery(searchString)
	query.SetFuzziness(1)
	searchRequest := bleve.NewSearchRequest(query)
	result, err := index.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Search for %s returned:\n", searchString)
	fmt.Printf("%+v\n", result)

}

// Struct Config contains config.json file
type Config struct {
	Objects []Object
}

type Object struct {
	Id       string   `json:"id"`
	Keywords []string `json:"keywords"`
}

func ReadObjectsFromConfig(fileName string) ([]Object, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var objects []Object
	err = json.NewDecoder(file).Decode(&objects)
	if err != nil {
		return nil, err
	}

	return objects, nil
}
