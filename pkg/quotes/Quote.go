package quotes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
)

type Quote struct {
	Id    int    `json:"id"`
	Author string `json:"author"`
	Book   string `json:"book"`
	Categories []string `json:"categories"`
	Quote string `json:"quote"`
}

func getAllQuotes() []Quote {
	var quotes []Quote

	absPath, err := filepath.Abs("pkg/quotes/quotes.json")
	if err != nil {
        log.Fatal(err)
		return quotes
    }
	file, readFileErr := ioutil.ReadFile(absPath)
	if readFileErr != nil {
		log.Fatal("Error when opening file: ", readFileErr)
	}

	unmarshalErr := json.Unmarshal(file, &quotes)

	if unmarshalErr != nil {
        log.Fatal("Error when unmarshalling: ", unmarshalErr)
		return quotes
    }

	return quotes
}


func GetRandomQuote() Quote {
    quotes := getAllQuotes()
    return quotes[rand.Intn(len(quotes))]
}
