package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mb-14/gomarkov"
)

var (
	quote string
	chain *gomarkov.Chain
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	chain, err := loadModel()
	if err != nil {
		fmt.Println(err)
		return
	}
	quote = generateStatement(chain)
	quote = strings.Replace(quote, "$", "", -1)
	fmt.Fprintf(w, "<h1>Vauva.fi Markov</h1><p>%s</p>", quote)
}
