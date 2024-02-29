package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store 包含了所有数据库相关的接口，包括事务
type Store interface {
	Querier // 嵌入所有生成的接口
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// Store 提供了数据库查询方式并且支持事务
type SQLStore struct {
	*Queries // 嵌入而非继承
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
