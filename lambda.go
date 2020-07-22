package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

type myEvent struct {
	Name string `json:"name"`
}

func handleRequest(_ context.Context, name myEvent) (string, error) {
	questions, err := getQuestions("", "")
	if err != nil {
		log.Fatal(err)
	}

	questionVertices, answerVertices, peopleVertices, edges := toVerticesAndEdges(questions)

	s3Region, ok := os.LookupEnv("REGION")
	if !ok {
		log.Println("REGION env variable is not set, using default region us-east-1")
		s3Region = "us-east-1"
	}
	s3Bucket, ok := os.LookupEnv("BUCKET")
	if !ok {
		log.Fatalf("BUCKET env variable is not set")
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	if err != nil {
		log.Fatal(err)
	}

	uploadResult, err := writeToS3(questionVertices, "question", sess, s3Bucket)
	if err != nil {
		log.Fatalf("Failed to upload questions: %v", err)
	}
	log.Printf("Questions successfully uploaded! Location: %v, version: %v, upload: %v\n", uploadResult.Location, uploadResult.VersionID, uploadResult.UploadID)

	uploadResult, err = writeToS3(answerVertices, "question", sess, s3Bucket)
	if err != nil {
		log.Fatalf("Failed to upload answers: %v", err)
	}
	log.Printf("Answers successfully uploaded! Location: %v, version: %v, upload: %v\n", uploadResult.Location, uploadResult.VersionID, uploadResult.UploadID)

	uploadResult, err = writeToS3(peopleVertices, "question", sess, s3Bucket)
	if err != nil {
		log.Fatalf("Failed to upload people: %v", err)
	}
	log.Printf("People successfully uploaded! Location: %v, version: %v, upload: %v\n", uploadResult.Location, uploadResult.VersionID, uploadResult.UploadID)

	uploadResult, err = writeToS3(edges, "question", sess, s3Bucket)
	if err != nil {
		log.Fatalf("Failed to upload edges: %v", err)
	}
	log.Printf("Eges successfully uploaded! Location: %v, version: %v, upload: %v\n", uploadResult.Location, uploadResult.VersionID, uploadResult.UploadID)

	return fmt.Sprintf("Hello 4 %s!", name.Name), nil
}

func writeToS3(records [][]string, key string, sess *session.Session, bucket string) (*s3manager.UploadOutput, error) {
	log.Printf("Write %d records, key: %v, bucket: %v\n", len(records), key, bucket)
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	defer writer.Flush()

	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			log.Fatalf("Failed to write record %v. Err: %v\n", record, err.Error())
		}
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return result, err
}
