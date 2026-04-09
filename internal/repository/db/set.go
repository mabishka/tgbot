package db

import (
	"context"
	"tgbot/internal/logger"
	"tgbot/internal/model"
	"unsafe"

	"go.uber.org/zap"
)

const (
	qAddUser   = "insert into users(id, chat_id, name, created) values ($1, $2, $3, current_timestamp) ON CONFLICT DO NOTHING"
	qAddTask   = "insert int tasks(user_id, chat_id, task_id, input_file_id, output_file_id, result, result_short) values($1, $2, $3, $4, $5, $6, $7)"
	qSetInfo   = "update task set task_id = $1, file_id = $2 where id = $3"
	qSetResult = "update task set result = $1 where id = $2 "

	qAddWord = "insert int words(id, word) values($1, $2)"
)

// SaveUser - регистрация пользователя – запоминаем его идентификатор.
func AddUser(ctx context.Context, conn model.Connection, id, chatId int64, name string) error {
	if err := conn.Execute(ctx, qAddUser, id, chatId, name); err != nil {
		logger.Log().Error("error AddUser - Query", zap.Error(err))
		return err
	}
	return nil
}

func AddTask(ctx context.Context, conn model.Connection, x *model.SpeachTaskResponse, short string) error {
	return conn.Execute(ctx, qAddTask, x.User, x.ChatID, x.TaskID, x.InFileID, x.OutFileID, unsafe.String(unsafe.SliceData(x.Output), len(x.Output)), short)
}

func SetTaskInfo(ctx context.Context, conn model.Connection, id, fileID, taskID string) error {
	return conn.Execute(ctx, qSetInfo, taskID, fileID, id)
}

func SetTaskResult(ctx context.Context, conn model.Connection, id string, x []byte) error {
	return conn.Execute(ctx, qSetResult, x, id)
}

func AddWord(ctx context.Context, conn model.Connection, words []model.Words) error {
	tx, err := conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	for _, v := range words {
		if err := conn.Execute(ctx, qAddWord, v.ID, v.Word); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
