# Development Commands

## Build & Run
```bash
# プロジェクトをビルド
go build -o git-cz .

# 直接実行
go run main.go

# インストール (GOPATH/binに配置)
go install
```

## Development Tools
```bash
# 開発ツールのセットアップ (mise使用)
mise install

# 利用可能なツール一覧
mise list

# 現在のツールバージョン確認
mise current
```

## Code Quality
```bash
# リント実行
golangci-lint run

# コードフォーマット
go fmt ./...

# インポート整理
goimports -w .

# モジュール整理
go mod tidy
```

## Testing
```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/...

# テストカバレッジ
go test -cover ./...

# 詳細なテスト結果
go test -v ./...
```

## Debugging
```bash
# delveデバッガーで実行
dlv debug

# 特定のブレークポイントでデバッグ
dlv debug -- [args]
```

## Dependencies Management
```bash
# 依存関係の追加
go get [package]

# 依存関係の更新
go get -u ./...

# モジュール情報確認
go mod graph
go mod why [package]
```