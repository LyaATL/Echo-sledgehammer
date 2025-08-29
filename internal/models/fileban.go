package models

import "time"

type FileBan struct {
	ID          string    `db:"id" json:"id"`
	Filename    string    `db:"filename" json:"filename"`
	Hash        string    `db:"hash" json:"hash"`
	Signature   string    `db:"signature" json:"signature"`
	Reason      string    `db:"reason" json:"reason"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	SubmittedBy string    `db:"submitted_by" json:"submittedBy"`
}
