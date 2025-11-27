package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PointsLedgerModel = (*customPointsLedgerModel)(nil)

type (
	// PointsLedgerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPointsLedgerModel.
	PointsLedgerModel interface {
		pointsLedgerModel
		withSession(session sqlx.Session) PointsLedgerModel
	}

	customPointsLedgerModel struct {
		*defaultPointsLedgerModel
	}
)

// NewPointsLedgerModel returns a model for the database table.
func NewPointsLedgerModel(conn sqlx.SqlConn) PointsLedgerModel {
	return &customPointsLedgerModel{
		defaultPointsLedgerModel: newPointsLedgerModel(conn),
	}
}

func (m *customPointsLedgerModel) withSession(session sqlx.Session) PointsLedgerModel {
	return NewPointsLedgerModel(sqlx.NewSqlConnFromSession(session))
}
