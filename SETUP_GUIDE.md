# Go Aether - 独立 Module 完整指南

## 🎯 概述

这个方案创建一个完全独立的 Go Module 仓库,用户可以直接 `go get` 使用,无需安装 Rust 工具链。

## 📁 仓库结构

### 主仓库 (xiaozuhui/aether)

```
aether/                           # Rust 主仓库
├── src/
├── scripts/
│   └── build-all-libs.sh         # 构建所有平台的预编译库
├── Cargo.toml
└── target/
    └── releases/                 # 预编译库输出目录
        ├── darwin-arm64/
        │   └── libaether.a
        ├── darwin-amd64/
        │   └── libaether.a
        ├── linux-amd64/
        │   └── libaether.a
        └── windows-amd64/
            └── aether.lib
```

### Go Module 仓库 (xiaozuhui/go-aether)

```
go-aether/                        # 新仓库
├── go.mod                         # module github.com/xiaozuhui/go-aether
├── version.go                     # 版本号
├── aether.go                      # 主代码
├── lib_loader.go                  # 自动下载库
├── aether_test.go
├── scripts/
│   └── fetch-lib.sh               # 下载预编译库
├── Makefile
├── README.md
└── examples/
    └── basic/
        └── main.go
```

## 🚀 发布流程

### 1. 主仓库发布新版本

```bash
cd aether

# 1. 更新版本号
vim Cargo.toml
# version = "0.4.5"

# 2. 构建预编译库
./scripts/build-all-libs.sh

# 3. 创建 GitHub Release
gh release create v0.4.5 \
  --title "Aether v0.4.5" \
  --notes "Release notes here..." \
  target/releases/*/*/*.a

# 4. 推送 tag
git tag v0.4.5
git push origin v0.4.5
```

### 2. Go Module 发布新版本

```bash
cd go-aether

# 1. 更新版本号
vim version.go
# const Version = "v1.0.1"

# 2. 更新 CHANGELOG
vim CHANGELOG.md

# 3. 提交
git add .
git commit -m "chore: bump version to v1.0.1"

# 4. 打 tag
git tag v1.0.1

# 5. 推送
git push origin main
git push origin v1.0.1
```

## 👥 用户使用流程

### 安装

```bash
# 方式 1: 安装最新版本
go get github.com/xiaozuhui/go-aether@latest

# 方式 2: 安装特定版本
go get github.com/xiaozuhui/go-aether@v1.0.0
```

### 首次使用自动下载库

用户第一次 import 时,会自动下载预编译库:

```go
package main

import aether "github.com/xiaozuhui/go-aether"

func main() {
    // 首次调用会自动下载库到 ~/.aether/lib/
    engine := aether.New()
    defer engine.Close()

    result, _ := engine.Eval("Set X 10\n(X + 20)")
    println(result)
}
```

**输出:**
```
ℹ️  Aether Go 绑定 - 预编译库下载工具
ℹ️  正在查询最新版本...
✅  最新版本: 0.4.5
ℹ️  检测到平台: darwin-arm64
ℹ️  正在下载预编译库...
✅ 下载完成!
```

### 手动下载库

```bash
# 下载最新版本
go run github.com/xiaozuhui/go-aether/scripts/fetch-lib.sh

# 下载特定版本
go run github.com/xiaozuhui/go-aether/scripts/fetch-lib.sh -v 0.4.4

# 查看可用版本
go run github.com/xiaozuhui/go-aether/scripts/fetch-lib.sh --list
```

## 📦 文件清单

### 需要创建的文件

1. **go-aether/go.mod**
```go
module github.com/xiaozuhui/go-aether

go 1.21
```

2. **go-aether/version.go**
```go
package aether

// Version 是 Go 绑定的版本号
const Version = "v1.0.0"
```

3. **go-aether/aether.go**
   - 从 `bindings/go/aether.go` 复制,修改模块路径

4. **go-aether/lib_loader.go**
   - 自动下载库的逻辑 (上面已创建)

5. **go-aether/scripts/fetch-lib.sh**
   - 下载脚本 (上面已创建)

6. **go-aether/Makefile**
```makefile
.PHONY: test fetch-lib clean

# 运行测试
test: fetch-lib
	go test -v -race

# 手动下载库
fetch-lib:
	@./scripts/fetch-lib.sh

# 清理
clean:
	go clean
	rm -rf ~/.aether/lib

# 显示库信息
info:
	@go run -tags=info main.go
```

## 🔧 技术细节

### 自动下载机制

1. **init 钩子**: 包初始化时检查库是否存在
2. **自动检测**: 检测平台 (OS + 架构)
3. **智能查找**: 在 `~/.aether/lib/{platform}/` 中查找最新版本
4. **环境变量**: 设置 `AETHER_LIB_DIR` 供 CGO 使用

### 库查找路径

```
~/.aether/lib/
├── darwin-arm64/
│   ├── v0.4.4/
│   │   └── libaether.a
│   └── v0.4.5/
│       └── libaether.a
├── darwin-amd64/
│   └── ...
└── linux-amd64/
    └── ...
```

### CGO 集成

```c
/*
#cgo LDFLAGS: -L${AETHER_LIB_DIR} -laether
#cgo darwin LDFLAGS: -framework Security -framework CoreFoundation
*/
import "C"
```

`${AETHER_LIB_DIR}` 由 Go 代码在运行时设置。

## 📊 优势

| 特性 | 传统方式 | 本方案 |
|------|---------|--------|
| 用户安装 | 需要 Rust 工具链 | `go get` 即可 |
| 首次使用 | 需要编译 | 自动下载库 |
| 跨平台 | 需要交叉编译 | 预编译库 |
| 版本管理 | 依赖 Git tag | 独立版本号 |
| 用户体验 | ⭐⭐ 复杂 | ⭐⭐⭐⭐⭐ 简单 |

## 🎯 下一步

你需要:

1. **创建新仓库**: `go-aether`
2. **复制代码**: 从 `bindings/go` 复制文件
3. **修改 go.mod**: 改为 `module github.com/xiaozuhui/go-aether`
4. **添加脚本**: `lib_loader.go`, `fetch-lib.sh`
5. **测试**: 确保自动下载正常工作
6. **发布**: 推送到 GitHub,用户就可以 `go get` 了

需要我帮你创建完整的文件吗?
