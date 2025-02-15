package repository

import "context"

type TxOrConn interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) (Row, error)
}

type CommandTag interface {
	RowsAffected() int64
}

type Row interface {
	Scan(dest ...interface{}) error
}
