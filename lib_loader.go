package aether

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

// 自动检查预编译库
func init() {
	if !IsLibraryAvailable() {
		// 不在 init 中自动下载，而是提供清晰的安装指引
		fmt.Fprintf(os.Stderr, "\n⚠️  Aether 库文件未找到\n\n")
		fmt.Fprintf(os.Stderr, "请运行以下命令下载预编译库：\n\n")
		fmt.Fprintf(os.Stderr, "  go run github.com/xiaozuhui/aether-go/cmd/fetch@latest\n\n")
		fmt.Fprintf(os.Stderr, "或者从源码构建：\n")
		fmt.Fprintf(os.Stderr, "  git clone https://github.com/xiaozuhui/aether.git\n")
		fmt.Fprintf(os.Stderr, "  cd aether && cargo build --release\n")
		fmt.Fprintf(os.Stderr, "  cp target/release/libaether.a <your-project>/lib/\n\n")
	}
}

// getLibDir 获取库文件目录
func getLibDir() (string, error) {
	// 优先使用环境变量
	if dir := os.Getenv("AETHER_LIB_DIR"); dir != "" {
		return dir, nil
	}

	// 查找项目根目录（包含 go.mod 的目录）
	moduleRoot := findModuleRoot()

	// 使用项目内的 lib 目录
	libDir := filepath.Join(moduleRoot, "lib")

	// 检查库文件是否存在
	libFile := filepath.Join(libDir, "libaether.a")
	if _, err := os.Stat(libFile); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("库文件未下载,请运行: go run github.com/xiaozuhui/aether-go/cmd/fetch@latest")
		}
		return "", fmt.Errorf("无法访问库目录: %w", err)
	}

	return libDir, nil
}

// detectPlatform 检测当前平台
func detectPlatform() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// 转换架构名称
	switch arch {
	case "amd64":
		// 保持不变
	case "arm64":
		// 保持不变
	default:
		arch = "amd64" // 默认
	}

	return fmt.Sprintf("%s-%s", os, arch)
}

// ensureLibrary 确保库文件存在
func ensureLibrary() error {
	libDir, err := getLibDir()
	if err != nil {
		return err
	}

	// 检查库文件是否存在
	libFile := filepath.Join(libDir, "libaether.a")
	if _, err := os.Stat(libFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("库文件不存在: %s", libFile)
		}
		return fmt.Errorf("无法访问库文件: %w", err)
	}

	// 设置 AETHER_LIB_DIR 环境变量供 CGO 使用
	os.Setenv("AETHER_LIB_DIR", libDir)

	return nil
}

// FetchLibrary 下载预编译库
func FetchLibrary() error {
	// 查找项目根目录
	moduleRoot := findModuleRoot()
	libDir := filepath.Join(moduleRoot, "lib")

	// 创建 lib 目录
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("无法创建 lib 目录: %w", err)
	}

	// 检测平台
	platform := detectPlatform()

	// 构建 URL
	libName := "libaether.a"
	if platform == "windows-amd64" {
		libName = "aether.lib"
	}

	// 下载最新版本的库
	url := fmt.Sprintf("https://github.com/xiaozuhui/aether-go/releases/latest/download/%s/%s", platform, libName)
	outputFile := filepath.Join(libDir, libName)

	// 检查是否已存在
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Printf("库文件已存在: %s\n", outputFile)
		return nil
	}

	// 下载文件
	fmt.Printf("正在从 %s 下载...\n", url)
	cmd := exec.Command("curl", "-L", "-f", "--progress-bar", url, "-o", outputFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("下载失败: %w\n请检查网络连接或手动下载: %s", err, url)
	}

	fmt.Printf("\n库文件已下载到: %s\n", outputFile)
	return nil
}

// findModuleRoot 查找模块根目录
func findModuleRoot() string {
	// 从当前文件路径向上查找 go.mod
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达根目录
			break
		}
		dir = parent
	}

	return "."
}

// GetLibraryInfo 获取库文件信息
func GetLibraryInfo() map[string]interface{} {
	moduleRoot := findModuleRoot()
	libDir := filepath.Join(moduleRoot, "lib")
	libFile := filepath.Join(libDir, "libaether.a")

	info := map[string]interface{}{
		"lib_dir": libDir,
		"platform": detectPlatform(),
	}

	// 检查库文件是否存在
	if fi, err := os.Stat(libFile); err == nil {
		info["size"] = fi.Size()
		info["modified"] = fi.ModTime()
		info["exists"] = true
	} else {
		info["exists"] = false
		info["error"] = err.Error()
	}

	return info
}

// IsLibraryAvailable 检查库是否可用
func IsLibraryAvailable() bool {
	libDir, err := getLibDir()
	if err != nil {
		return false
	}

	libFile := filepath.Join(libDir, "libaether.a")
	_, err = os.Stat(libFile)
	return err == nil
}

// EnsureLibraryError 提供更详细的库错误信息
func EnsureLibraryError() error {
	if IsLibraryAvailable() {
		return nil
	}

	moduleRoot := findModuleRoot()
	libDir := filepath.Join(moduleRoot, "lib")
	libFile := filepath.Join(libDir, "libaether.a")

	if _, err := os.Stat(libFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("库文件不存在: %s\n解决方法: 运行 go run github.com/xiaozuhui/aether-go/cmd/fetch@latest", libFile)
		}
		return fmt.Errorf("无法访问库文件: %w", err)
	}

	return nil
}

// SyscallStat 用于检查库文件
type SyscallStat_t struct {
	Dev uint64
	Ino uint64
	Mode uint32
	Nlink uint32
	Uid uint32
	Gid uint32
	Rdev uint64
	Size int64
	Atim int64
	Mtim int64
	Ctim int64
	Blksize int64
	Blocks int64
}

// getLibFileStat 获取库文件状态 (使用 syscall,避免 CGO 递归)
func getLibFileStat(path string) (*SyscallStat_t, error) {
	var stat SyscallStat_t
	err := syscall.Stat(path, (*syscall.Stat_t)(unsafe.Pointer(&stat)))
	return &stat, err
}
