package database

import (
	"time"

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
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		character TEXT NOT NULL,
		world TEXT NOT NULL,
		lodestone_id TEXT,
		reason TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		submitted_by TEXT NOT NULL,
		approved_by TEXT,
		approved_at TIMESTAMP,
		status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
		 
		FOREIGN KEY (approved_by) REFERENCES management_users(username)
    );

	CREATE TABLE IF NOT EXISTS file_bans (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    filename TEXT NOT NULL,
	    hash TEXT,
	    signature TEXT,
        reason TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        submitted_by TEXT NOT NULL,
	    approved_by TEXT,
		approved_at TIMESTAMP,
		status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
		
		FOREIGN KEY (approved_by) REFERENCES management_users(username)
	);
	
	CREATE TABLE IF NOT EXISTS management_users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    username TEXT NOT NULL UNIQUE,
	    password TEXT NOT NULL,
	    active BOOLEAN NOT NULL DEFAULT TRUE,
	    role TEXT NOT NULL CHECK (role IN ('admin', 'moderator'))
	);

	CREATE INDEX IF NOT EXISTS client_bans_character_world_idx ON client_bans (character, world);
	CREATE INDEX IF NOT EXISTS file_bans_hash_idx ON file_bans (hash);`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) RequestClientBan(b models.ClientBan) error {
	b.CreatedAt = time.Now().UTC()

	_, err := s.db.NamedExec(`
        INSERT INTO client_bans (character, world, lodestone_id, reason, submitted_by)
        VALUES (character, :world, :lodestone_id, :reason, :submitted_by)
    `, b)
	return err
}

func (s *Store) RequestFileBan(b models.FileBan) error {

	_, err := s.db.NamedExec(`
        INSERT INTO file_bans (filename, reason, submitted_by)
        VALUES (:filename, :reason, :submitted_by)
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
    		FROM  client_bans 
    		WHERE character = ? AND world = ?`,
		character, world)

	return ban, err
}

func (s *Store) GetPasswordHashAndRole(username string) (string, string, error) {
	var set struct {
		PasswordHash string `db:"password_hash"`
		Role         string `db:"role"`
	}

	err := s.db.Get(&set, "SELECT password, role FROM management_users WHERE username = ? ", username)
	if err != nil {
		return "", "", err
	}
	return set.PasswordHash, set.Role, nil
}
