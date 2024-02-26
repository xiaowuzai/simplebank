package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreateUser func(user User) error
}

type CreateUserTxResult struct {
	User User
}

// CreateUserTx 转移金额事务
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreateUser(result.User)
	})

	return result.User, err
}
