package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AirdropRoundsModel = (*customAirdropRoundsModel)(nil)

type (
	// AirdropRoundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAirdropRoundsModel.
	AirdropRoundsModel interface {
		airdropRoundsModel
		withSession(session sqlx.Session) AirdropRoundsModel
	}

	customAirdropRoundsModel struct {
		*defaultAirdropRoundsModel
	}
)

// NewAirdropRoundsModel returns a model for the database table.
func NewAirdropRoundsModel(conn sqlx.SqlConn) AirdropRoundsModel {
	return &customAirdropRoundsModel{
		defaultAirdropRoundsModel: newAirdropRoundsModel(conn),
	}
}

func (m *customAirdropRoundsModel) withSession(session sqlx.Session) AirdropRoundsModel {
	return NewAirdropRoundsModel(sqlx.NewSqlConnFromSession(session))
}
