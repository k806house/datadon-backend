package lib

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetBody(req map[string]interface{}, body interface{}) error {
	if val, ok := req["body"]; ok {
		if valStr, ok := val.(string); ok {
			err := json.Unmarshal([]byte(valStr), &body)
			if err != nil {
				return err
			}
		} else {
			return errors.New("wrong body format")
		}
	} else {
		return errors.New("body not found")
	}

	return nil
}

func GetAuthToken(req map[string]interface{}) (string, error) {
	authKey := ""
	if headers, ok := req["headers"]; ok {
		if headersMap, ok := headers.(map[string]interface{}); ok {
			log.Warn().Msgf("headers: %v", headersMap)
			if auth, ok := headersMap["authorization"]; ok {
				if authStr, ok := auth.(string); ok {
					authKey = authStr
				} else {
					log.Warn().Msg("wrong auth format")
					return "", errors.New("wrong auth format")
				}
			} else {
				log.Warn().Msg("Authorization header not found")
				return "", errors.New("auth not found")
			}
		} else {
			log.Warn().Msg("headers is not map[string]interface{}")
			return "", errors.New("wrong headers format")
		}
	} else {
		log.Warn().Msg("headers not found")
		return "", errors.New("headers not found")
	}

	return authKey, nil
}

func GetAuthUserID(ctx context.Context, req map[string]interface{}) (int, error) {
	authKey, err := GetAuthToken(req)
	if err != nil {
		return -1, err
	}

	userID := -1
	err = GetDB(ctx).Get(&userID, "SELECT user_id FROM public.session WHERE auth_key = $1 AND NOT expired", authKey)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, errors.New("no user found for auth key")
	}

	return userID, err
}
