# Code Style and Conventions

## Editor Configuration (.editorconfig)
```
root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

[*.go]
indent_style = tab
```

## Go Specific Rules
- **インデント**: Goファイルはタブを使用（.editorconfigで定義）
- **その他のファイル**: 2スペースインデント
- **改行**: LF（Unix形式）
- **文字エンコーディング**: UTF-8
- **末尾空白**: 削除する
- **ファイル末尾**: 改行を挿入

## Package Structure Conventions
```
cmd/           # CLIコマンド定義
config/        # 設定関連
internal/      # 内部パッケージ（外部公開しない）
  app/         # アプリケーションロジック
  entity/      # データモデル/エンティティ
```

## Import Organization
1. 標準ライブラリ
2. 外部ライブラリ
3. プロジェクト内パッケージ

## Error Handling Pattern
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err)
    os.Exit(1)
}
```

## Configuration Pattern
- YAML設定ファイル使用
- `config.LoadConfig(path)` パターン
- 構造体ベースの設定管理

## Bubble Tea Model Pattern
```go
type Model struct {
    // fields
}

func (m *Model) Init() tea.Cmd { ... }
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (m *Model) View() string { ... }
```