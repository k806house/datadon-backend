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

type EventStudyGetRequest struct {
	StudyID int `json:"exam_id,omitempty"`
}

type EventStudyGetResponse struct {
	Studies []model.Study
}

func (e EventStudyGetResponse) Encode() (string, error) {
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

	body := EventStudyGetRequest{}
	_ = lib.GetBody(req, &body)

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("public.study")

	if body.StudyID != 0 {
		stmt = stmt.Where(sq.Eq{"id": body.StudyID})
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	studies := make([]model.Study, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &studies, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyGetResponse{
		Studies: studies,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
