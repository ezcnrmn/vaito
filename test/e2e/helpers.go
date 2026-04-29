package e2e

import (
	"database/sql"
	"fmt"
	"testing"
)

func seedDatabase(db *sql.DB) error {
	// Пароль для всех тестовых пользователей: 12345678
	query := `
	INSERT INTO users (id, name, email, password_hash, role_id, created_at, version) 
	VALUES 
	(1, 'Admin', 'admin@test.com', '\x24326124313224744f486d706559746a726a774c456c6468357231644f626c62692f68416d6664526d544e386c6e386536783158384f334935423843', (SELECT id FROM roles WHERE name='Administrator'), NOW(), 1),
	(2, 'User', 'user@test.com', '\x24326124313224744f486d706559746a726a774c456c6468357231644f626c62692f68416d6664526d544e386c6e386536783158384f334935423843', (SELECT id FROM roles WHERE name='User'), NOW(), 1);
	SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users));`
	_, err := db.Exec(query)
	return err
}

func cleanTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func cleanUserTable(t *testing.T, db *sql.DB) {
	t.Helper()

	query := `DELETE FROM users WHERE id > 2;`
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
