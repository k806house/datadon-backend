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

type EventExamTagsSetRequest struct {
	ExamID int           `json:"exam_id"`
	Tags   model.TagList `json:"tags"`
}

type EventExamTagsSetResponse struct {
}

func (e EventExamTagsSetResponse) Encode() (string, error) {
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

	body := EventExamTagsSetRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	if body.ExamID == 0 {
		return "", errors.New("invalid request")
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("public.exam").Set("tags", body.Tags).
		Where(sq.Eq{"id": body.ExamID}, sq.Eq{"user_id": userID})

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	res, err := lib.GetDB(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return "", err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if affected == 0 {
		return "", errors.New("no exam found")
	}

	return EventExamTagsSetResponse{}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
