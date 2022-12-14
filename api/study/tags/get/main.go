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

type EventStudyTagsGetRequest struct {
	Study int `json:"study_id"`
}

type EventStudyTagsGetResponse struct {
	Tags model.TagList `json:"tags"`
}

func (e EventStudyTagsGetResponse) Encode() (string, error) {
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

	body := EventStudyTagsGetRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	if body.Study == 0 {
		return "", errors.New("invalid request")
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("tags").From("public.study").
		Where(sq.Eq{"id": body.Study}).Where(sq.Eq{"user_id": userID})

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	tags := make(model.TagList, 0)
	err = lib.GetDB(ctx).GetContext(ctx, &tags, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("no study found")
	}

	if err != nil {
		return "", err
	}

	return EventStudyTagsGetResponse{Tags: tags}.Encode()

}

func main() {
	lambda.Start(HandleRequest)
}
