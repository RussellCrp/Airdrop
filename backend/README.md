## 后端服务说明（Airdrop Backend）

本目录包含 Airdrop 项目的后端服务，基于 **Go 语言 + go-zero** REST 框架，使用 **MySQL / GORM** 作为数据存储，并集成以太坊链上空投监听逻辑。

后端主要功能：
- **钱包登录**：通过钱包地址登录，生成 JWT 访问令牌。
- **积分系统**：记录用户积分、登录天数与任务完成情况。
- **任务系统**：用户完成任务（如登录任务）获取积分。
- **空投轮次管理**：管理员创建/查询空投轮次、配置 Merkle Root 与代币信息。
- **空投证明生成**：为用户生成 Merkle Proof，供链上合约验证。
- **链上监听（可选）**：根据配置监听以太坊链上空投合约的领取事件。

后端服务代码路径：`backend/service/airdrop`

---

## 技术栈

- **语言**：Go `1.24.x`
- **框架**：`github.com/zeromicro/go-zero` REST
- **数据库**：MySQL（通过 `gorm.io/gorm` 访问）
- **ORM & 驱动**：
  - `gorm.io/gorm`
  - `gorm.io/driver/mysql`
- **区块链相关**：
  - `github.com/ethereum/go-ethereum`
- **认证与安全**：
  - JWT：`github.com/golang-jwt/jwt/v5`

---

## 目录结构

```text
backend/
  service/
    airdrop/
      airdrop.go              # 程序入口，启动 HTTP 服务
      go.mod / go.sum         # Go 依赖管理
      api/
        airdrop.api           # go-zero API 描述文件
        airdrop_request.http  # 接口调试示例
      etc/
        airdrop-api.yaml      # 配置文件（端口、DB、Auth、链上配置等）
      internal/
        authctx/              # 认证上下文相关
        config/               # 配置结构定义
        contract/             # 合约交互封装
        entity/               # GORM 实体（数据库表结构）
        handler/              # HTTP handler 层（路由对应）
        listener/             # 链上事件监听（空投领取等）
        logic/                # 业务逻辑层
        middleware/           # JWT、Admin 等中间件
        security/             # JWT 管理等安全相关
        svc/                  # ServiceContext，资源与依赖管理
        tasks/                # 任务处理逻辑（积分任务等）
        types/                # API 请求/响应结构体
        util/                 # 工具函数
        testutil/             # 测试辅助
      sql/
        ...                   # 数据库初始化或迁移脚本（如提供）
```

---

## 配置说明（`etc/airdrop-api.yaml`）

示例配置（已存在文件简化版）：

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
  AccessSecret: "replace-with-strong-secret"
  AccessExpire: 72000

Admin:
  Wallets:
    - "0xa27b6d5f1c0Ce106428B128307b652Ba6d1ba6c5"

Eth:
  Enabled: false
  RPC: ""
  DistributorAddress: ""
  StartBlock: 0
  PollInterval: 15s
  OwnerPrivateKey: ""
  ChainID: 11155111
  GasLimit: 0
  GasPrice: 0
```

**关键字段说明：**
- **Mysql.DSN**：MySQL 连接串，需要根据本地/线上环境进行修改。
- **Auth.AccessSecret**：JWT 签名密钥，务必在生产环境替换为强随机值。
- **Auth.AccessExpire**：JWT 过期时间（秒）。
- **Admin.Wallets**：拥有管理员权限的钱包地址列表。
- **Eth.Enabled**：是否启用链上监听功能（启用前需同时配置 RPC、合约地址等）。

---

## 安装与运行

### 1. 准备环境

- 已安装 **Go 1.24.x**（或兼容版本）
- 已安装并启动 **MySQL**，并创建数据库：

```sql
CREATE DATABASE airdrop CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

