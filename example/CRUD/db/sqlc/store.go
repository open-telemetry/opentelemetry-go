package db

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
}

func NewStore(dbtx DBTX) Store {
	return &SQLStore{
		Queries: New(dbtx),
	}
}
