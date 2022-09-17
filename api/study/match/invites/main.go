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

type EventStudyMatchInvitesRequest struct {
}

type EventStudyMatchInvitesResponse struct {
	Studies []model.Study `json:"studies"`
}

const (
	Wait     = "wait"
	Approved = "approved"
	Declined = "declined"
)

func (e EventStudyMatchInvitesResponse) Encode() (string, error) {
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

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id").From("exam").
		Where(sq.Eq{"user_id": userID})

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	examsIDs := make([]int, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &examsIDs, query, args...)
	if err != nil {
		return "", err
	}

	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("study_id").From("match").
		Where(sq.Eq{"exam_id": examsIDs}).Where(sq.Eq{"match.user": Wait})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	studiesIDs := make([]int, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &studiesIDs, query, args...)
	if err != nil {
		return "", err
	}

	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("study").
		Where(sq.Eq{"id": studiesIDs})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	studies := make([]model.Study, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &studies, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyMatchInvitesResponse{
		Studies: studies,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
