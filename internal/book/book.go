package book

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"time"
)

type Book struct {
	Quotes []Record `json:"quotes"`
}

type Record struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

//go:embed quotes.json
var quotesRaw []byte

func New() (Book, error) {
	var b Book
	if err := json.Unmarshal(quotesRaw, &b); err != nil {
		return Book{}, err
	}

	return b, nil
}

func (q Book) RandomQuote() string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	idx := r.Intn(len(q.Quotes))
	rec := q.Quotes[idx].Quote + " by " + q.Quotes[idx].Author

	return rec
}
