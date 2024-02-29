package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// postgres 会返回固定的错误代码: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	ForeignKeyViolationCode = "23503"
	UniqueViolationCode     = "23505"
)

// 数据未找到错误
var ErrRecordNotFound = pgx.ErrNoRows

// 唯一键错误
var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolationCode,
}

// 外键错误
var ErrForeignKeyViolation = &pgconn.PgError{
	Code: ForeignKeyViolationCode,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
