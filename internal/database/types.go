package database

type Field struct {
	Name  string
	Value any
}

type Scanner interface {
	Scan(dest ...any) error
}
