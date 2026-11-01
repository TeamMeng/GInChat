package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type BizError struct {
	Msg string
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(msg string) error {
	return &BizError{Msg: msg}
}

func IsDuplicateKey(err error) bool {
	var sqlErr *pgconn.PgError
	return errors.As(err, &sqlErr) && sqlErr.Code == "23505"
}
