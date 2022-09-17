package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/k806house/datadon-backend/lib"
)

type EventExamUploadResponse struct {
	UploadLink  string `json:"upload_link"`
	TmpFileName string `json:"tmp_file_name"`
}

func (e EventExamUploadResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func getRandomName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, err := lib.GetAuthUserID(ctx, req)
	if err != nil {
		return "", err
	}

	if userID == -1 {
		return "", errors.New("Unauthorized")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-west-1"),
	)
	if err != nil {
		return "", err
	}

	if err != nil {
		log.Fatalf("failed to create AWS session, %v", err)
	}

	fileName := getRandomName(32)
	awsS3Client := s3.NewFromConfig(cfg)
	pr := s3.NewPresignClient(awsS3Client)
	prReq, err := pr.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("datadon-data"),
		Key:    aws.String(fmt.Sprintf("tmp/%s", fileName)),
	})
	if err != nil {
		return "", err
	}

	return EventExamUploadResponse{
		UploadLink:  prReq.URL,
		TmpFileName: fileName,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
