package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

type EventTagGetRequest struct {
}

type EventTagGetResponse struct {
	Tags model.TagList `json:"tags"`
}

func (e EventTagGetResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("public.tag")

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	tags := make(model.TagList, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &tags, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("no exam found")
	}

	if err != nil {
		return "", err
	}

	return EventTagGetResponse{Tags: tags}.Encode()

}

func main() {
	lambda.Start(HandleRequest)
}
