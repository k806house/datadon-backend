package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
)

type EventLogoutUserRequest struct {
	AuthorizationToken string `json:"auth_token"`
}

type EventLogoutUserResponse struct {
}

func (e EventLogoutUserResponse) Encode() (string, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

func HandleRequest(ctx context.Context, req map[string]interface{}) (string, error) {
	userID, err := lib.GetAuthUserID(ctx, req)
	if err != nil {
		return "", err
	}

	token, err := lib.GetAuthToken(req)
	if err != nil {
		return "", err
	}

	if userID != -1 {
		_, err := lib.GetDB(ctx).Exec("UPDATE public.session SET expired = TRUE WHERE auth_key = $1", token)
		if err != nil {
			return "", err
		}
	}

	return EventLogoutUserResponse{}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
