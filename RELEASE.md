# Go Module 发布指南

本文档说明如何将 `aether-go` 发布到 Go Module Proxy。

## 前置条件

1. ✅ `go.mod` 已正确配置为 `github.com/xiaozuhui/aether-go`
2. ✅ 代码已推送到 GitHub 仓库: https://github.com/xiaozuhui/aether-go
3. ✅ LICENSE 文件已添加 (Apache-2.0)
4. ✅ README.md 已更新，包含正确的安装说明

## 发布步骤

### 1. 确保代码已提交并推送

```bash
# 查看当前状态
git status

# 添加所有更改
git add .

# 提交
git commit -m "准备发布 v1.0.0

- 修复 go.mod 模块路径
- 添加 LICENSE 文件
- 更新 README 安装说明
- 修复 lib_loader.go 编译错误"

# 推送到远程仓库
git push origin main
```

### 2. 创建版本标签

Go Module Proxy 通过 Git 标签来识别版本。

```bash
# 创建版本标签（使用语义化版本号）
git tag v1.0.0

# 推送标签到远程仓库
git push origin v1.0.0
```

### 3. 验证发布

等待几秒到几分钟后，Go Module Proxy 会自动索引你的模块。

#### 检查方法 1: 使用 pkg.go.dev

访问: https://pkg.go.dev/github.com/xiaozuhui/aether-go

如果页面显示你的包信息，说明发布成功！

#### 检查方法 2: 使用 go list

```bash
# 在临时目录测试
mkdir /tmp/test-aether-go
cd /tmp/test-aether-go
go mod init test

# 尝试获取模块
go get github.com/xiaozuhui/aether-go@v1.0.0

# 查看是否成功
cat go.mod
```

#### 检查方法 3: 直接安装

```bash
# 创建测试项目
mkdir /tmp/test-aether-project
cd /tmp/test-aether-project

# 初始化模块
go mod init example.com/test

# 安装 aether-go
go get github.com/xiaozuhui/aether-go@v1.0.0
```

### 4. 用户如何使用

用户可以通过以下方式使用你的模块:

```bash
# 方式 1: 获取最新版本
go get github.com/xiaozuhui/aether-go@latest

# 方式 2: 获取特定版本
go get github.com/xiaozuhui/aether-go@v1.0.0

# 方式 3: 在 go.mod 中指定
require github.com/xiaozuhui/aether-go v1.0.0
```

在代码中导入:

```go
import "github.com/xiaozuhui/aether-go"

func main() {
    engine := aether.New()
    defer engine.Close()

    result, err := engine.Eval("Set X 10\n(X + 20)")
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // 30
}
```

## 版本管理

### 语义化版本号

推荐使用语义化版本号 (Semantic Versioning):

- **v1.0.0** - 第一个稳定版本
- **v1.0.1** - Bug 修复 (向后兼容)
- **v1.1.0** - 新功能 (向后兼容)
- **v2.0.0** - 重大更新 (不向后兼容)

### 发布新版本

```bash
# 1. 更新代码
git add .
git commit -m "feat: 添加新功能"

# 2. 创建新标签
git tag v1.1.0

# 3. 推送标签
git push origin main
git push origin v1.1.0
```

## 常见问题

### Q: 标签推送后多久能在 pkg.go.dev 上看到?

A: 通常需要几秒到几分钟。如果超过 10 分钟还没看到，可以访问 https://pkg.go.dev/github.com/xiaozuhui/aether-go/@v/list 查看所有可用版本。

### Q: 如何删除错误的标签?

A: Go Module Proxy 会缓存所有标签，即使删除也不会移除。正确的做法是发布新的标签修复问题。

```bash
# 删除本地标签
git tag -d v1.0.0

# 删除远程标签
git push origin :refs/tags/v1.0.0

# 创建正确的标签
git tag v1.0.1
git push origin v1.0.1
```

### Q: 用户需要先编译 Rust 库吗?

A: 是的。有两种方式:

1. **使用预编译库（推荐）**: 运行 `./scripts/fetch-lib.sh`
2. **从源码构建**: 克隆 aether 仓库并运行 `cargo build --release`

在 README.md 中已经包含了详细的安装说明。

### Q: 如何查看模块的所有版本?

A: 访问: https://pkg.go.dev/github.com/xiaozuhui/aether-go/@v/list

或者使用命令:

```bash
go list -m -versions github.com/xiaozuhui/aether-go
```

## 检查清单

发布前检查:

- [ ] `go.mod` 模块路径正确
- [ ] 所有代码已提交
- [ ] 代码可以正常编译 (`go build`)
- [ ] README.md 包含安装说明
- [ ] LICENSE 文件存在
- [ ] Git 标签已创建
- [ ] Git 标签已推送到远程
- [ ] 在 pkg.go.dev 上可以查到

## 下一步

发布成功后:

1. ✅ 在 pkg.go.dev 上验证
2. ✅ 在本地测试安装
3. ✅ 更新主 aether 仓库的 README，添加 Go 绑定链接
4. ✅ 考虑发布到其他平台（如 Reddit、Hacker News 等）

## 相关链接

- Go Module Proxy: https://proxy.golang.org/
- pkg.go.dev: https://pkg.go.dev/
- 语义化版本: https://semver.org/
- 你的仓库: https://github.com/xiaozuhui/aether-go
