package models

import "time"

type FileBan struct {
	ID          int       `db:"id" json:"id"`
	Filename    string    `db:"filename" json:"filename"`
	Hash        string    `db:"hash" json:"hash"`
	Signature   string    `db:"signature" json:"signature,omitempty"`
	Reason      string    `db:"reason" json:"reason"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	SubmittedBy string    `db:"submitted_by" json:"submittedBy"`
	ApprovedBy  string    `db:"approved_by" json:"approvedBy,omitempty"`
	ApprovedAt  time.Time `db:"approved_at" json:"approvedAt,omitempty"`
	Status      string    `db:"status" json:"status"`
}
