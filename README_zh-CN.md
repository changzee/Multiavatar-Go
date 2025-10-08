# multiavatar-go

[![Go Reference](https://pkg.go.dev/badge/github.com/changzee/multiavatar-go.svg)](https://pkg.go.dev/github.com/changzee/multiavatar-go)

一個用於生成多文化頭像的 Go 函式庫。這是對原始 JavaScript 函式庫 [Multiavatar](https://github.com/multiavatar/Multiavatar) 的移植版本。

該函式庫根據輸入的字串生成唯一的、確定性的頭像。其核心演算法與原始版本保持一致，確保相同的輸入字串始終生成相同的頭像。

![Logo](https://raw.githubusercontent.com/multiavatar/Multiavatar/main/logo.png)

## 功能特性

- **確定性演算法**: 對於相同的輸入，始終生成相同的頭像。
- **無依賴**: 函式庫是自包含的，不需要任何外部依賴。
- **函式選項模式**: 易於使用的 API，支援函式選項模式 (Functional Options)。
- **可自訂**: 支援生成帶有或不帶有背景的頭像。
- **執行緒安全**: 所有公開的函式都是為並行使用而設計的，執行緒安全。

## 安裝

使用 `go get` 來安裝此函式庫：

```bash
go get github.com/changzee/multiavatar-go
```

## 使用方法

在你的專案中匯入此函式庫：

```go
import "github.com/changzee/multiavatar-go"
```

### 基本範例

若要生成一個預設頭像，只需使用一個字串呼叫 `Generate` 函式。

```go
package main

import (
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// 為字串 "Binx Bond" 生成一個頭像
	svgCode := multiavatar.Generate("Binx Bond")

	// 將 SVG 儲存到檔案
	err := os.WriteFile("avatar.svg", []byte(svgCode), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

###帶有選項的範例

該函式庫使用函式選項模式來設定頭像的生成。例如，你可以使用 `WithoutBackground()` 選項來生成一個透明背景的頭像。

```go
package main

import (
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// 生成一個透明背景的頭像
	svgCode := multiavatar.Generate("John Doe", multiavatar.WithoutBackground())

	// 將 SVG 儲存到檔案
	err := os.WriteFile("avatar_transparent.svg", []byte(svgCode), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

## API 參考

### `Generate(input string, options ...Option) string`

這是生成 SVG 頭像的主要函式。

- `input string`: 一個 UTF-8 字串，作為生成頭像的種子。相同的輸入將始終產生相同的頭像。
- `options ...Option`: 一組可變參數的函式選項，用於自訂生成過程。

返回一個包含完整、格式正確的 SVG 頭像程式碼的字串。

### 選項

#### `WithoutBackground() Option`

此選項會移除頭像的彩色背景，使其變為透明。

## 授權

本專案採用 MIT 授權 - 詳情請參閱 [LICENSE](LICENSE) 檔案。原始的 Multiavatar 專案有其自身的授權，應予以遵守。
