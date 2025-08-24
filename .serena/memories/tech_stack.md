# Technology Stack

## Dependencies (go.mod)
### Direct Dependencies:
- `github.com/charmbracelet/bubbletea v1.3.6` - TUIフレームワーク
- `github.com/spf13/cobra v1.9.1` - CLIフレームワーク
- `gopkg.in/yaml.v3 v3.0.1` - YAML設定ファイル管理

### Key Indirect Dependencies:
- `github.com/charmbracelet/lipgloss` - TUIスタイリング
- `github.com/charmbracelet/colorprofile` - カラープロファイル
- Various terminal handling libraries

## Development Tools
- **Runtime**: Go 1.24.6
- **Package Manager**: Go modules
- **Tool Manager**: mise (mise.toml)
- **Linter**: golangci-lint
- **Debugger**: delve (dlv)

## File Structure
```
├── main.go                 # エントリーポイント
├── cmd/                    # CLIコマンド定義
│   └── root.go
├── config/                 # 設定管理
│   ├── config.go           # 設定構造体とローダー
│   └── config.yaml         # デフォルト設定
├── internal/               # 内部パッケージ
│   ├── app/               # アプリケーションロジック
│   │   └── main.go
│   └── entity/            # データモデル
│       └── model.go       # Bubble Teaモデル
├── go.mod                 # Go modules定義
├── go.sum                 # 依存関係チェックサム
├── mise.toml              # 開発ツール管理
└── .editorconfig          # エディター設定
```