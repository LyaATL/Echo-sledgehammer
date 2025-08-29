package models

import "time"

// ClientBan Ban represents a single global ban entry.
type ClientBan struct {
	ID          int       `db:"id" json:"id"`
	Character   string    `db:"character" json:"character"`
	World       string    `db:"world" json:"world"`
	LodestoneID string    `db:"lodestone_id" json:"lodestoneId,omitempty"`
	Reason      string    `db:"reason" json:"reason"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	SubmittedBy string    `db:"submitted_by" json:"submittedBy"`
	ApprovedBy  string    `db:"approved_by" json:"approvedBy,omitempty"`
	ApprovedAt  time.Time `db:"approved_at" json:"approvedAt,omitempty"`
	Status      string    `db:"status" json:"status"`
}
