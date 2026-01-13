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

// 自动下载预编译库
func init() {
	if err := ensureLibrary(); err != nil {
		// 不要在 init 中 panic,只记录警告
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
		fmt.Fprintf(os.Stderr, "请运行: ./scripts/fetch-lib.sh\n")
	}
}

// getLibDir 获取库文件目录
func getLibDir() (string, error) {
	// 优先使用环境变量
	if dir := os.Getenv("AETHER_LIB_DIR"); dir != "" {
		return dir, nil
	}

	// 检测平台
	platform := detectPlatform()

	// 查找已下载的库
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户目录: %w", err)
	}

	libBaseDir := filepath.Join(homeDir, ".aether", "lib", platform)

	// 查找最新的版本目录
	entries, err := os.ReadDir(libBaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			// 目录不存在,需要下载
			return "", fmt.Errorf("库文件未下载,请运行: ./scripts/fetch-lib.sh")
		}
		return "", fmt.Errorf("无法读取库目录: %w", err)
	}

	// 找到最新的版本目录
	var latestVersion string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 0 && entry.Name()[0] == 'v' {
			if latestVersion == "" || entry.Name() > latestVersion {
				latestVersion = entry.Name()
			}
		}
	}

	if latestVersion == "" {
		return "", fmt.Errorf("未找到已下载的库文件,请运行: ./scripts/fetch-lib.sh")
	}

	return filepath.Join(libBaseDir, latestVersion), nil
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
func FetchLibrary(version string) error {
	// 查找脚本
	scriptPath := filepath.Join(findModuleRoot(), "scripts", "fetch-lib.sh")

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("下载脚本不存在: %s", scriptPath)
	}

	// 构建命令
	cmd := exec.Command(scriptPath)
	if version != "" && version != "latest" {
		cmd.Args = append(cmd.Args, "-v", version)
	}

	// 设置标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行脚本
	return cmd.Run()
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
	libDir, err := getLibDir()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	info := map[string]interface{}{
		"lib_dir": libDir,
		"platform": detectPlatform(),
	}

	// 获取库文件大小
	libFile := filepath.Join(libDir, "libaether.a")
	if fi, err := os.Stat(libFile); err == nil {
		info["size"] = fi.Size()
		info["modified"] = fi.ModTime()
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

	libDir, err := getLibDir()
	if err != nil {
		return fmt.Errorf("库未初始化: %w\n解决方法:\n1. 运行: ./scripts/fetch-lib.sh\n2. 或从源码编译: git clone https://github.com/xiaozuhui/aether.git && cd aether && cargo build --release", err)
	}

	libFile := filepath.Join(libDir, "libaether.a")
	if _, err := os.Stat(libFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("库文件不存在: %s\n解决方法: 运行 ./scripts/fetch-lib.sh", libFile)
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
