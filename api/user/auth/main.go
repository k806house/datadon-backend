package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
	"github.com/k806house/datadon-backend/model"
	"github.com/pkg/errors"
)

type EventAuthRequest struct {
	UserID int `json:"user_id"`
}

type EventAuthResponse struct {
	AuthorizationToken string `json:"authorization_token"`
}

func (e EventAuthResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	log.Printf("Input request: %v", req)

	auth := EventAuthRequest{}
	err := lib.GetBody(req, &auth)
	if err != nil {
		return "", err
	}

	user := model.User{}
	err = lib.GetDB(ctx).Get(&user, "SELECT * FROM public.user WHERE id = $1", auth.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("user not found")
	}

	if err != nil {
		return "", err
	}

	authToken := randToken()
	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("public.session").
		Columns("user_id", "auth_key").
		Values(user.ID, authToken)

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}
	_, err = lib.GetDB(ctx).Query(query, args...)

	if err != nil {
		return "", err
	}

	return EventAuthResponse{
		AuthorizationToken: authToken,
	}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
