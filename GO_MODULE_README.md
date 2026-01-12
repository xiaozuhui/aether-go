# Go Aether - 独立 Go Module

Aether DSL 的 Go 语言绑定,独立仓库,开箱即用。

## 🚀 快速开始

### 安装

```bash
go get github.com/xiaozuhui/aether-go@latest
```

### 使用

```go
package main

import (
    "fmt"
    aether "github.com/xiaozuhui/aether-go"
)

func main() {
    engine := aether.New()
    defer engine.Close()

    result, err := engine.Eval(`
        Set X 10
        Set Y 20
        (X + Y)
    `)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // 30
}
```

## 📦 预编译库

首次使用时,会自动下载对应平台的预编译库:

- ✅ macOS (Intel & Apple Silicon)
- ✅ Linux (x86_64)
- ✅ Windows (x86_64)

库文件会被缓存到 `~/.aether/lib/`。

### 手动下载

```bash
# 从 GitHub Release 下载
./scripts/fetch-lib.sh
```

## 🔨 开发

### 前置要求

- Go 1.21+
- (可选) Rust - 用于重新编译本地库

### 构建预编译库

```bash
# 克隆主仓库
git clone https://github.com/xiaozuhui/aether.git
cd aether

# 编译所有平台的库
./scripts/build-all-libs.sh

# 文件会生成到 target/releases/
```

## 📝 许可证

Apache-2.0,与主仓库保持一致
