package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TasksModel = (*customTasksModel)(nil)

type (
	// TasksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTasksModel.
	TasksModel interface {
		tasksModel
		withSession(session sqlx.Session) TasksModel
	}

	customTasksModel struct {
		*defaultTasksModel
	}
)

// NewTasksModel returns a model for the database table.
func NewTasksModel(conn sqlx.SqlConn) TasksModel {
	return &customTasksModel{
		defaultTasksModel: newTasksModel(conn),
	}
}

func (m *customTasksModel) withSession(session sqlx.Session) TasksModel {
	return NewTasksModel(sqlx.NewSqlConnFromSession(session))
}
