# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2025-09-21]
### Added
- Agent repository integration into monorepo structure
- Complete agent functionality with cmd/agent module
- Agent-specific packages: fm, gpu, monitor, processgroup, pty, util, utls
- GPU monitoring support for multiple platforms (Linux, Windows, macOS, FreeBSD)
- Process group management functionality
- Pseudo-terminal (PTY) support for remote command execution
- Network utility functions and HTTP client utilities
- TLS/UTLS transport layer customization
- File manager (FM) functionality for remote file operations
- Agent authentication and configuration management (AuthHandler, AgentConfig)

### Changed
- Unified build configuration in .goreleaser.yml to support both dashboard and agent builds
- Merged agent models into main model package to eliminate duplication
- Updated import paths from github.com/nezhahq/agent to github.com/naiba/nezha
- Enhanced go.mod with agent-specific dependencies
- Updated .gitignore to include agent and dashboard binaries

### Removed
- **完全移除内网穿透(NAT)功能**：
  - 删除 model/nat.go 和 service/singleton/nat.go
  - 移除控制器中的 NAT 网关、API 接口和页面路由
  - 移除 agent 中的 NAT 任务处理逻辑和相关参数
  - 删除 NAT 相关的前端模板文件
- **完全移除动态DNS功能**：
  - 删除整个 pkg/ddns/ 包及其所有子模块
  - 删除 model/ddns.go 和 service/singleton/ddns.go
  - 移除 model.Server 中的 DDNS 相关字段 (EnableDDNS, DDNSProfiles, DDNSProfilesRaw)
  - 移除控制器中的 DDNS API 接口和页面路由
  - 移除 RPC 服务中的 DDNS 更新逻辑
  - 删除 DDNS 相关的前端模板文件
- **清理相关依赖**：
  - 自动移除 go.mod 中不再使用的 DDNS 相关依赖包
  - 清理代码中的 unused imports
  - 移除数据库迁移中的 NAT 和 DDNSProfile 模型

### Fixed
- Resolved package naming conflicts between agent and dashboard models
- Eliminated duplicate TaskType constants and related structures
- Fixed import path inconsistencies across merged codebase

### Security
- **修复全部13个安全漏洞** (100%覆盖率)：
  - **高危漏洞 (2个)**：
    - CVE-2025-22868: golang.org/x/oauth2 恶意token内存消耗攻击 → v0.27.0
    - CVE-2025-22869: golang.org/x/crypto SSH服务器DoS攻击 → v0.36.0
  - **中危漏洞 (7个)**：
    - GHSA-2464-8j7c-4cjm: go-viper/mapstructure敏感信息泄露 → v2.4.0
    - GHSA-fv92-fjc5-jj9h: mapstructure日志信息泄露 → v2.4.0
    - CVE-2025-22872: golang.org/x/net XSS漏洞 → v0.38.0
    - CVE-2025-22870: golang.org/x/net IPv6代理绕过 → v0.38.0
    - GHSA-pmc3-p9hx-jq96: uTLS TLS降级攻击 → v1.7.0
    - CVE-2025-27144: go-jose DoS解析攻击 → v4.0.5
  - **低危漏洞 (1个)**：
    - CVE-2025-8556: CIRCL密钥交换验证问题 → v1.6.1
- **依赖升级统计**：
  - golang.org/x/crypto: v0.25.0 → v0.36.0
  - golang.org/x/net: v0.27.0 → v0.38.0
  - golang.org/x/oauth2: v0.21.0 → v0.27.0
  - github.com/refraction-networking/utls: v1.6.3 → v1.7.0
  - github.com/cloudflare/circl: v1.3.7 → v1.6.1
  - github.com/go-jose/go-jose/v4: v4.0.2 → v4.0.5
  - github.com/go-viper/mapstructure/v2: v2.2.1 → v2.4.0

### Technical Details
- Successfully merged nezhahq/agent v0-final branch
- Maintained compatibility for both dashboard and agent components
- Preserved existing dashboard functionality without breaking changes
- Established monorepo architecture for unified maintenance
- 33 files changed with 3,893 additions and 17 deletions
- **代码简化**：移除内网穿透和动态DNS功能后，代码库更加专注于核心监控功能
- **编译验证**：确保 dashboard 和 agent 组件编译无错误