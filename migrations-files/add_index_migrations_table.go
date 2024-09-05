package migrations_files

type AddIndexNameMigrations struct{}

func (m *AddIndexNameMigrations) Up() string {
	return `
		ALTER TABLE 
		    migrations
		ADD CONSTRAINT idx_name_migration 
		    UNIQUE (name);
	`
}

func (m *AddIndexNameMigrations) Down() string {
	return `
		ALTER TABLE
			migrations
		DROP INDEX idx_name_migration;
	`
}
