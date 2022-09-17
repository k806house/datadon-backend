package main

import (
	"context"
	"encoding/json"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/jackc/pgx/v4"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

type EventGetUserRequest struct {
}

type EventGetUserResponse struct {
	User []model.User `json:"user"`
}

func (e EventGetUserResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, _ := lib.GetAuthUserID(ctx, req)

	if userID == -1 {
		log.Println("Start processing request")

		stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("public.user")

		query, args, err := stmt.ToSql()
		if err != nil {
			return "", err
		}

		users := make([]model.User, 0)
		err = lib.GetDB(ctx).SelectContext(ctx, &users, query, args...)
		if err != nil {
			return "", err
		}

		return EventGetUserResponse{User: users}.Encode()
	} else {
		stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("public.user").Where(sq.Eq{"id": userID})
		query, args, err := stmt.ToSql()
		if err != nil {
			return "", err
		}

		users := make([]model.User, 0)
		err = lib.GetDB(ctx).SelectContext(ctx, &users, query, args...)
		if err != nil {
			return "", err
		}

		return EventGetUserResponse{User: users}.Encode()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
