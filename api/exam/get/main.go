package main

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

type EventExamGetRequest struct {
	ExamID int `json:"exam_id,omitempty"`
}

type EventExamGetResponse struct {
	Exams []model.Exam
}

func (e EventExamGetResponse) Encode() (string, error) {
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

	body := EventExamGetRequest{}
	_ = lib.GetBody(req, &body)

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("public.exam")

	if body.ExamID != 0 {
		stmt = stmt.Where(sq.Eq{"id": body.ExamID})
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	exams := make([]model.Exam, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &exams, query, args...)
	if err != nil {
		return "", err
	}

	return EventExamGetResponse{
		Exams: exams,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
