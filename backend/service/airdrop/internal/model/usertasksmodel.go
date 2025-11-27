package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserTasksModel = (*customUserTasksModel)(nil)

type (
	// UserTasksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTasksModel.
	UserTasksModel interface {
		userTasksModel
		withSession(session sqlx.Session) UserTasksModel
	}

	customUserTasksModel struct {
		*defaultUserTasksModel
	}
)

// NewUserTasksModel returns a model for the database table.
func NewUserTasksModel(conn sqlx.SqlConn) UserTasksModel {
	return &customUserTasksModel{
		defaultUserTasksModel: newUserTasksModel(conn),
	}
}

func (m *customUserTasksModel) withSession(session sqlx.Session) UserTasksModel {
	return NewUserTasksModel(sqlx.NewSqlConnFromSession(session))
}
