package main

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
	"github.com/k806house/datadon-backend/repository/exam"
)

type EventStudyDownloadRequest struct {
	StudyID int `json:"study_id,omitempty"`
}

type EventStudyDownloadResponse struct {
	FileList []string `json:"file_list"`
}

const (
	Wait     = "wait"
	Approved = "approved"
	Declined = "declined"
)

func (e EventStudyDownloadResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, err := lib.GetAuthUserID(ctx, req)
	if err != nil {
		return "", err
	}

	if userID == -1 {
		return "", errors.New("unauthorized")
	}

	body := EventStudyDownloadRequest{}
	_ = lib.GetBody(req, &body)

	if body.StudyID <= 0 {
		return "", errors.New("invalid study id")
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("public.study").
		Where(sq.Eq{"id": body.StudyID}).Where(sq.Eq{"user_id": userID})

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	study := model.Study{}
	err = lib.GetDB(ctx).GetContext(ctx, &study, query, args...)
	if err != nil {
		return "", err
	}

	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("exam_id").From("public.match").
		Where(sq.Eq{"study_id": body.StudyID}).Where(sq.Eq{"match.user": Approved}).Where(sq.Eq{"match.researcher": Approved})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	examsIDs := make([]int, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &examsIDs, query, args...)
	if err != nil {
		return "", err
	}

	allFiles := make([]string, 0)
	for _, examID := range examsIDs {
		curFiles, err := exam.GetExamFiles(ctx, examID)
		if err != nil {
			return "", err
		}
		allFiles = append(allFiles, curFiles...)
	}

	return EventStudyDownloadResponse{
		FileList: allFiles,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
