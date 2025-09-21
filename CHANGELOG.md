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

### Fixed
- Resolved package naming conflicts between agent and dashboard models
- Eliminated duplicate TaskType constants and related structures
- Fixed import path inconsistencies across merged codebase

### Technical Details
- Successfully merged nezhahq/agent v0-final branch
- Maintained compatibility for both dashboard and agent components
- Preserved existing dashboard functionality without breaking changes
- Established monorepo architecture for unified maintenance
- 33 files changed with 3,893 additions and 17 deletions