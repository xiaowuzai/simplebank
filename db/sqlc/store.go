package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 包含了所有数据库相关的接口，包括事务
type Store interface {
	Querier // 嵌入所有生成的接口
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error)
}

// Store 提供了数据库查询方式并且支持事务
type SQLStore struct {
	*Queries // 嵌入而非继承
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx 执行事务
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // 使用默认的隔离级别
	if err != nil {
		return nil
	}

	query := New(tx)
	if err := fn(query); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}
	return tx.Commit()
}
