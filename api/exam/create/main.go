package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

type EventExamCreateRequest struct {
	model.Exam
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
		Columns("user_id", "name", "description", "created_at", "file_list").
		Values(userID, body.Name, body.Description, time.Now().UTC(), "")

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	_, err = lib.GetDB(ctx).Query(query, args...)
	if err != nil {
		return "", err
	}

	return EventExamCreateResponse{
		Exam: model.Exam{},
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
