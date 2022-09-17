package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/k806house/datadon-backend/lib"
)

type EventExamUploadRequest struct {
	ExamID   int    `json:"exam_id"`
	FileName string `json:"file_name"`
}

type EventExamUploadResponse struct {
	UploadLink string `json:"upload_link"`
}

func (e EventExamUploadResponse) Encode() (string, error) {
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

	body := EventExamUploadRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	if body.ExamID == 0 || body.FileName == "" {
		return "", errors.New("invalid request")
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("public.exam").
		Where(sq.Eq{"id": body.ExamID}, sq.Eq{"user_id": userID})

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	cnt := 0
	err = lib.GetDB(ctx).GetContext(ctx, &cnt, query, args...)
	if errors.Is(err, sql.ErrNoRows) || cnt == 0 {
		return "", errors.New("no exam found")
	}

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
	pr := s3.NewPresignClient(awsS3Client)
	prReq, err := pr.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("datadon-data"),
		Key:    aws.String(fmt.Sprint("exam/%d/%s", body.ExamID, body.FileName)),
	})
	if err != nil {
		return "", err
	}

	return EventExamUploadResponse{
		UploadLink: prReq.URL,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
