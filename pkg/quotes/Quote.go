package quotes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Category string

const (
	Mythology Category = "mythology"
	Education          = "education"
)

type JsonQuote struct {
	Author     string     `json:"author"`
	Book       *string    `json:"book"`
	Categories []Category `json:"categories"`
	Quote      string     `json:"quote"`
}

type JsonQuotes map[string]JsonQuote

type Quote struct {
	JsonQuote
	Id string `json:"id"`
}

func getAllQuotes() JsonQuotes {
	var quotes JsonQuotes

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

func GetQuoteById(id string) Quote {
	quotes := getAllQuotes()
	return Quote{
		JsonQuote: quotes[id],
        Id:        id,
	}
}

func GetRandomQuote() Quote {
	quotes := getAllQuotes()
	for k := range quotes {
		return Quote{
			JsonQuote: quotes[k],
            Id:        k,
		}
	}
	return Quote{
		JsonQuote: quotes["1"],
        Id:        "1",
	}
}
