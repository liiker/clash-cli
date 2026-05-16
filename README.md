# clash-cli

基于 Clash RESTful API 的命令行代理管理工具。

通过简洁的命令行界面管理本地或远程运行的 Clash 实例，支持查看节点、切换代理、测速等功能。

## 功能特性

- 📋 **查看代理列表** — 树状展示所有代理组与节点
- 🔄 **切换代理节点** — 一键切换代理组的活动节点
- ⚡ **节点测速** — 测试单个节点或整组节点的延迟
- 📌 **查看当前节点** — 查看所有代理组的活动节点
- 🌐 **远程管理** — 支持连接远程 Clash 实例
- 🔐 **认证支持** — 支持 Clash RESTful API Secret 认证

## 安装

### 前置要求

- Go 1.18+
- 已运行的 Clash 实例，且开启了 RESTful API

### 编译安装

```bash
git clone <repo-url> && cd clash-cli
go build -o clash-cli .
```

也可以直接安装到 `$GOPATH/bin`：

```bash
go install .
```

## 使用方法

### 全局参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--api-url` | Clash RESTful API 地址 | `http://127.0.0.1:9090` |
| `--secret` | API 认证密钥 | 空 |

### 查看所有代理节点

```bash
clash-cli proxies
```

### 切换代理节点

```bash
# 基本用法
clash-cli switch <代理组名> <节点名>

# 示例
clash-cli switch "🔰 选择节点" "🇭🇰 香港A01"
clash-cli switch "Proxy" "HK-Node-01"
```

### 测试节点延迟

```bash
# 测试单个节点
clash-cli delay "🇭🇰 香港A01"

# 自定义测速 URL 和超时时间
clash-cli delay "🇭🇰 香港A01" --url https://www.gstatic.com/generate_204 --timeout 5000

# 测试整个代理组
clash-cli delay-group "🔰 选择节点"
```

### 查看当前活动节点

```bash
# 查看所有代理组
clash-cli current

# 查看指定代理组
clash-cli current "🔰 选择节点"
```

### 连接远程 Clash

```bash
clash-cli proxies --api-url http://192.168.1.100:9090 --secret your-secret-key
```

## 输出说明

### 延迟质量评定

| 标识 | 延迟范围 | 评价 |
|------|----------|------|
| 🟢 优秀 | < 200ms | 推荐 |
| 🟡 一般 | 200–500ms | 可用 |
| 🟠 较差 | 500–1000ms | 不推荐 |
| 🔴 很差 | > 1000ms | 极差 |
| 🔴 不可达 | 超时 | 无法连接 |

## 开发

```
clash-cli/
├── main.go            # 程序入口
├── client/
│   ├── client.go      # HTTP 客户端封装
│   └── proxy.go       # API 数据结构与方法
├── cmd/
│   └── root.go        # Cobra 命令定义
├── go.mod
└── go.sum
```

## License

MIT
