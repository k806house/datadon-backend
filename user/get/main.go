package main

import (
	"context"
	"encoding/json"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
)

var conn *sqlx.DB

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	log.Println("Start processing request")

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("public.user")

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	users := make([]model.User, 0)
	err = conn.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(users)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {
	ctx := context.Background()

	var err error
	conn, err = lib.ConnectToDB(ctx)
	if err != nil {
		panic(err)
	}

	lambda.Start(HandleRequest)
}
