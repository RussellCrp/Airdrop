# Airdrop Platform

一个基于 go-zero + GORM + MySQL 的任务空投服务，并配套 Foundry Solidity 合约（Merkle 防女巫）。

## 目录结构

- `backend/service/airdrop/`：Go 服务，包含 API、逻辑、事件监听与单元测试。
- `contracts/`：Foundry 工程，包含 `AirdropToken` 与 `AirdropDistributor` 合约、部署脚本与测试。
- `backend/service/airdrop/model/airdrop.sql`：完整的数据库 DDL，可直接导入。

## 后端快速开始

1. **准备 MySQL**
   ```bash
   mysql -uroot -p123456 < backend/service/airdrop/model/airdrop.sql
   ```
2. **配置**：复制 `etc/airdrop-api.yaml` 按需修改 `Mysql.DSN`、`Auth.AccessSecret`、`Admin.Wallets` 等。
3. **运行服务**：
   ```bash
   cd backend/service/airdrop
   go run airdrop.go -f etc/airdrop-api.yaml
   ```
4. **运行单元测试**：
   ```bash
   go test ./...
   ```

### 登录与 JWT

- 客户端需构造消息 `airdrop-login:<wallet>:<timestamp>` 并使用钱包私钥 (personal_sign) 进行签名。
- 登录成功后返回 `accessToken`，后续接口在 `Authorization: Bearer <token>` 中携带。
- 系统记录连续登录天数并按 100~500 积分阶梯奖励。

### 任务与空投流程

- 管理员使用 `/api/v1/admin/tasks/award` 按任务代码（`PROMO`、`INVEST`、`REFERRAL`）发放积分。
- `/api/v1/admin/airdrop/start` 会冻结所有用户当前积分并写入 `round_points`，用于当轮空投。
- 用户使用 `/api/v1/airdrop/claim` 记录链上领取意图。链上 `Claimed` 事件由 `claim_watcher` 实时监听，成功后扣减冻结积分并更新 `claims` 状态。

## 智能合约 (Foundry)

1. 安装依赖后运行测试：
   ```bash
   cd contracts
   forge test
   ```
2. 部署脚本示例：`script/Deploy.s.sol` 读取 `PRIVATE_KEY`，部署 `AirdropToken` 与 `AirdropDistributor` 并为分发器预铸代币。
3. `AirdropDistributor` 通过存储每轮的 Merkle Root 与截止时间，`claim(roundId, amount, proof)` 会校验证明、防止重复领取并发出 `Claimed` 事件。

## 事件监听

`etc/airdrop-api.yaml` 中 `Eth` 段落开启后，服务启动会连接指定 RPC，订阅 `Claimed` 日志，自动同步 `claims` 表与 `round_points.claimed_points`。

## 常用命令

```bash
# 运行 go-zero 服务
cd backend/service/airdrop && go run airdrop.go -f etc/airdrop-api.yaml

# 运行 Go 单元测试
go test ./...

# 运行 Foundry 测试
cd contracts && forge test
```
