package database

const (
	SortByAsc  = "ASC"
	SortByDesc = "DESC"
)

type FindOptions struct {
	Limit  uint64
	Offset uint64
}

type Field struct {
	Name  string
	Value any
}

type Scanner interface {
	Scan(dest ...any) error
}
