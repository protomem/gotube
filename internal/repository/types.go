package repository

type FindOptions struct {
	Limit  uint64
	Offset uint64
}

type Scanable interface {
	Scan(dest ...any) error
}
