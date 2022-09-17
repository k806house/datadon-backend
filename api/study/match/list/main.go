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

type EventStudyMatchListRequest struct {
	StudyID int `json:"study_id,omitempty"`
}

type EventStudyMatchListResponse struct {
	Exams []model.Exam `json:"exams"`
}

const (
	Wait     = "wait"
	Approved = "approved"
	Declined = "declined"
)

func (e EventStudyMatchListResponse) Encode() (string, error) {
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

	body := EventStudyMatchListRequest{}
	_ = lib.GetBody(req, &body)

	if body.StudyID <= 0 {
		return "", errors.New("Invalid study id")
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("exam_id").From("match").
		Where(sq.Eq{"study_id": body.StudyID}).Where(sq.Eq{"match.user": Approved})

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
		Select("*").From("exam").
		Where(sq.Eq{"id": examsIDs})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	exams := make([]model.Exam, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &exams, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyMatchListResponse{
		Exams: exams,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
