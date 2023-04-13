package pgstore

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pershin-daniil/TimeSlots/pkg/models"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"strings"
)

//go:embed migrations
var migrations embed.FS

type Store struct {
	log *logrus.Entry
	db  *sqlx.DB
}

func New(ctx context.Context, log *logrus.Logger, dsn string) (*Store, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Store{
		log: log.WithField("module", "pgstore"),
		db:  db,
	}, nil
}

func (s *Store) Migrate(direction migrate.MigrationDirection) error {
	assetDir := func() func(string) ([]string, error) {
		return func(path string) ([]string, error) {
			dirEntry, er := migrations.ReadDir(path)
			if er != nil {
				return nil, er
			}
			entries := make([]string, 0)
			for _, e := range dirEntry {
				entries = append(entries, e.Name())
			}

			return entries, nil
		}
	}()
	asset := migrate.AssetMigrationSource{
		Asset:    migrations.ReadFile,
		AssetDir: assetDir,
		Dir:      "migrations",
	}
	_, err := migrate.Exec(s.db.DB, "postgres", asset, direction)
	if err != nil {
		return fmt.Errorf("err migrating: %w", err)
	}
	s.log.Infof("migration success")
	return nil
}

func (s *Store) User(ctx context.Context, userID int64) (models.User, error) {
	query := `
SELECT id, last_name, first_name, status, created_at, updated_at
FROM users
WHERE id = $1`
	var result models.User
	if err := s.db.GetContext(ctx, &result, query, userID); err != nil {
		return models.User{}, fmt.Errorf("get user faild: %w", err)
	}
	return result, nil
}

func (s *Store) CreateUser(ctx context.Context, user models.UserRequest) (models.User, error) {
	query := `
INSERT INTO users (id, last_name, first_name)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO NOTHING
RETURNING id, last_name, first_name, status, created_at, updated_at`
	var newUser models.User
	if err := s.db.GetContext(ctx, &newUser, query, user.ID, user.LastName, user.FirstName); err != nil {
		return models.User{}, fmt.Errorf("create user faild: %w", err)
	}
	return newUser, nil
}

func (s *Store) UpdateUser(ctx context.Context, user models.UserRequest) (models.User, error) {
	query := `
UPDATE users
	SET last_name = $2,
		first_name = $3,
	    status = $4
WHERE id = $1
RETURNING id, last_name, first_name, status;`
	var updatedUser models.User
	if err := s.db.GetContext(ctx, &updatedUser, query, user.ID, user.LastName, user.FirstName, user.Status); err != nil {
		return models.User{}, fmt.Errorf("update user faild: %w", err)
	}
	return updatedUser, nil
}

func (s *Store) Status(ctx context.Context, userID int64) (string, error) {
	query := `
SELECT status
FROM users
WHERE id = $1`
	var status string
	if err := s.db.GetContext(ctx, &status, query, userID); err != nil {
		return "", fmt.Errorf("get status faild: %w", err)
	}
	return status, nil
}

//func (s *Store) NewMsg(ctx context.Context, UserID int64, Msg int) (models.Msg, error) {
//	query := `
//INSERT INTO messages (user_id, msg_id)
//VALUES ($1, $2)
//`
//}

func (s *Store) ResetTables(ctx context.Context, tables []string) error {
	_, err := s.db.ExecContext(ctx, `TRUNCATE TABLE`+` `+strings.Join(tables, `, `))
	for _, table := range tables {
		_, err = s.db.ExecContext(ctx, fmt.Sprintf(`ALTER SEQUENCE %s_id_seq RESTART`, table))
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Store) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *Store) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Store) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
