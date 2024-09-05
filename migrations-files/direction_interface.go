package migrations_files

type Direction interface {
	Up() string
	Down() string
}
