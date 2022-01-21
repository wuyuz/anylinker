package db

import (
	"context"
	"testing"
)

func TestNewDb(t *testing.T) {

	err := NewDb(Drivename("sqlite3"),
		Dsn("sqlite3.db"),
		MaxIdleConnection(10),
		MaxQueryTime(3),
		MaxQueryTime(3),
		MaxOpenConnection(3),
	)
	if err != nil {
		t.Fatalf("NewDb Err: %v", err)
	}
	conn, err := GetConn(context.Background())
	if err != nil {
		t.Fatalf("Get Conn Err: %v", err)
	}
	conn.Close()
	//_ = os.Remove("sqlite3.db")
}


func TestNewMySQLDb(t *testing.T) {
	err := NewDb(Drivename("mysql"),
		Dsn("root:root@tcp(127.0.0.1:3306)/anylinker?charset=utf8mb4&parseTime=True&loc=Local"),
		MaxIdleConnection(10),
		MaxQueryTime(3),
		MaxQueryTime(3),
		MaxOpenConnection(3),
	)
	if err != nil {
		t.Fatalf("NewDb Err: %v", err)
	}
	conn, err := GetConn(context.Background())
	if err != nil {
		t.Fatalf("Get Conn Err: %v", err)
	}
	conn.Close()
	//_ = os.Remove("sqlite3.db")
}
