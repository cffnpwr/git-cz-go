# Project Overview

## Purpose
Git CZ Go は、Conventional Commit規則に従ったコミットメッセージを作成するためのCLIツールです。BubbleTeaを使用したTUIインターフェースを提供し、対話的にコミットメッセージを作成できます。

## Tech Stack
- **言語**: Go 1.24.6
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea
- **設定管理**: YAML形式
- **開発ツール管理**: mise
- **リンター**: golangci-lint
- **デバッガー**: delve (dlv)

## Architecture
プロジェクトはClean Architectureに基づいた構造になっています：
- `main.go`: エントリーポイント
- `cmd/`: Cobraを使用したCLIコマンド定義
- `config/`: YAMLベースの設定管理
- `internal/`: 内部パッケージ
  - `internal/app/`: アプリケーションロジック
  - `internal/entity/`: Bubble Teaモデル定義

## Key Features
- Conventional Commits形式のコミットメッセージ生成
- YAMLベースの設定ファイル（config/config.yaml）
- TUIを使用した対話型インターフェース
- カスタマイズ可能なメッセージテンプレート