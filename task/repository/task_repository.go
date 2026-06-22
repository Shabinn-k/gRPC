package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type TaskRespository struct {
	DB *pgx.Conn
}

func (r *TaskRespository)CreateTask(title ,description string)(int32,error){
	var id int32

	fmt.Println("Title:", title)
	fmt.Println("Description:", description)
	err:=r.DB.QueryRow(
		context.Background(),
		`INSERT INTO tasks(title,description)
		VALUES($1,$2)
		RETURNING id`,
		title,
		description,
	).Scan(&id)

	return id,err
}