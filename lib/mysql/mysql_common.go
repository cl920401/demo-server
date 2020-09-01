package mysql

import (
	"database/sql"

	"demo-server/lib/log"
)

// Commit : 封装提交方法，出现错误时输出日志
func Commit(tx *sql.Tx) {
	if tx == nil {
		return
	}
	if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
		log.Error(err)
	}
}

// Rollback : 封装回滚方法，出现错误时输出日志
func Rollback(tx *sql.Tx) {
	if tx == nil {
		return
	}
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Error(err)
	}
}

// CloseRows : 封装关闭连接方法，出现错误时输出日志
func CloseRows(rows *sql.Rows) {
	if rows == nil {
		return
	}
	if err := rows.Close(); err != nil {
		log.Error("CloseRows error:", err)
	}
}
