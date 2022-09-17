package approve

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k806house/datadon-backend/lib"
)

type EventStudyMatchResponseRequest struct {
	StudyID  int    `json:"study_id,omitempty"`
	ExamID   int    `json:"exam_id,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Response string `json:"action,omitempty"`
}

type EventStudyMatchResponseResponse struct {
}

const (
	Wait     = "wait"
	Approved = "approved"
	Declined = "declined"

	User       = "user"
	Researcher = "researcher"
)

func (e EventStudyMatchResponseResponse) Encode() (string, error) {
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

	body := EventStudyMatchResponseRequest{}
	err = lib.GetBody(req, &body)
	if err != nil {
		return "", err
	}

	if body.StudyID <= 0 || body.ExamID <= 0 {
		return "", errors.New("invalid study id or exam id")
	}

	if body.Response != Approved && body.Response != Declined {
		return "", errors.New("response kind must be " + Approved + " or " + Declined)
	}

	stmt := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("match")
	if body.Kind == User {
		stmt = stmt.
			Set("user", body.Response).
			Where(sq.Eq{"study_id": body.StudyID}, sq.Eq{"exam_id": body.ExamID})
	} else if body.Kind == Researcher {
		stmt = stmt.
			Set("researcher", body.Response).
			Where(sq.Eq{"study_id": body.StudyID}, sq.Eq{"exam_id": body.ExamID})
	} else {
		return "", errors.New("invalid kind. must be " + Researcher + " or " + User)
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return "", err
	}

	_, err = lib.GetDB(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}

	return EventStudyMatchResponseResponse{}.Encode()
}

func main() {
	lambda.Start(HandleRequest)
}
