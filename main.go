package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type opts struct {
	StackExchangeAccessToken string `long:"access-token" env:"STACK_EXCHANGE_ACCESS_TOKEN" description:"Stack Exchange Access Token, obtainable through OAuth2, see https://api.stackexchange.com/docs/authentication"`
	StackExchangeKey         string `long:"key" env:"STACK_EXCHANGE_KEY" description:"Stack Exchange Application Key, if you don't have a request key you can obtain one by registering your application on Stack Apps."`
}

func main() {
	if len(os.Args) == 0 {
		lambda.Start(handleRequest)
	}

	var opts opts
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	questionsToLocalFiles(opts)
}

func questionsToLocalFiles(opts opts) {
	questions, err := getQuestions(opts.StackExchangeAccessToken, opts.StackExchangeKey)
	if err != nil {
		log.Fatal(err)
	}

	questionsJson, err := json.Marshal(questions)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("output.json", questionsJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

	questionVertices, answerVertices, peopleVertices, edges := toVerticesAndEdges(questions)

	err = writeToLocalFile(questionVertices, "question_vertices.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	log.Printf("Written %v\n", "question_vertices.csv")

	err = writeToLocalFile(answerVertices, "answer_vertices.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	log.Printf("Written %v\n", "answer_vertices.csv")

	err = writeToLocalFile(peopleVertices, "people_vertices.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	log.Printf("Written %v\n", "people_vertices.csv")

	err = writeToLocalFile(edges, "edges.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	log.Printf("Written %v\n", "edges.csv")

}

func writeToLocalFile(records [][]string, fileName string) error {
	log.Printf("Write %d records, file name: %v\n", len(records), fileName)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			log.Fatalf("Failed to write record %v. Err: %v\n", record, err.Error())
		}
	}
	return nil
}
