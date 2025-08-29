package clientbans

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// InitSchema TODO: Fix obvious sqli
func (s *Store) InitSchema() error {
	schema := `
    CREATE TABLE IF NOT EXISTS clientbans (
        id TEXT PRIMARY KEY,
        character TEXT NOT NULL,
        world TEXT NOT NULL,
        lodestone_id TEXT,
        reason TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        submitted_by TEXT NOT NULL
    );`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) AddClientBan(b ClientBan) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now().UTC()

	_, err := s.db.NamedExec(`
        INSERT INTO clientbans (id, character, world, lodestone_id, reason, created_at, submitted_by)
        VALUES (:id, :character, :world, :lodestone_id, :reason, :created_at, :submitted_by)
    `, b)
	return err
}

func (s *Store) List() ([]ClientBan, error) {
	var bans []ClientBan
	err := s.db.Select(&bans, `SELECT * FROM clientbans ORDER BY created_at DESC`)
	return bans, err
}
