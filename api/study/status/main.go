package main

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
	"github.com/rs/zerolog/log"
)

type EventStudyStatusRequest struct {
	StudyID int `json:"study_id,omitempty"`
}

type EventStudyStatusResponse struct {
	Found                     int `json:"found"`
	WaitingUserDecision       int `json:"waiting_user_decision"`
	WaitingResearcherDecision int `json:"waiting_researcher_decision"`
	RejectedByUser            int `json:"rejected_by_user"`
	RejectedByResearcher      int `json:"rejected_by_researcher"`
	Ready                     int `json:"ready"`
}

const (
	Wait     = "wait"
	Approved = "approved"
	Declined = "declined"
)

func (e EventStudyStatusResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func Intersect[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		if containsGeneric(b, v) {
			set = append(set, v)
		}
	}

	return set
}

func containsGeneric[T comparable](b []T, e T) bool {
	for _, v := range b {
		if v == e {
			return true
		}
	}
	return false
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, err := lib.GetAuthUserID(ctx, req)
	if err != nil {
		return "", err
	}

	if userID == -1 {
		return "", errors.New("Unauthorized")
	}

	body := EventStudyStatusRequest{}
	_ = lib.GetBody(req, &body)

	if body.StudyID <= 0 {
		return "", errors.New("Invalid study id")
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

	tagList := make([]string, 0)
	for _, tag := range study.Tags {
		tagList = append(tagList, tag.Name)
	}

	if len(tagList) == 0 {
		return EventStudyStatusResponse{}.Encode()
	}

	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").From("exam")

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}
	exams := make([]model.Study, 0)
	err = lib.GetDB(ctx).SelectContext(ctx, &exams, query, args...)
	if err != nil {
		return "", err
	}

	for _, ex := range exams {
		exTagList := make([]string, 0)
		for _, tag := range ex.Tags {
			exTagList = append(exTagList, tag.Name)
		}

		intersect := Intersect(tagList, exTagList)
		if len(intersect) == len(tagList) {
			stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
				Insert("match").Columns("study_id", "exam_id").Values(study.ID, ex.ID).Suffix("ON CONFLICT DO NOTHING")

			query, args, err = stmt.ToSql()
			if err != nil {
				return "", err
			}

			_, err = lib.GetDB(ctx).ExecContext(ctx, query, args...)
			if err != nil {
				return "", err
			}
		}
	}

	// Found
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"study_id": body.StudyID})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	found := 0
	err = lib.GetDB(ctx).GetContext(ctx, &found, query, args...)
	if err != nil {
		return "", err
	}

	// WaitingUserDecision
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"match.study_id": body.StudyID}).Where(sq.Eq{"match.user": Wait})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	waitingUserDecision := 0
	err = lib.GetDB(ctx).GetContext(ctx, &waitingUserDecision, query, args...)
	if err != nil {
		return "", err
	}

	// WaitingResearcherDecision
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"match.study_id": body.StudyID}).Where(sq.Eq{"match.user": Approved}).Where(sq.Eq{"match.researcher": Wait})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	log.Print(query)
	waitingResearcherDecision := 0
	err = lib.GetDB(ctx).GetContext(ctx, &waitingResearcherDecision, query, args...)
	if err != nil {
		return "", err
	}

	// RejectedByUser
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"match.study_id": body.StudyID}).Where(sq.Eq{"match.user": Declined})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	rejectedByUser := 0
	err = lib.GetDB(ctx).GetContext(ctx, &rejectedByUser, query, args...)
	if err != nil {
		return "", err
	}

	// RejectedByResearcher
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"match.study_id": body.StudyID}).Where(sq.Eq{"match.researcher": Declined})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	rejectedByResearcher := 0
	err = lib.GetDB(ctx).GetContext(ctx, &rejectedByResearcher, query, args...)
	if err != nil {
		return "", err
	}
	// Ready
	stmt = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").From("match").
		Where(sq.Eq{"match.study_id": body.StudyID}).Where(sq.Eq{"match.researcher": Approved}).Where(sq.Eq{"match.user": Approved})

	query, args, err = stmt.ToSql()
	if err != nil {
		return "", err
	}

	ready := 0
	err = lib.GetDB(ctx).GetContext(ctx, &ready, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyStatusResponse{
		Found:                     found,
		WaitingUserDecision:       waitingUserDecision,
		WaitingResearcherDecision: waitingResearcherDecision,
		RejectedByUser:            rejectedByUser,
		RejectedByResearcher:      rejectedByResearcher,
		Ready:                     ready,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
