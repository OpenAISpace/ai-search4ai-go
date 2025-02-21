#!/bin/bash

# 定义服务名称和路径
SERVICE_NAME="search4ai"
BINARY_PATH="/usr/local/bin/search4ai"
SERVICE_PATH="/etc/systemd/system/${SERVICE_NAME}.service"
GO_BINARY="search4ai"
ENV_SOURCE=".env"
ENV_DEST="/usr/local/bin/.env"

# 检查是否以 root 权限运行
if [ "$EUID" -ne 0 ]; then 
    echo "请使用 root 权限运行此脚本"
    exit 1
fi

# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo "未检测到 Go 环境，请先安装 Go"
    exit 1
fi

# 检查 .env 文件是否存在
if [ ! -f "${ENV_SOURCE}" ]; then
    echo "错误：${ENV_SOURCE} 文件不存在！"
    exit 1
fi

# 编译项目
echo "开始编译项目..."
go build -o ${GO_BINARY} main.go
if [ $? -ne 0 ]; then
    echo "编译失败！"
    exit 1
fi

# 停止已存在的服务
if systemctl is-active --quiet ${SERVICE_NAME}; then
    echo "停止已存在的服务..."
    systemctl stop ${SERVICE_NAME}
fi

# 复制二进制文件
echo "复制二进制文件到 ${BINARY_PATH}..."
cp ${GO_BINARY} ${BINARY_PATH}
chmod +x ${BINARY_PATH}

# 复制 .env 文件
echo "复制 .env 文件到 ${ENV_DEST}..."
cp ${ENV_SOURCE} ${ENV_DEST}
chmod 600 ${ENV_DEST}  # 设置适当的权限，只允许 root 读写

# 创建服务文件
echo "创建系统服务文件..."
cat > ${SERVICE_PATH} << EOF
[Unit]
Description=Search4AI Service
After=network.target

[Service]
Type=simple
ExecStart=${BINARY_PATH}
Restart=always
RestartSec=10
User=root
WorkingDirectory=/usr/local/bin
Environment="ENV_FILE=${ENV_DEST}"

[Install]
WantedBy=multi-user.target
EOF

# 重新加载 systemd
echo "重新加载 systemd..."
systemctl daemon-reload

# 启用并启动服务
echo "启用并启动服务..."
systemctl enable ${SERVICE_NAME}
systemctl start ${SERVICE_NAME}

# 检查服务状态
echo "检查服务状态..."
systemctl status ${SERVICE_NAME}

# 清理编译产物
echo "清理编译文件..."
rm ${GO_BINARY}

echo "安装完成！" 