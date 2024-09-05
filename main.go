package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	migrationsfiles "github.com/jhamiltonjunior/migration-golang/migrations-files"
)

func main() {
	NewMigrationConnection()
}

func NewMigrationConnection() *MigrationConnection {
	migrationConnection := &MigrationConnection{}
	migrationConnection.Connect()
	migrationConnection.Migrate()

	toBeExecuted := map[string]map[string]migrationsfiles.Direction{
		"up": {
			"create_user_table":            &migrationsfiles.CreateUserTable{},
			"create_index_name_migrations": &migrationsfiles.AddIndexNameMigrations{},
		},
		"down": {
			"drop_user_table": &migrationsfiles.CreateUserTable{},
		},
	}

	for direction, migrations := range toBeExecuted {
		for identifier, migration := range migrations {
			if migrationConnection.HasMigration(identifier) {
				continue
			}

			if direction == "up" {
				query := migration.Up()
				migrationConnection.Commit(query)
				migrationConnection.CreateMigration(identifier)
			} else if direction == "down" {
				query := migration.Down()
				migrationConnection.Rollback(query)
				migrationConnection.CreateMigration(identifier)
			}
		}
	}

	return migrationConnection
}

type MigrationConnection struct {
	conn *sql.DB
}

func (m *MigrationConnection) Connect() {
	db, err := sql.Open("mysql", "root:0000@tcp(localhost:3306)/test")
	if err != nil {
		panic(err)
	}
	m.conn = db
}

func (m *MigrationConnection) Close() {
	err := m.conn.Close()
	if err != nil {
		panic(err)
	}
}

func (m *MigrationConnection) CreateMigrateTable() {
	_, err := m.conn.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
			    id INT AUTO_INCREMENT PRIMARY KEY,
			    name VARCHAR(255) NOT NULL,
			    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
	`)
	if err != nil {
		panic(err)
	}
}

func (m *MigrationConnection) CreateMigration(name string) {
	_, err := m.conn.Exec("INSERT INTO migrations (name) VALUES (?)", name)
	if err != nil {
		panic(err)
	}
}

func (m *MigrationConnection) GetMigrations() []string {
	rows, err := m.conn.Query("SELECT name FROM migrations")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		migrations = append(migrations, name)
	}

	return migrations
}

func (m *MigrationConnection) HasMigration(name string) bool {
	rows, err := m.conn.Query("SELECT name FROM migrations WHERE name = ?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return rows.Next()
}

func (m *MigrationConnection) RollbackMigration(name string) {
	_, err := m.conn.Exec("DELETE FROM migrations WHERE name = ?", name)
	if err != nil {
		panic(err)
	}
}

func (m *MigrationConnection) Migrate() {
	m.CreateMigrateTable()

	migrations := m.GetMigrations()
	if len(migrations) == 0 {
		m.CreateMigration("initial")
	}
}

func (m *MigrationConnection) Rollback(query string) {
	_, err := m.conn.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (m *MigrationConnection) Commit(query string) {
	_, err := m.conn.Exec(query)
	if err != nil {
		panic(err)
	}
}
