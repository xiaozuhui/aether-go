// Aether 库文件下载工具
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("Aether 预编译库下载工具")
	fmt.Println()

	// 调用下载函数
	if err := fetchLibrary(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n下载完成！现在可以使用 aether-go 了。")
}

func fetchLibrary() error {
	fmt.Println("检测平台...")

	// 检测当前平台
	platform := detectPlatform()
	fmt.Printf("平台: %s\n", platform)

	// 检查是否为 macOS 平台
	if !isMacOS(platform) {
		fmt.Printf("\n注意: 当前仅提供 macOS 平台的预编译库。\n")
		fmt.Printf("其他平台需要从源码编译 Aether Rust 库。\n")
		fmt.Printf("\n编译步骤:\n")
		fmt.Printf("  1. 克隆 Aether 仓库: git clone https://github.com/xiaozuhui/aether.git\n")
		fmt.Printf("  2. 构建 Rust 库: cd aether && cargo build --release\n")
		fmt.Printf("  3. 复制库文件: mkdir -p lib/ && cp target/release/libaether.a lib/\n")
		return fmt.Errorf("不支持的平台: %s (仅支持 macOS)", platform)
	}

	// 创建 lib 目录
	libDir := "lib"
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("无法创建 lib 目录: %w", err)
	}

	// 构建下载文件名（包含架构信息）
	// 文件名格式: libaether-darwin-arm64.a 或 libaether-darwin-amd64.a
	libName := fmt.Sprintf("libaether-%s.a", platform)
	url := fmt.Sprintf("https://github.com/xiaozuhui/aether-go/releases/latest/download/%s", libName)

	// 本地保存路径（统一命名为 libaether.a）
	outputFile := fmt.Sprintf("%s/libaether.a", libDir)

	// 检查是否已存在
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Printf("库文件已存在: %s\n", outputFile)
		return nil
	}

	// 下载文件
	fmt.Printf("\n正在从 GitHub 下载...\n")
	fmt.Printf("URL: %s\n", url)

	// 使用 curl 或 wget 下载
	if err := downloadFile(url, outputFile); err != nil {
		return err
	}

	fmt.Printf("\n库文件已下载到: %s\n", outputFile)
	return nil
}

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

func isMacOS(platform string) bool {
	return platform == "darwin-amd64" || platform == "darwin-arm64"
}

func downloadFile(url, outputFile string) error {
	// 尝试使用 curl
	if _, err := exec.LookPath("curl"); err == nil {
		cmd := exec.Command("curl", "-L", "-f", "--progress-bar", url, "-o", outputFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// 尝试使用 wget
	if _, err := exec.LookPath("wget"); err == nil {
		cmd := exec.Command("wget", "--show-progress", "-O", outputFile, url)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return fmt.Errorf("需要 curl 或 wget 来下载文件")
}
