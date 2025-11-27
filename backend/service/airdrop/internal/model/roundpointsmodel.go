package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RoundPointsModel = (*customRoundPointsModel)(nil)

type (
	// RoundPointsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoundPointsModel.
	RoundPointsModel interface {
		roundPointsModel
		withSession(session sqlx.Session) RoundPointsModel
	}

	customRoundPointsModel struct {
		*defaultRoundPointsModel
	}
)

// NewRoundPointsModel returns a model for the database table.
func NewRoundPointsModel(conn sqlx.SqlConn) RoundPointsModel {
	return &customRoundPointsModel{
		defaultRoundPointsModel: newRoundPointsModel(conn),
	}
}

func (m *customRoundPointsModel) withSession(session sqlx.Session) RoundPointsModel {
	return NewRoundPointsModel(sqlx.NewSqlConnFromSession(session))
}
