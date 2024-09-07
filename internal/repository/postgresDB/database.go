package postgresDB

import (
	"github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pgxPool      *pgxpool.Pool
	queryBuilder squirrel.StatementBuilderType
}

func NewStorage(pgxPool *pgxpool.Pool) *Storage {
	return &Storage{
		pgxPool:      pgxPool,
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
