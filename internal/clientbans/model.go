package clientbans

import "time"

// ClientBan Ban represents a single global ban entry.
type ClientBan struct {
	ID          string    `db:"id" json:"id"`
	Character   string    `db:"character" json:"character"`
	World       string    `db:"world" json:"world"`
	LodestoneID string    `db:"lodestone_id" json:"lodestoneId"`
	Reason      string    `db:"reason" json:"reason"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	SubmittedBy string    `db:"submitted_by" json:"submittedBy"`
}
