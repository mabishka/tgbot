package db

import (
	"context"
	"tgbot/internal/model"
)

const q = `create table if not exists users (
    id bigint primary key,
    chat_id bigint,
    name text,
    created timestamp default now()
);

create table if not exists tasks (
    user_id bigint,
	chat_id bigint,
    task_id uuid,
    input_file_id uuid,
    output_file_id uuid,
    result text,
	result_short text,
    created timestamp default now()
);
`

// Create -  создание структуры таблиц
func Create(ctx context.Context, conn model.Connection) error {
	return conn.Execute(ctx, q)
}

/*
// List - список сохраненных встреч.
func List(ctx context.Context, conn Connection, user string, meet *model.Meet) []*model.Meet {
	resp := make([]*model.Meet, 0)

	rows, err := conn.Query(ctx, "select name from t_task t where user = $1 order by created", user)
	return resp
}

// Find - получение текста встречи.
func Find(ctx context.Context, conn Connection, user string, meet *model.Meet) []*model.Meet {
	resp := make([]*model.Meet, 0)

	rows, err := conn.Query(ctx, "select name from t_task t inner join t_task_result r on t.id = r.id and r.value ~ $1", data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return resp
}

// GetMeet - получение текста встречи.
func GetMeet(ctx context.Context, conn Connection, id string) ([]*model.Meet, error) {
	resp := make([]*model.Meet, 0)

	rows, err := conn.Query(ctx, "select name, r.value from t_task t inner join t_task_result r on t.id = r.id and r.id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return resp
}

func SaveTask(ctx context.Context, conn Connection, taskID, name, user string) error {
	return conn.Execute(ctx, "insert into t_task(id, user, created) values ($1, $2, current_timestamp)", taskID, user)
}

func SaveFile(ctx context.Context, conn Connection, taskID string, x []byte) error {
	return conn.Execute(ctx, "insert into t_task_result(id, value) values($1, $2)", taskID, x)
}
*/
