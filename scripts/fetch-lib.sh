#!/bin/bash

# Aether 预编译库下载脚本
# 自动检测平台并下载对应的预编译库

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 默认版本 (如果用户不指定)
DEFAULT_VERSION="latest"

# 库文件目录
LIB_DIR="${HOME}/.aether/lib"

# GitHub API 基础 URL
GITHUB_API="https://api.github.com/repos/xiaozuhui/aether"

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

# 获取最新版本号
get_latest_version() {
    info "正在查询最新版本..."

    local latest_version
    latest_version=$(curl -s "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$latest_version" ]; then
        error "无法获取最新版本号"
        return 1
    fi

    # 移除 v 前缀
    echo "$latest_version" | sed 's/^v//'
}

# 检测平台
detect_platform() {
    local os
    local arch

    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            os="windows"
            ;;
        *)
            error "不支持的操作系统: $(uname -s)"
            return 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            error "不支持的架构: $(uname -m)"
            return 1
            ;;
    esac

    echo "${os}-${arch}"
}

# 下载库文件
download_lib() {
    local version=$1
    local platform=$2
    local lib_name="libaether.a"

    if [ "$platform" = "windows-amd64" ]; then
        lib_name="aether.lib"
    fi

    # 构建 URL
    local url
    if [ "$version" = "latest" ]; then
        # 使用最新版本的下载 URL
        url="https://github.com/xiaozuhui/aether/releases/latest/download/${platform}/${lib_name}"
    else
        url="https://github.com/xiaozuhui/aether/releases/download/v${version}/${platform}/${lib_name}"
    fi

    local output_dir="${LIB_DIR}/${platform}/v${version}"
    local output_file="${output_dir}/${lib_name}"

    # 检查是否已存在
    if [ -f "$output_file" ]; then
        success "库文件已存在: $output_file"
        echo "$output_file"
        return 0
    fi

    # 创建目录
    mkdir -p "$output_dir"

    # 下载
    info "正在下载预编译库..."
    info "版本: v${version}"
    info "平台: $platform"

    if command -v curl >/dev/null 2>&1; then
        # 显示进度条
        curl -L -f --progress-bar "$url" -o "$output_file" || {
            rm -f "$output_file"
            echo ""
            error "下载失败!"
            echo ""
            echo "可能的原因:"
            echo "1. 网络连接问题"
            echo "2. 版本 v${version} 不存在"
            echo "3. 平台 ${platform} 暂不支持"
            echo ""
            echo "请检查: https://github.com/xiaozuhui/aether/releases"
            echo ""
            echo "或从源码编译:"
            echo "  git clone https://github.com/xiaozuhui/aether.git"
            echo "  cd aether"
            echo "  cargo build --release"
            return 1
        }
    elif command -v wget >/dev/null 2>&1; then
        wget --show-progress -O "$output_file" "$url" || {
            rm -f "$output_file"
            error "下载失败!"
            return 1
        }
    else
        error "需要 curl 或 wget 来下载文件"
        return 1
    fi

    success "库文件已下载到: $output_file"
    echo "$output_file"
}

# 显示帮助
show_help() {
    cat << EOF
用法: $0 [选项]

下载 Aether 预编译库

选项:
    -v, --version VERSION    指定版本号 (默认: latest)
    -p, --platform PLATFORM   指定平台 (默认: 自动检测)
    -l, --list               列出可用版本
    -h, --help               显示此帮助信息

示例:
    $0                          # 下载最新版本
    $0 -v 0.4.4                # 下载指定版本
    $0 -v 0.4.4 -p darwin-arm64 # 下载指定平台

支持的平台:
    - darwin-amd64    (macOS Intel)
    - darwin-arm64    (macOS Apple Silicon)
    - linux-amd64     (Linux x86_64)
    - windows-amd64   (Windows x86_64)

EOF
}

# 列出可用版本
list_versions() {
    info "查询可用版本..."

    local versions
    versions=$(curl -s "${GITHUB_API}/releases" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | head -10)

    if [ -z "$versions" ]; then
        error "无法获取版本列表"
        return 1
    fi

    echo ""
    echo "可用的版本:"
    echo "$versions" | while read -r version; do
        echo "  - $version"
    done
    echo ""
}

# 主函数
main() {
    local VERSION="$DEFAULT_VERSION"
    local PLATFORM=""

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -p|--platform)
                PLATFORM="$2"
                shift 2
                ;;
            -l|--list)
                list_versions
                exit 0
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done

    echo "Aether Go 绑定 - 预编译库下载工具"
    echo ""

    # 如果是 latest,查询实际版本号
    if [ "$VERSION" = "latest" ]; then
        VERSION=$(get_latest_version) || exit 1
        success "最新版本: v${VERSION}"
    fi

    # 检测平台
    if [ -z "$PLATFORM" ]; then
        PLATFORM=$(detect_platform) || exit 1
        info "检测到平台: $PLATFORM"
    else
        info "使用平台: $PLATFORM"
    fi
    echo ""

    # 下载库
    local lib_file
    lib_file=$(download_lib "$VERSION" "$PLATFORM") || exit 1
    echo ""

    # 输出结果
    success "下载完成!"
    echo ""
    echo "库文件位置: $lib_file"
    echo ""
    echo "在 Go 代码中使用时,会自动找到此库文件"
    echo "如果需要手动指定,请设置环境变量:"
    echo ""
    echo "export CGO_LDFLAGS=\"-L${LIB_DIR}/${PLATFORM}/v${VERSION}/\""
    echo ""
}

# 运行
main "$@"
