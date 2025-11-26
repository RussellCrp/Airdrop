// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	// Mysql 数据库配置
	Mysql struct {
		DSN string `json:"dsn,optional"`
	} `json:"Mysql"`
}
