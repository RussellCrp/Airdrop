## 项目说明：任务制 ERC20 空投系统（go-zero + gorm + Foundry）

### 后端技术栈
- **语言**：Go
- **框架**：`go-zero`（REST API） + `gorm`（MySQL ORM）
- **数据库**：MySQL，连接：`127.0.0.1:3306/airdrop`，账号/密码：`root/123456`
- **目录**：
  - `api/airdrop.api`：API 描述，用 `goctl api go` 生成服务骨架
  - `service/airdrop`：后端服务
    - `internal/model`：gorm 实体
    - `internal/logic`：业务逻辑（任务、积分、快照、Merkle）
    - `internal/merkle`：Merkle Tree 构建与 leaf 编码
    - `internal/svc`：`ServiceContext`，内含 `*gorm.DB`
    - `etc/airdrop-api.yaml`：配置（含 MySQL DSN）
  - `sql/airdrop_schema.sql`：建库建表 SQL
  - `sql/airdrop_seed.sql`：基础任务与示例数据

### Solidity / Foundry
- **目录**：
  - `contracts/AirdropMerkle.sol`：Merkle 白名单空投合约（防女巫）
  - `script/DeployAirdrop.s.sol`：部署脚本（使用环境变量 `AIRDROP_TOKEN`、`AIRDROP_MERKLE_ROOT`）
  - `test/AirdropMerkle.t.sol`：Foundry 单元测试
  - `foundry.toml`：Foundry 配置，包含 OpenZeppelin remapping
- **合约特性**：
  - 使用已有 ERC20 代币地址
  - `merkleRoot` 控制一轮空投白名单
  - `claim(index, account, amount, proof)`：使用 `keccak256(abi.encodePacked(index, account, amount))` + MerkleProof 验证
  - `claimed[index]` 防重复领取
  - 管理员可更新 `merkleRoot`、紧急提币

### 启动步骤（后端）
1. **准备 MySQL**
   - 确保本地 MySQL 运行在 `127.0.0.1:3306`，创建用户 `root/123456`
2. **初始化数据库**
   - 登录 MySQL：
     - `mysql -uroot -p123456`
   - 执行建表脚本（会自动创建 `airdrop` 库）：
     - `SOURCE sql/airdrop_schema.sql;`
   - 切库：
     - `USE airdrop;`
   - 执行基础数据脚本：
     - `SOURCE sql/airdrop_seed.sql;`
3. **生成 / 更新 go-zero 代码（如需）**
   - `cd service/airdrop`
   - `goctl api go -api ../../api/airdrop.api -dir .`
4. **启动服务**
   - 在 `service/airdrop` 目录：
   - `go run airdrop.go -f etc/airdrop-api.yaml`
   - 默认监听：`http://0.0.0.0:8888`

### 主要接口（HTTP 简要）
- `POST /api/v1/users`：创建/注册用户（绑定钱包地址）
- `GET /api/v1/users/{id}`：查询用户信息 + 当前积分 & 档位
- `GET /api/v1/users/{id}/tasks`：查询用户任务完成情况
- `GET /api/v1/users/{id}/score`：查询积分 & 档位
- 任务上报：
  - `POST /api/v1/tasks/promote`
  - `POST /api/v1/tasks/quant_invest`
  - `POST /api/v1/tasks/trade_volume`
  - `POST /api/v1/tasks/referral`
  - `POST /api/v1/tasks/login_streak`
- 空投快照 & Merkle：
  - `POST /api/v1/airdrop/snapshot`：生成快照（管理员）
  - `GET /api/v1/airdrop/snapshot/{id}`：查询快照
  - `GET /api/v1/airdrop/snapshot/{id}/proof?address=...`：查询某地址的 Merkle proof

### Go 单元测试（后端）
- 测试文件在 `service/airdrop/internal/logic/*.go` 中：
  - `logic_test.go`：公共 TestMain，使用真实 MySQL 连接 `airdrop` 库
  - `user_logic_test.go`：用户创建 & 查询
  - `task_logic_test.go`：任务上报 & 积分更新
  - `snapshot_logic_test.go`：快照生成 & Merkle proof 查询
- **运行前要求**：
  - 按“初始化数据库”步骤，确保 `airdrop` 库已经按脚本建好并写入 `airdrop_seed.sql` 数据
- 运行测试：
  - `cd service/airdrop`
  - `go test ./internal/logic -run .`

### Foundry 测试与部署
1. **安装依赖**
   - `forge install OpenZeppelin/openzeppelin-contracts --no-commit`
2. **运行合约测试**
   - `forge test`
3. **部署示例**
   - 设置环境变量（示例）：
     - `export AIRDROP_TOKEN=0xYourErc20Address`
     - `export AIRDROP_MERKLE_ROOT=0xYourMerkleRoot`
   - 运行：
     - `forge script script/DeployAirdrop.s.sol:DeployAirdrop --rpc-url <RPC_URL> --private-key <PK> --broadcast`



