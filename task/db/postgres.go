package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn,error){
	conn,err:=pgx.Connect(
		context.Background(),
		"postgres://postgres:shabin@localhost:5432/taskdb?sslmode=disable",
	)
	if err!=nil{
		return nil,err
	}
	return conn,nil
}