package model

import "time"

type Exam struct {
	ID          int       `json:"id,omitempty" db:"id"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	FileList    string    `json:"file_list" db:"file_list"`
	Name        string    `json:"name"  db:"name"`
	Description string    `json:"description"  db:"description"`
	UserID      int       `json:"user_id"  db:"user_id"`
}
