# Confirm Component

## Overview

この確認コンポーネントは、Bubble Teaベースの対話型UIで、Yes/No形式の確認ダイアログを提供する。Lipglossによるスタイリングと柔軟なキーバインドオプションを提供する。

## Features

- バイナリ選択: Yes/Noの2択選択
- 複数入力方法: キーナビゲーション、直接入力、確認キーをサポート
- カスタマイズ可能なスタイル: Lipglossによる見た目のカスタマイズ
- 柔軟なキーマップ: キーバインドのカスタマイズ
- Builder Pattern: メソッドチェーンによる設定

## Quick Start

### Basic Usage

```go
import (
    "github.com/cffnpwr/git-cz-go/pkg/component/confirm"
    tea "github.com/charmbracelet/bubbletea"
)

// 確認モデルを初期化
model := confirm.InitModel().
    SetPrompt("Do you want to continue?")

// Bubble Teaプログラムとして実行
p := tea.NewProgram(model)
p.Run()
```

### Getting Confirmation Result

確認結果を取得するには`GetValue()`メソッドを使用する。

```go
// 確認後に結果を取得
result := model.GetValue()
if result {
    fmt.Println("User confirmed: Yes")
} else {
    fmt.Println("User confirmed: No")
}
```

## API Reference

### Constructor

#### `InitModel() Model`

新しい確認モデルを初期化する。デフォルトでは`false`（No）が選択された状態で開始される。

### Configuration Methods

#### `SetPrompt(s string) Model`

プロンプトメッセージを設定する。デフォルトは"Confirm"。

#### `SetKeyMap(km KeyMap) Model`

キーマップをカスタマイズする。

### State Methods

#### `GetValue() bool`

現在の確認状態を取得する。`true`の場合はYes、`false`の場合はNoを表す。

## Default Key Bindings

| Key                           | Action                 |
| ----------------------------- | ---------------------- |
| `←` / `→` / `Tab` / `h` / `l` | 選択を切り替え         |
| `y`                           | Yes を直接選択して確定 |
| `n`                           | No を直接選択して確定  |
| `Enter` / `Space`             | 現在の選択を確定       |
| `Ctrl+C` / `q` / `Esc`        | 終了                   |

## Visual Design

- 選択されていない項目: 通常のボーダーで表示
- 選択された項目: 紫色の背景（`#bb9af7`）、白色のテキスト、太字で表示
- レイアウト: YesとNoのボタンが水平に並んで表示
