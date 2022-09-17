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

type EventStudyCreateRequest struct {
	model.Study
}

type EventStudyCreateResponse struct {
	StudyID int `json:"study_id"`
}

func (e EventStudyCreateResponse) Encode() (string, error) {
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

	body := EventStudyCreateRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("public.study").
		Columns("user_id", "name", "description", "created_at", "tags").
		Values(userID, body.Name, body.Description, time.Now().UTC(), body.Tags).Suffix("RETURNING id")

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	studyID := 0
	err = lib.GetDB(ctx).Get(&studyID, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyCreateResponse{
		StudyID: studyID,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
