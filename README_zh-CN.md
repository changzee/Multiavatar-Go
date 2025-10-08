# multiavatar-go

[![Go Reference](https://pkg.go.dev/badge/github.com/changzee/multiavatar-go.svg)](https://pkg.go.dev/github.com/changzee/multiavatar-go)

一个用于生成多文化头像的 Go 函式庫。这是对原始 JavaScript 函式庫 [Multiavatar](https://github.com/multiavatar/Multiavatar) 的 Go 语言移植版本。

该函式庫根据输入的字串生成唯一的、确定性的 SVG 头像。其核心演算法与原始版本保持一致，确保相同的输入字串始终生成相同的头像。

![Logo](https://raw.githubusercontent.com/multiavatar/Multiavatar/main/logo.png)

## 功能特性

- **确定性演算法**: 对于相同的输入，始终生成相同的头像。
- **功能丰富**: 支持性别预设、主题控制、部件更换、颜色自定义等多种高级选项。
- **无外部依赖**: 函式庫完全自包含。
- **函式选项模式**: 易于使用的 API，支持函式选项模式 (Functional Options Pattern)。
- **可自訂**: 支持生成带或不带背景的头像，并能移除任意部件。
- **执行绪安全**: 所有公开的函式都为并发使用而设计，是线程安全的。

## 安装

使用 `go get` 来安装此函式庫：

```bash
go get github.com/changzee/multiavatar-go
```

## 使用方法

### 基本用法

若要生成一个默认头像，只需提供一个字符串。

```go
package main

import (
	"log"
	"os"

	"github.com/changzee/multiavatar-go"
)

func main() {
	// 为字符串 "Binx Bond" 生成一个头像
	svgCode := multiavatar.Generate("Binx Bond")

	// 将 SVG 保存到文件
	err := os.WriteFile("avatar.svg", []byte(svgCode), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
```

### 高级用法

该函式庫使用函式选项模式来自定义头像的生成。你可以组合使用多个选项。

#### 示例 1: 生成透明背景的头像

使用 `WithoutBackground()` 选项来移除背景。

```go
svgCode := multiavatar.Generate("John Doe", multiavatar.WithoutBackground())
// 将其保存到 "avatar_transparent.svg"
```

#### 示例 2: 使用性别预设

使用 `WithGender()` 选项可以生成具有特定性别特征的头像（会影响发型、眼睛等的选择范围）。

```go
// 生成一个女性风格的头像
svgCodeFemale := multiavatar.Generate("Jane Doe", multiavatar.WithGender("female"))

// 生成一个男性风格的头像
svgCodeMale := multiavatar.Generate("John Doe", multiavatar.WithGender("male"))
```

#### 示例 3: 自定义部件和颜色

你可以组合使用多个选项来达到更精细的控制，例如移除顶部（头发），并自定义皮肤颜色。

```go
svgCode := multiavatar.Generate(
    "No Hair Avatar",
    multiavatar.WithoutPart("top"),          // 移除 "top" 部件（头发）
    multiavatar.WithSkinColor("#f2c280"),    // 设置皮肤颜色
)
```

## API 参考

### `Generate(input string, options ...Option) string`

这是生成 SVG 头像的主要函式。

- `input string`: 一个 UTF-8 字串，作为生成头像的种子。
- `options ...Option`: 一组可变参数的函式选项，用于自订生成过程。

返回一个包含完整 SVG 头像代码的字符串。

### 可用选项 (Options)

#### 背景与部件

- `WithoutBackground() Option`: 移除头像的彩色背景，使其变为透明。
- `WithoutPart(partName string) Option`: 禁用并移除一个指定的部件，例如 `"top"`, `"eyes"`, `"mouth"`。

#### 预设

- `WithGender(gender string) Option`: 应用性别预设，可选值为 `"female"` 或 `"male"` (以及别名如 `"f"`, `"m"`)。这会影响部件的选择范围以符合特定风格。

#### 主题控制

- `WithTheme(theme string) Option`: 强制所有部件使用同一个主题 (`"A"`, `"B"`, or `"C"`)。
- `WithPartTheme(partName, theme string) Option`: 为单个部件强制指定主题。
- `WithAllowedThemes(partName string, themesList []string) Option`: 限制单个部件只能从指定的主题列表中选择。

#### 版本控制

- `WithPartVersion(partName, partVersion string) Option`: 强制单个部件使用指定的版本（`"00"` 到 `"15"`）。
- `WithAllowedVersions(partName string, versions []string) Option`: 限制单个部件只能从指定的版本列表中选择。
- `WithAllowedHeadVersions(versions ...string) Option`: `WithAllowedVersions` 的便捷用法，专用于 `"head"` 部件。
- `WithAllowedEyesVersions(versions ...string) Option`: `WithAllowedVersions` 的便捷用法，专用于 `"eyes"` 部件。
- `WithAllowedTopVersions(versions ...string) Option`: `WithAllowedVersions` 的便捷用法，专用于 `"top"` 部件。

#### 颜色自定义

- `WithPartColors(partName string, colors []string) Option`: 为单个部件覆盖颜色数组。
- `WithSkinColor(hex string) Option`: 设置皮肤颜色（`"head"` 部件的主色）。
- `WithEnvColor(hex string) Option`: 设置背景颜色。
- `WithClothesColors(colors ...string) Option`: 设置衣物颜色。
- `WithTopColors(colors ...string) Option`: 设置头发/顶部颜色。
- `WithEyesColors(colors ...string) Option`: 设置眼睛颜色。
- `WithMouthColors(colors ...string) Option`: 设置嘴部颜色。

## 授权

本专案采用 MIT 授权。原始的 Multiavatar 专案有其自身的授权，使用时应予以遵守。