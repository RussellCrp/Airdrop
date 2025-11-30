# Airdrop Platform

一个基于 go-zero + GORM + MySQL 的任务空投服务，并配套 Foundry Solidity 合约（Merkle 防女巫）。

## 📋 项目简介

Airdrop Platform 是一个完整的链上空投解决方案，包含：

- **后端服务**：基于 go-zero 框架的 RESTful API 服务，提供用户认证、任务管理、积分系统、空投管理等功能
- **智能合约**：基于 Foundry 开发的 Solidity 合约，使用 Merkle Tree 实现防女巫攻击的空投分发机制
- **事件监听**：实时监听链上事件，自动同步链上数据到数据库

## 🏗️ 技术栈

### 后端
- **框架**：go-zero v1.9.3
- **数据库**：MySQL (GORM v1.31.1)
- **认证**：JWT (golang-jwt/jwt/v5)
- **区块链交互**：go-ethereum v1.16.7
- **Go 版本**：1.24.0

### 智能合约
- **框架**：Foundry
- **Solidity 版本**：^0.8.24
- **依赖**：OpenZeppelin Contracts

## 📁 项目结构

```
Airdrop/
├── backend/
│   └── service/
│       └── airdrop/          # Go 后端服务
│           ├── api/          # API 定义文件
│           ├── etc/          # 配置文件
│           ├── internal/     # 内部逻辑
│           ├── sql/          # 数据库 DDL
│           ├── go.mod        # Go 依赖
│           └── airdrop.go    # 服务入口
├── contracts/                # Foundry 智能合约工程
│   ├── src/                 # 合约源码
│   │   ├── AirdropToken.sol
│   │   ├── AirdropDistributor.sol
│   │   └── MerkleProof.sol
│   ├── script/              # 部署脚本
│   │   └── Deploy.s.sol
│   ├── test/                # 合约测试
│   └── foundry.toml         # Foundry 配置
└── README.md
```

## 🚀 快速开始

### 环境要求

- Go 1.24.0+
- MySQL 5.7+
- Foundry (用于智能合约开发)

### 后端服务

#### 1. 安装依赖

```bash
cd backend/service/airdrop
go mod download
```

#### 2. 初始化数据库

```bash
mysql -uroot -p123456 < sql/airdrop.sql
```

或者手动创建数据库并导入：

```bash
mysql -uroot -p123456
CREATE DATABASE airdrop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE airdrop;
SOURCE sql/airdrop.sql;
```

#### 3. 配置服务

编辑 `etc/airdrop-api.yaml` 配置文件：

```yaml
Name: airdrop-api
Host: 0.0.0.0
Port: 8888
Mode: dev

Mysql:
  DSN: root:123456@tcp(127.0.0.1:3306)/airdrop?charset=utf8mb4&parseTime=true&loc=Local
  MaxIdle: 5
  MaxOpen: 10

Auth:
  AccessSecret: "replace-with-strong-secret"  # 请修改为强密钥
  AccessExpire: 72000  # JWT 过期时间（秒）

Admin:
  Wallets:
    - "0xa27b6d5f1c0Ce106428B128307b652Ba6d1ba6c5"  # 管理员钱包地址

Eth:
  Enabled: false  # 是否启用事件监听
  RPC: ""  # 以太坊 RPC 地址
  DistributorAddress: ""  # AirdropDistributor 合约地址
  StartBlock: 0  # 开始监听的区块号
  PollInterval: 15s  # 轮询间隔
```

**重要配置项说明**：
- `Mysql.DSN`：MySQL 数据库连接字符串
- `Auth.AccessSecret`：JWT 签名密钥，生产环境请使用强密钥
- `Admin.Wallets`：管理员钱包地址列表，用于管理员操作
- `Eth`：区块链事件监听配置（可选）

#### 4. 运行服务

```bash
cd backend/service/airdrop
go run airdrop.go -f etc/airdrop-api.yaml
```

服务启动后，默认监听 `http://0.0.0.0:8888`

#### 5. 运行测试

```bash
go test ./...
```

### 智能合约

#### 1. 安装 Foundry

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

#### 2. 安装依赖

```bash
cd contracts
forge install
```

#### 3. 编译合约

```bash
forge build
```

#### 4. 运行测试

```bash
forge test
```

#### 5. 部署合约

编辑 `script/Deploy.s.sol` 并设置环境变量：

```bash
export PRIVATE_KEY=your_private_key
export RPC_URL=your_rpc_url

forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --verify
```

## 🔐 认证机制

### 钱包登录

客户端需要构造登录消息并使用钱包私钥签名：

1. **构造消息**：`airdrop-login:<wallet>:<timestamp>`
   - `wallet`：钱包地址
   - `timestamp`：Unix 时间戳（秒）

2. **签名**：使用钱包的 `personal_sign` 方法对消息进行签名

3. **请求登录**：
   ```http
   POST /api/airdrop/v1/auth/login
   Content-Type: application/json
   
   {
     "wallet": "0x...",
     "signature": "0x...",
     "timestamp": 1234567890
   }
   ```

4. **响应**：返回 `accessToken` 和过期时间
   ```json
   {
     "code": 0,
     "msg": "success",
     "data": {
       "accessToken": "eyJ...",
       "expiresAt": 1234567890,
       "loginDays": 1,
       "points": 0
     }
   }
   ```

5. **使用 Token**：后续请求在 Header 中携带
   ```http
   Authorization: Bearer <accessToken>
   ```

### 连续登录奖励

系统会自动记录用户的连续登录天数，奖励积分计算公式为：**奖励 = 连续登录天数 × 100**