- 根据实际情况修改 `etc/airdrop-api.yaml` 中的：
  - `Mysql.DSN`
  - `Auth.AccessSecret`
  - `Admin.Wallets`
  - （可选）`Eth` 相关配置

### 2. 安装依赖

在 `backend/service/airdrop` 目录下执行：

```bash
cd backend/service/airdrop
go mod tidy
```

### 3. 启动服务

在 `backend/service/airdrop` 目录下执行：

```bash
go run airdrop.go -f etc/airdrop-api.yaml
```

启动成功后，默认监听：`http://0.0.0.0:8888`

日志中会输出：

```text
Starting server at 0.0.0.0:8888...
```

---

## 主要接口概览

API 描述文件：`api/airdrop.api`，核心接口包括：

- **用户认证**
  - `POST /api/airdrop/v1/auth/login`
    - 请求体：`LoginRequest{ wallet, signature, timestamp }`
    - 响应体：`LoginResponse{ data: { accessToken, expiresAt, loginDays, points } }`

- **用户积分**
  - `GET /api/airdrop/v1/me/points`
    - 需要在请求头携带 `Authorization: Bearer <token>`
    - 返回当前用户积分、登录天数等信息。

- **任务提交**
  - `POST /api/airdrop/v1/airdrop/task/submit`
    - 请求体：`SubmitTaskRequest{ wallet, taskCode, proveParams }`
    - 用于提交任务完成证明，增加积分。

- **空投轮次（管理员）**
  - `POST /api/airdrop/v1/admin/airdrop/start`
    - 创建空投轮次，设置 `roundName`、`merkleRoot`、`tokenAddress`、`claimDeadline` 等。
  - `GET /api/airdrop/v1/admin/airdrop/round`
    - 查询当前/指定空投轮次信息。

- **空投证明**
  - `GET /api/airdrop/v1/airdrop/me/proof`
    - 请求参数：`roundId`
    - 返回当前用户地址的 Merkle Proof、可领取额度等。

建议使用 `api/airdrop_request.http` 文件中的示例请求配合 VS Code / IDEA 插件进行调试。

---

## 运行模式与部署建议

- **开发模式（本地）**
  - `Mode: dev`，日志更加详细。
  - 建议使用 SQLite 或本地 MySQL，调整 `Mysql.DSN`。
  - `Eth.Enabled` 可保持 `false`，仅验证业务逻辑与 API。

- **生产/测试环境**
  - 使用受控的 MySQL 实例。
  - 将 `Auth.AccessSecret` 配置为安全随机值，并通过环境变量或配置管理系统注入。
  - 若需要链上监听，将 `Eth.Enabled` 设为 `true` 并正确填充 RPC、合约地址、OwnerPrivateKey 等字段。

---

## 开发说明

- **路由与接口定义**
  - 修改 `api/airdrop.api` 后，可使用 `goctl` 重新生成 handler、types 等代码（注意保留已写好的业务逻辑）。
  - 当前生效的路由注册位于 `internal/handler/routes.go`。

- **业务逻辑**
  - 登录逻辑位于 `internal/logic/loginlogic.go`，会：
    - 归一化钱包地址；
    - 校验登录时间戳；
    - 创建/查找用户记录；
    - 根据是否为管理员钱包设置不同角色；
    - 生成 JWT；
    - 触发登录任务以增加用户积分。

- **ServiceContext 与资源管理**
  - `internal/svc/servicecontext.go` 负责：
    - 初始化数据库连接（GORM + MySQL）；
    - 初始化 JWTManager；
    - 构建中间件（JWT、Admin）；
    - 管理管理员钱包集合；
    - 在 `Eth.Enabled` 时启动链上事件监听。

---

## TODO / 后续优化方向

- 增加更多任务类型与积分规则的配置化支持。
- 为接口补充更完整的 OpenAPI / Swagger 文档。
- 增加单元测试与集成测试覆盖率（`internal/testutil` 可复用）。
- 支持多链/多合约空投配置。


