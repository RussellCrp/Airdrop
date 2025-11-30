#!/bin/bash
# 从 Foundry 编译结果中提取 ABI 和 bin 文件的脚本

CONTRACT_NAME=${1:-"AirdropDistributor"}
OUTPUT_DIR="out/${CONTRACT_NAME}.sol"
JSON_FILE="${OUTPUT_DIR}/${CONTRACT_NAME}.json"
ABI_FILE="${CONTRACT_NAME}.abi.json"
BIN_FILE="${CONTRACT_NAME}.bin"

if [ ! -f "$JSON_FILE" ]; then
    echo "错误: 找不到文件 $JSON_FILE"
    echo "请先运行: forge build"
    exit 1
fi

# 提取 ABI
python3 << EOF
import json
import sys

try:
    with open("$JSON_FILE", "r") as f:
        data = json.load(f)
    
    # 提取 ABI
    abi = data.get("abi", [])
    if abi:
        with open("$ABI_FILE", "w") as out:
            json.dump(abi, out, indent=2)
        print(f"✓ 已生成: $ABI_FILE")
        print(f"  ABI 包含 {len(abi)} 个项")
    else:
        print("错误: 找不到 ABI 字段")
        sys.exit(1)
    
    # 提取字节码
    bytecode = data.get("bytecode", {}).get("object", "")
    if bytecode:
        with open("$BIN_FILE", "w") as out:
            out.write(bytecode)
        print(f"✓ 已生成: $BIN_FILE")
        print(f"  字节码长度: {len(bytecode)} 字符")
    else:
        print("错误: 找不到 bytecode 字段")
        sys.exit(1)
        
except Exception as e:
    print(f"错误: {e}")
    sys.exit(1)
EOF

echo ""
echo "文件位置:"
echo "  ABI: $(pwd)/$ABI_FILE"
echo "  BIN: $(pwd)/$BIN_FILE"

