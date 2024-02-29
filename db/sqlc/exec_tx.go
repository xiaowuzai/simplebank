package db

import (
	"context"
	"fmt"
)

// execTx 执行事务
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx) // 使用默认的隔离级别
	if err != nil {
		return nil
	}

	query := New(tx)
	if err := fn(query); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}
	return tx.Commit(ctx)
}
