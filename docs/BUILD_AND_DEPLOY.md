# 构建与部署文档

## 目录
- [项目概述](#项目概述)
- [架构说明](#架构说明)
- [环境要求](#环境要求)
- [本地开发](#本地开发)
- [构建指南](#构建指南)
- [部署指南](#部署指南)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

## 项目概述

本项目是一个 AI API 网关，支持多种 AI 服务提供商（OpenAI、Claude、Gemini 等），提供统一的 API 接口、用户管理、计费系统等功能。

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22+, Gin, GORM |
| 前端 | React 19, TypeScript, Rsbuild |
| 数据库 | SQLite, MySQL, PostgreSQL |
| 缓存 | Redis |
| 容器化 | Docker, Docker Compose |

## 架构说明

### 单二进制部署

本项目采用 **前端嵌入后端** 的架构，前端构建产物通过 Go 的 `embed` 指令嵌入到后端二进制文件中，最终生成一个可执行文件。

```
┌─────────────────────────────────────────────────────┐
│                   单个二进制文件                       │
│  ┌───────────────────────────────────────────────┐  │
│  │              Go 后端代码                        │  │
│  │  - API 路由                                    │  │
│  │  - 中间件 (鉴权、限流、CORS)                    │  │
│  │  - 业务逻辑 (计费、渠道路由)                     │  │
│  │  - 数据库交互                                   │  │
│  └───────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────┐  │
│  │              嵌入的前端资源                      │  │
│  │  - HTML/CSS/JavaScript                        │  │
│  │  - React 应用                                  │  │
│  │  - 静态资源 (图片、字体)                        │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### 构建流程

```
1. 构建前端
   web/default/ → web/default/dist/
   web/classic/ → web/classic/dist/

2. 嵌入到 Go 代码
   //go:embed web/default/dist
   var buildFS embed.FS

3. 编译为单个二进制文件
   go build -o server .
```

## 环境要求

### 后端开发
- Go 1.22 或更高版本
- GCC（CGO 依赖，SQLite 需要）
- Git

### 前端开发
- Node.js 18 或更高版本
- Bun（推荐）或 npm/yarn
- Git

### 生产部署
- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本

## 本地开发

### 1. 克隆项目

```bash
git clone <repository-url>
cd your-project
```

### 2. 启动后端

```bash
# 安装依赖
go mod tidy

# 运行开发服务器
go run .
```

后端将在 `http://localhost:3000` 启动。

### 3. 启动前端

```bash
# 进入前端目录
cd web/default

# 安装依赖
bun install

# 启动开发服务器
bun run dev
```

前端将在 `http://localhost:3001` 启动，自动代理 API 请求到后端。

### 4. 配置环境变量

复制 `.env.example` 为 `.env` 并根据需要修改：

```bash
cp .env.example .env
```

## 构建指南

### 构建单个二进制文件

```bash
# 一键构建（前端 + 后端）
make build
```

这将：
1. 构建默认主题前端 (`web/default/dist/`)
2. 构建经典主题前端 (`web/classic/dist/`)
3. 编译 Go 后端，将前端嵌入到二进制文件中
4. 输出可执行文件到 `bin/server`

### 手动构建

```bash
# 1. 构建前端
cd web/default && bun install && bun run build
cd web/classic && bun install && bun run build

# 2. 构建后端
CGO_ENABLED=0 go build -o bin/server .
```

### Docker 构建

```bash
# 构建 Docker 镜像
docker build -t your-project .

# 或使用 docker-compose
docker-compose build
```

## 部署指南

### 方式一：直接部署

#### 1. 构建项目

```bash
make build
```

#### 2. 配置环境变量

```bash
export GIN_MODE=release
export PORT=3000
export SQL_DSN="user:password@tcp(localhost:3306)/dbname"
export REDIS_CONN_STRING="localhost:6379"
export SESSION_SECRET="your-secret-key"
```

#### 3. 启动服务

```bash
./bin/server
```

### 方式二：Docker 部署

#### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件配置数据库等信息
```

#### 2. 启动服务

```bash
# 使用 docker-compose
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方式三：生产环境部署

#### 1. 服务器准备

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. 部署项目

```bash
# 克隆项目
git clone <repository-url>
cd your-project

# 配置环境变量
cp .env.example .env
nano .env

# 启动服务
docker-compose up -d
```

#### 3. 配置反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

#### 4. 配置 SSL（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 配置说明

### 数据库配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SQL_DSN` | 数据库连接字符串 | 空（使用 SQLite） |

#### SQLite（默认）
```bash
SQL_DSN=
```

#### MySQL
```bash
SQL_DSN=user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

#### PostgreSQL
```bash
SQL_DSN=host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai
```

### Redis 配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `REDIS_CONN_STRING` | Redis 连接地址 | 空（不使用 Redis） |
| `REDIS_PASSWORD` | Redis 密码 | 空 |
| `REDIS_DB` | Redis 数据库 | 0 |

### 安全配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SESSION_SECRET` | Session 密钥 | 必须配置 |
| `CRYPTO_SECRET` | 加密密钥 | 必须配置 |
| `CORS_ALLOW_ORIGINS` | CORS 允许的来源 | `*` |

### 其他配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `PORT` | 服务端口 | 3000 |
| `GIN_MODE` | Gin 运行模式 | debug |
| `TZ` | 时区 | Asia/Shanghai |

## 常见问题

### 1. 数据库连接失败

**问题**: `Failed to connect to database`

**解决方案**:
- 检查 `SQL_DSN` 配置是否正确
- 确保数据库服务正在运行
- 检查防火墙设置

### 2. Redis 连接失败

**问题**: `Warning: Failed to connect to Redis`

**解决方案**:
- Redis 是可选的，不配置会使用内存缓存
- 如果需要 Redis，确保 `REDIS_CONN_STRING` 正确
- 检查 Redis 服务是否运行

### 3. 前端构建失败

**问题**: `bun install` 或 `bun run build` 失败

**解决方案**:
- 确保 Node.js 版本 >= 18
- 清除缓存: `rm -rf node_modules && bun install`
- 检查网络连接

### 4. 端口被占用

**问题**: `Listen tcp :3000: bind: address already in use`

**解决方案**:
- 修改 `.env` 中的 `PORT` 配置
- 或停止占用端口的进程: `lsof -i :3000`

### 5. Docker 构建失败

**问题**: Docker build 失败

**解决方案**:
- 确保 Docker 正在运行
- 检查 Docker 版本 >= 20.10
- 清除 Docker 缓存: `docker system prune -a`

### 6. 权限问题

**问题**: `Permission denied`

**解决方案**:
- Linux/Mac: `chmod +x bin/server`
- 检查数据目录权限: `chmod -R 755 data/`

## 监控与日志

### 查看日志

```bash
# Docker 日志
docker-compose logs -f

# 或查看日志文件
tail -f logs/app.log
```

### 健康检查

```bash
curl http://localhost:3000/api/status
```

## 备份与恢复

### 数据库备份

```bash
# SQLite
cp data.db data.db.backup

# MySQL
mysqldump -u user -p dbname > backup.sql

# PostgreSQL
pg_dump -U user dbname > backup.sql
```

### 数据库恢复

```bash
# SQLite
cp data.db.backup data.db

# MySQL
mysql -u user -p dbname < backup.sql

# PostgreSQL
psql -U user dbname < backup.sql
```

## 更新与回滚

### 更新项目

```bash
# 拉取最新代码
git pull

# 重新构建
make build

# 重启服务
docker-compose restart
```

### 回滚版本

```bash
# 切换到指定版本
git checkout <commit-hash>

# 重新构建
make build

# 重启服务
docker-compose restart
```
