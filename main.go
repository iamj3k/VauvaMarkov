package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mb-14/gomarkov"
)

var (
	trainBool bool
	topics    []string
	threads   []string
	comments  []string
	page      int
	filename  string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Input parameter: <true/false>\ntrue: Generate a model based on statements.txt\nfalse: Generate a chain based on the generated model")
		os.Exit(1)
	}
	if os.Args[1] == "vauva" {
		topics = GetTopics()
		for _, topic := range topics {
			threads = GetThreads(topic)
			for _, thread := range threads {
				comments = GetComments(thread)
				for _, comment := range comments {
					writeStatement(comment)
				}
			}

		}
		os.Exit(0)
	}
	if os.Args[1] == "server" {
		http.HandleFunc("/vauva/", viewHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
	trainBool, err := strconv.ParseBool(os.Args[1])
	if err != nil {
		fmt.Println("\nINVALID SYNTAX!\nInput parameter: <true/false>\ntrue: Generate a model based on statements.txt\nfalse: Generate a chain based on the generated model")
		fmt.Println()
		os.Exit(1)
	}
	train := flag.Bool("train", trainBool, "Train the markov chain")
	order := flag.Int("order", 10, "Chain order to use")

	flag.Parse()
	if *train {
		chain := buildModel(*order)
		saveModel(chain)
	} else {
		chain, err := loadModel()
		if err != nil {
			fmt.Println(err)
			return
		}
		generateStatement(chain)
	}
}

func writeStatement(statement string) {
	filename = "statements.txt"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(statement); err != nil {
		panic(err)
	}
}

func buildModel(order int) *gomarkov.Chain {
	chain := gomarkov.NewChain(order)
	for _, data := range getDataset("statements.txt") {
		chain.Add(split(data))
	}
	return chain
}

func split(str string) []string {
	return strings.Split(str, "")
}

func getDataset(fileName string) []string {
	file, _ := os.Open(fileName)
	scanner := bufio.NewScanner(file)
	var list []string
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	return list
}

func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func saveModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func generateStatement(chain *gomarkov.Chain) string {
	order := chain.Order
	tokens := make([]string, 0)
	for i := 0; i < order; i++ {
		tokens = append(tokens, gomarkov.StartToken)
	}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - order):])
		tokens = append(tokens, next)
	}
	fmt.Println(strings.Join(tokens[1:len(tokens)-1], ""))
	return fmt.Sprintf(strings.Join(tokens[1:len(tokens)-1], ""))
}
