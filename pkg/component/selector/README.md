# Selector Component

## Overview

このセレクタコンポーネントは、Bubble Teaベースの対話型UIで、高さ制限付きの循環/非循環ナビゲーションをサポートする選択リストを提供する。Lipglossによるスタイリングと柔軟な設定オプションを提供する。

## Features

- 循環ナビゲーション: 上下移動でリストの先頭と末尾を循環
- カスタマイズ可能なスタイル: Lipglossによる見た目のカスタマイズ
- 柔軟な設定: プロンプト、表示モード、キーマップのカスタマイズ
- 表示範囲制限: 同時表示する項目数の指定
- Builder Pattern: メソッドチェーンによる設定

## Quick Start

### Basic Usage

```go
import (
    "github.com/cffnpwr/git-cz-go/pkg/component/selector"
    tea "github.com/charmbracelet/bubbletea"
)

// SelectItemインターフェースを実装
type StringItem string
func (s StringItem) ToString() string { return string(s) }

// アイテムを準備
items := []selector.SelectItem{
    StringItem("Option 1"),
    StringItem("Option 2"), 
    StringItem("Option 3"),
}

// セレクタモデルを初期化
model := selector.InitModel(items, 5).
    SetPrompt("Choose an option").
    SetShowSelectedItem(true).
    SetCyclic(true)

// Bubble Teaプログラムとして実行
p := tea.NewProgram(model)
p.Run()
```

### Getting Selected Item

選択されたアイテムを取得するには`GetSelectedItem()`メソッドを使用する。

```go
// 選択後にアイテムを取得
selectedItem := model.GetSelectedItem()
if selectedItem != nil {
    fmt.Printf("Selected: %s\n", selectedItem.ToString())
}
```

## API Reference

### Constructor

#### `InitModel(items []SelectItem, displaySize int) Model`

新しいセレクタモデルを初期化する。

- `items`: 選択可能なアイテムのスライス
- `displaySize`: 同時表示する最大項目数

### Configuration Methods

#### `SetPrompt(s string) Model`

プロンプトメッセージを設定する。

#### `SetShowSelectedItem(b bool) Model`  

選択後にアイテムを表示するかどうかを設定する。

#### `SetCyclic(b bool) Model`

循環ナビゲーションの有効/無効を設定する。

#### `SetKeyMap(km KeyMap) Model`

キーマップをカスタマイズする。

### State Methods

#### `GetSelectedItem() SelectItem`

現在選択されているアイテムを取得する。選択されていない場合は`nil`を返す。

### Interfaces

#### `SelectItem`

選択可能なアイテムに必要なインターフェース。

```go
type SelectItem interface {
    ToString() string
}
```

## Default Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | 上に移動 |
| `↓` / `j` | 下に移動 |
| `Enter` / `Space` | アイテムを選択 |
| `Ctrl+C` / `Esc` | 終了 |
