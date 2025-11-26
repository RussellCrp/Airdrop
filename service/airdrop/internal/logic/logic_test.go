package logic

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"airdrop/internal/config"
	"airdrop/internal/svc"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/conf"
)

var testSvcCtx *svc.ServiceContext

func TestMain(m *testing.M) {
	// 加载 etc/airdrop-api.yaml 配置
	var c config.Config
	conf.MustLoad("../../etc/airdrop-api.yaml", &c)
	testSvcCtx = svc.NewServiceContext(c)

	code := m.Run()
	os.Exit(code)
}

func resetTestDB() error {
	// 仅检查 MySQL 连接是否可用，具体建表和填充由 sql 脚本负责
	airdropDSN := "root:123456@tcp(127.0.0.1:3306)/airdrop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", airdropDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Ping()
}

func newCtx(t *testing.T) context.Context {
	t.Helper()
	return context.Background()
}
