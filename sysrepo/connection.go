package sysrepo

/*
#include "lib.h"
*/
import "C"

import "fmt"

// Connection к Sysrepo, удобная обёртка для работы с подключением на уровне Golang
type Connection struct {
	conn *C.sr_conn_ctx_t
}

func NewConnection() (*Connection, error) {
	core := Core{}

	var conn *C.sr_conn_ctx_t
	if err := core.Connect(C.uint(C.SR_CONN_DEFAULT), &conn); err != nil {
		return nil, fmt.Errorf("failed to connect to sysrepo: %w", err)
	}

	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) Disconnect() error {
	core := Core{}

	return core.Disconnect(c.conn)
}
