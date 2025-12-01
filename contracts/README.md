## 智能合约（contracts）

本目录是 Airdrop 平台的 **智能合约工程**，使用 Foundry 开发和测试。

### 目录结构

```text
contracts/
├── src/               # 合约源码
│   ├── AirdropToken.sol
│   ├── AirdropDistributor.sol
│   └── MerkleProof.sol
├── script/            # 部署 & 运维脚本
│   ├── AirdropTokenDeploy.s.sol
│   ├── AirdropDeploy.s.sol
│   ├── AirdropTokenMint.s.sol
│   ├── ClaimAirdropToken.s.sol
│   └── GetAirdropTokenBanlance.s.sol
├── test/              # 合约测试
│   ├── AirdropDistributor.t.sol
│   └── AirdropDistributorMerkle.t.sol
├── broadcast/         # Forge 脚本广播记录
├── foundry.toml       # Foundry 配置
└── foundry.lock
```

### 环境准备

- 已安装 Foundry（`forge`, `cast`, `anvil` 等）  
  如未安装，可参照官方文档：`https://book.getfoundry.sh/`

### 安装依赖

```bash
cd contracts
forge install
```

### 编译合约

```bash
cd contracts
forge build
```

### 运行测试

```bash
cd contracts
forge test -vvv
```

### 常用脚本

- **部署代币合约**

  ```bash
  cd contracts
  forge script script/AirdropTokenDeploy.s.sol:AirdropTokenDeploy \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
  ```

- **部署 AirdropDistributor**

  ```bash
  cd contracts
  forge script script/AirdropDeploy.s.sol:AirdropDeploy \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
  ```

- **给 Distributor 铸造/转入代币**

  ```bash
  cd contracts
  forge script script/AirdropTokenMint.s.sol:AirdropTokenMint \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
  ```

- **链上领取空投**

  ```bash
  cd contracts
  forge script script/ClaimAirdropToken.s.sol:ClaimAirdropToken \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
  ```

- **查询空投代币余额**

  ```bash
  cd contracts
  forge script script/GetAirdropTokenBanlance.s.sol:GetAirdropTokenBanlance \
    --rpc-url $RPC_URL
  ```

### 提示

- 后端相关说明请查看 `backend/README.md`。  
- 根目录的 `README.md` 只作为整体项目的介绍与导航。  
