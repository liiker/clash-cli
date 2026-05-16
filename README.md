# clash-cli

基于 Clash RESTful API 的命令行代理管理工具。

通过简洁的命令行界面管理本地或远程运行的 Clash 实例，支持查看节点、切换代理、测速等功能。

## 功能特性

- 📋 **查看代理列表** — 树状展示所有代理组与节点
- 🔄 **切换代理节点** — 一键切换代理组的活动节点
- ⚡ **节点测速** — 测试单个节点或整组节点的延迟
- 📌 **查看当前节点** — 查看所有代理组的活动节点
- ⚙️ **运行配置** — 查看和修改运行模式、日志级别等
- 🔄 **热重载** — 热重载配置文件，无需重启
- 🔗 **连接管理** — 查看活跃连接、关闭指定或全部连接
- 📦 **订阅管理** — 查看和更新代理提供者（订阅源）
- 🌐 **远程管理** — 支持连接远程 Clash 实例
- 🔐 **认证支持** — 支持 Clash RESTful API Secret 认证

## 安装

### 前置要求

- Go 1.18+
- 已运行的 Clash 实例，且开启了 RESTful API

### 编译安装

```bash
git clone https://github.com/liiker/clash-cli.git && cd clash-cli

# 默认编译（版本号为 dev）
go build -o clash-cli .

# 指定版本号编译
go build -ldflags "-X clash-cli/cmd.Version=v1.0.0" -o clash-cli .
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
| `--version` | 查看版本号 | |

### 查看版本号

```bash
clash-cli --version
```

### 查看所有代理节点

```bash
clash-cli proxies
```

### 切换代理节点

```bash
# 基本用法
clash-cli switch <代理组名> <节点名>

# 示例
clash-cli switch "Proxy" "🇭🇰 香港 01"
clash-cli switch "Streaming" "🇯🇵 日本 03"
```

### 测试节点延迟

```bash
# 测试单个节点
clash-cli delay "🇭🇰 香港 01"

# 自定义测速 URL 和超时时间
clash-cli delay "🇭🇰 香港 01" --url https://www.gstatic.com/generate_204 --timeout 5000

# 测试整个代理组
clash-cli delay-group "Proxy"
```

### 查看当前活动节点

```bash
# 查看所有代理组
clash-cli current

# 查看指定代理组
clash-cli current "Proxy"
```

### 查看和修改运行配置

```bash
# 查看当前配置
clash-cli config

# 切换为全局代理模式
clash-cli config --mode global

# 切换回规则分流模式
clash-cli config --mode rule

# 修改日志级别
clash-cli config --log-level debug

# 启用 IPv6
clash-cli config --ipv6=true
```

### 热重载配置文件

```bash
# 使用当前配置重新加载
clash-cli reload

# 指定配置文件路径
clash-cli reload /path/to/config.yaml
```

### 查看和管理活跃连接

```bash
# 查看所有活跃连接
clash-cli connections

# 关闭所有连接
clash-cli connections --close

# 关闭指定 ID 的连接
clash-cli connections --kill <connection-id>
```

### 管理代理提供者（订阅源）

```bash
# 查看所有提供者
clash-cli providers

# 更新指定提供者的订阅
clash-cli providers --update "provider-name"

# 对指定提供者执行健康检查
clash-cli providers --healthcheck "provider-name"
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
