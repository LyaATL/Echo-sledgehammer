package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sledgehammer.echo-mesh.com/internal/models"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InitSchema() error {
	schema := `
    CREATE TABLE IF NOT EXISTS client_bans (
        id TEXT PRIMARY KEY,
        character TEXT NOT NULL,
        world TEXT NOT NULL,
        lodestone_id TEXT,
        reason TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        submitted_by TEXT NOT NULL
    );

	CREATE TABLE IF NOT EXISTS file_bans (
	    ID TEXT PRIMARY KEY,
	    hash TEXT NOT NULL,
	    signature TEXT,
        reason TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        submitted_by TEXT NOT NULL
	)`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) AddClientBan(b models.ClientBan) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now().UTC()

	_, err := s.db.NamedExec(`
        INSERT INTO client_bans (id, character, world, lodestone_id, reason, created_at, submitted_by)
        VALUES (:id, :character, :world, :lodestone_id, :reason, :created_at, :submitted_by)
    `, b)
	return err
}

func (s *Store) AddFileBan(b models.FileBan) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now().UTC()

	_, err := s.db.NamedExec(`
        INSERT INTO file_bans (id, hash, signature, reason, created_at, submitted_by)
        VALUES (:id, :hash, :signature, :reason, :created_at, :submitted_by)
    `, b)
	return err
}

func (s *Store) DoesClientBanExist(character string, world string) (bool, error) {
	var count int
	err := s.db.Get(&count, `
        SELECT COUNT(*) FROM client_bans
        WHERE character = ? AND world = ?
    `, character, world)
	return count > 0, err
}

func (s *Store) GetPlayerBanInfo(character string, world string) (models.ClientBan, error) {
	var ban models.ClientBan

	err := s.db.Get(&ban, `
    		SELECT character, world, reason, created_at 
    		FROM clientbans 
    		WHERE character = ? AND world = ?`,
		character, world)

	return ban, err
}
