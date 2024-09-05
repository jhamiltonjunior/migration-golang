package migrations_files

type CreateUserTable struct {
}

func (u *CreateUserTable) Up() string {
	return `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL
		)`
}

func (u *CreateUserTable) Down() string {
	return `DROP TABLE users`
}
