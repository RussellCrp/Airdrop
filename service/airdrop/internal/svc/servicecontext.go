// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"airdrop/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DSN), &gorm.Config{})
	if err != nil {
		logx.Errorf("failed to connect mysql: %v", err)
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
