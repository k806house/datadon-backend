package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

type EventExamCreateRequest struct {
	model.Exam
	Files []struct {
		Name    string `json:"name"`
		TmpName string `json:"tmp_name"`
	} `json:"files"`
}

type EventExamCreateResponse struct {
	model.Exam
}

func (e EventExamCreateResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, err := lib.GetAuthUserID(ctx, req)
	if err != nil {
		return "", err
	}

	if userID == -1 {
		return "", errors.New("Unauthorized")
	}

	body := EventExamCreateRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("public.exam").
		Columns("user_id", "name", "description", "created_at").
		Values(userID, body.Name, body.Description, time.Now().UTC()).Suffix("RETURNING id")

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	examID := 0
	err = lib.GetDB(ctx).Get(&examID, query, args...)
	if err != nil {
		return "", err
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

	awsS3Client := s3.NewFromConfig(cfg)

	for _, f := range body.Files {
		_, err = awsS3Client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String("datadon-data"),
			CopySource: aws.String(fmt.Sprintf("/datadon-data/tmp/%s", f.TmpName)),
			Key:        aws.String(fmt.Sprintf("exam/%d/%s", examID, f.Name)),
		})
		if err != nil {
			return "", err
		}
		_, _ = awsS3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String("datadon-data"),
			Key:    aws.String(fmt.Sprintf("tmp/%s", f.TmpName)),
		})
	}

	return EventExamCreateResponse{
		Exam: model.Exam{},
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
