package model

import "time"

type Exam struct {
	ID          int       `json:"id,omitempty" db:"id"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	Name        string    `json:"name"  db:"name"`
	Description string    `json:"description"  db:"description"`
	UserID      int       `json:"user_id"  db:"user_id"`
	Tags        TagList   `json:"tags" db:"tags"`
}
