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
	// 这里调用包内的下载函数
	// 由于这是一个独立的命令行工具，我们需要直接实现下载逻辑

	// TODO: 实现实际的下载逻辑
	// 可以将 lib_loader.go 中的 FetchLibrary 逻辑移到这里
	// 或者创建一个包级别的下载函数

	fmt.Println("检测平台...")

	// 检测当前平台
	platform := detectPlatform()
	fmt.Printf("平台: %s\n", platform)

	// 创建 lib 目录
	libDir := "lib"
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("无法创建 lib 目录: %w", err)
	}

	// 构建 URL
	libName := "libaether.a"
	if platform == "windows-amd64" {
		libName = "aether.lib"
	}

	url := fmt.Sprintf("https://github.com/xiaozuhui/aether-go/releases/latest/download/%s/%s", platform, libName)
	outputFile := fmt.Sprintf("%s/%s", libDir, libName)

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