- 连续登录 1 天：100 积分
- 连续登录 2 天：200 积分
- 连续登录 3 天：300 积分
- 连续登录 4 天：400 积分
- 连续登录 5 天及以上：500 积分（最大值）

**注意**：
- 如果昨天登录过，连续登录天数 +1（但最大值为 5）
- 如果昨天没登录，连续登录天数重置为 1
- 每天首次登录时发放奖励

## 📡 API 接口

### 用户接口

#### 1. 登录
```http
POST /api/airdrop/v1/auth/login
```

#### 2. 查询积分
```http
GET /api/airdrop/v1/me/points
Authorization: Bearer <token>
```

响应示例：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "wallet": "0x...",
    "available": 1000,
    "frozen": 500,
    "latestRound": 1,
    "loginStreak": 5
  }
}
```

#### 3. 提交任务
```http
POST /api/airdrop/v1/airdrop/task/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "wallet": "0x...",
  "taskCode": "PROMO",
  "proveParams": "..."
}
```

任务代码：
- `PROMO`：推广任务
- `INVEST`：投资任务
- `REFERRAL`：推荐任务

#### 4. 获取领取证明
```http
GET /api/airdrop/v1/airdrop/me/proof?roundId=1
Authorization: Bearer <token>
```

返回 Merkle 证明，用于链上领取：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "roundId": 1,
    "wallet": "0x...",
    "amount": 1000,
    "proof": ["0x...", "0x..."]
  }
}
```

### 管理员接口

#### 1. 启动空投轮次
```http
POST /api/airdrop/v1/admin/airdrop/start
Authorization: Bearer <token>
Content-Type: application/json

{
  "roundName": "Round 1",
  "tokenAddress": "0x...",
  "claimDeadline": 1234567890
}
```

此接口会：
- 冻结所有用户当前可用积分
- 创建新的空投轮次
- 将冻结积分写入 `round_points` 表

#### 2. 查询轮次信息
```http
GET /api/airdrop/v1/admin/airdrop/round?roundId=1
Authorization: Bearer <token>
```

## 🔄 空投流程

### 1. 任务积分发放

管理员通过 `/api/airdrop/v1/admin/tasks/award` 接口（如果存在）或直接操作数据库，按任务代码发放积分。

### 2. 启动空投轮次

管理员调用 `/api/airdrop/v1/admin/airdrop/start`：
- 系统冻结所有用户当前可用积分
- 生成 Merkle Tree
- 创建新的空投轮次记录
- 在链上调用 `AirdropDistributor.startRound()`

### 3. 用户领取

1. 用户调用 `/api/airdrop/v1/airdrop/me/proof` 获取 Merkle 证明
2. 用户在链上调用 `AirdropDistributor.claim(roundId, amount, proof)`
3. 系统的事件监听器（`claim_watcher`）监听到 `Claimed` 事件
4. 自动更新数据库：
   - 扣减用户的冻结积分
   - 更新 `claims` 表状态
   - 更新 `round_points.claimed_points`

## 🔗 智能合约

### AirdropToken

标准的 ERC20 代币合约，用于空投分发。

**主要功能**：
- `mint(address to, uint256 amount)`：铸造代币（仅 Owner）

### AirdropDistributor

基于 Merkle Tree 的空投分发合约，防止女巫攻击和重复领取。

**主要功能**：
- `startRound(uint256 roundId, bytes32 merkleRoot, uint64 claimDeadline)`：启动新的空投轮次
- `closeRound(uint256 roundId)`：关闭空投轮次
- `claim(uint256 roundId, uint256 amount, bytes32[] proof)`：领取空投
- `claimed(uint256 roundId, address account)`：查询是否已领取

**安全特性**：
- Merkle Proof 验证
- 防重复领取（每轮每个地址只能领取一次）
- 截止时间检查
- 轮次状态检查

## 📊 事件监听

如果启用了 `Eth.Enabled`，服务启动时会：

1. 连接到指定的以太坊 RPC
2. 订阅 `AirdropDistributor.Claimed` 事件
3. 实时监听链上领取操作
4. 自动同步数据到数据库

**配置示例**：
```yaml
Eth:
  Enabled: true
  RPC: "https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY"
  DistributorAddress: "0x..."
  StartBlock: 18000000
  PollInterval: 15s
```

## 🧪 测试

### 后端测试

```bash
cd backend/service/airdrop
go test ./...
```

### 合约测试

```bash
cd contracts
forge test -vvv  # 详细输出
```

## 📝 常用命令

```bash
# 运行后端服务
cd backend/service/airdrop && go run airdrop.go -f etc/airdrop-api.yaml

# 运行后端测试
cd backend/service/airdrop && go test ./...

# 编译合约
cd contracts && forge build

# 运行合约测试
cd contracts && forge test

# 格式化合约代码
cd contracts && forge fmt

# 生成 Gas 快照
cd contracts && forge snapshot

# 启动本地测试节点
anvil
```

## 🔧 开发指南

### 添加新的 API 接口

1. 在 `api/airdrop.api` 中定义接口
2. 使用 `goctl` 生成代码：
   ```bash
   goctl api go -api api/airdrop.api -dir . -style gozero
   ```
3. 在 `internal/logic` 中实现业务逻辑

### 添加新的合约功能

1. 在 `contracts/src/` 中编写合约
2. 在 `contracts/test/` 中编写测试
3. 运行 `forge test` 验证

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📮 联系方式

如有问题，请提交 Issue 或联系项目维护者。
