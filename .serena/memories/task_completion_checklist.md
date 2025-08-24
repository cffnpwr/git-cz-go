# Task Completion Checklist

タスク完了時に実行すべきコマンドのチェックリスト

## Code Quality Checks
```bash
# 1. コードフォーマット
go fmt ./...

# 2. インポート整理（必要に応じて）
goimports -w .

# 3. リント実行
golangci-lint run

# 4. モジュール整理
go mod tidy
```

## Testing
```bash
# 全テストの実行
go test ./...

# テストカバレッジ確認（必要に応じて）
go test -cover ./...
```

## Build Verification
```bash
# ビルドが通ることを確認
go build -o git-cz .

# 基本的な動作確認
./git-cz --help
```

## Pre-commit Checks
1. `.editorconfig` の設定に従っているか確認
2. 新しい依存関係が適切に `go.mod` に記録されているか
3. パッケージ構造が規約に従っているか
4. エラーハンドリングが適切に実装されているか

## Documentation Updates
- 新しい機能や重要な変更がある場合は、適切なコメントを追加
- 設定ファイルの変更がある場合は、`config/config.yaml` のサンプルを更新

## Development Tools Sync
```bash
# mise設定の同期確認
mise install

# 開発ツールのバージョン確認
mise current
```