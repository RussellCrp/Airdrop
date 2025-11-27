package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ClaimsModel = (*customClaimsModel)(nil)

type (
	// ClaimsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customClaimsModel.
	ClaimsModel interface {
		claimsModel
		withSession(session sqlx.Session) ClaimsModel
	}

	customClaimsModel struct {
		*defaultClaimsModel
	}
)

// NewClaimsModel returns a model for the database table.
func NewClaimsModel(conn sqlx.SqlConn) ClaimsModel {
	return &customClaimsModel{
		defaultClaimsModel: newClaimsModel(conn),
	}
}

func (m *customClaimsModel) withSession(session sqlx.Session) ClaimsModel {
	return NewClaimsModel(sqlx.NewSqlConnFromSession(session))
}
