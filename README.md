# sensitive-word-go

[sensitive-word-go](https://github.com/houbb/sensitive-word-go) 基于 DFA 算法实现的高性能敏感词工具，是 Java 项目 [sensitive-word](https://github.com/houbb/sensitive-word) 的 Go 重构实现。

[![Open Source Love](https://badges.frapsoft.com/os/v2/open-source.svg?v=103)](https://github.com/houbb/sensitive-word-go)
[![](https://img.shields.io/badge/license-Apache2-FF0080.svg)](https://github.com/houbb/sensitive-word-go/blob/master/LICENSE.txt)

## 创作目的

基于 DFA 算法实现，目前敏感词库内容收录 6W+（源文件 18W+，经过一次删减）。

后期将进行持续优化和补充敏感词库，并进一步提升算法的性能。

## 特性

- 6W+ 词库，且不断优化更新

- 基于 fluent-api 实现，使用优雅简洁

- 基于 DFA 算法，高性能匹配

- 支持敏感词的判断、返回、脱敏等常见操作

- 支持常见的格式转换

全角半角互换、英文大小写互换、数字常见形式的互换、中文繁简体互换、英文常见形式的互换、忽略重复词等

- 支持敏感词检测、邮箱检测、数字检测、网址检测、IPV4等

- 支持自定义替换策略

- 支持用户自定义敏感词和白名单

- 支持数据的数据动态更新（用户自定义），实时生效

- 支持敏感词的标签接口+内置分类实现

- 支持跳过一些特殊字符，让匹配更灵活

- 支持黑白名单单个的新增/删除，无需全量初始化

- 支持词匹配模式的两种模式（快速失败 / 全量遍历）

- 繁简体转换基于 [OpenCC-Go](https://github.com/yanmingcao/opencc-go) 实现，纯 Go，无需外部数据文件

## 变更日志

[CHANGE_LOG.md](https://github.com/houbb/sensitive-word-go/blob/master/CHANGE_LOG.md)

# 快速开始

## 准备

- Go 1.21+

## 引入

```bash
go get github.com/houbb/sensitive-word-go
```

## 核心方法

### 方式一：使用默认实例（懒加载，首次调用自动初始化）

```go
package main

import (
    "fmt"
    "github.com/houbb/sensitive-word-go"
)

func main() {
    text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"

    // 判断是否包含敏感词
    if sensitiveword.Default().Contains(text) {
        fmt.Println("包含敏感词")
    }

    // 返回第一个敏感词
    fmt.Println(sensitiveword.Default().FindFirst(text))

    // 返回所有敏感词
    fmt.Println(sensitiveword.Default().FindAll(text))

    // 默认替换策略（* 替换）
    fmt.Println(sensitiveword.Default().Replace(text))
}
```

### 方式二：自定义配置（fluent-api 风格）

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetIgnoreCase(true).
    SetIgnoreChineseStyle(true).
    SetEnableEmailCheck(true).
    SetEnableNumCheck(true).
    Init()
defer bs.Destroy()

if bs.Contains("联系 test@example.com") {
    fmt.Println("包含敏感信息")
}
```

## 核心方法说明

| 方法 | 参数 | 返回值 | 说明 |
|:-----|:-----|:-------|:-----|
| Contains(target) | 待验证的字符串 | 布尔值 | 验证字符串是否包含敏感词 |
| FindAll(target) | 待验证的字符串 | 字符串切片 | 返回字符串中所有敏感词（已去重保序） |
| FindFirst(target) | 待验证的字符串 | 字符串 | 返回字符串中第一个敏感词，未命中返回空字符串 |
| FindAllRaw(target) | 待验证的字符串 | []*WordResult | 返回所有敏感词的原始结果（含下标、类别信息） |
| FindFirstRaw(target) | 待验证的字符串 | *WordResult | 返回第一个敏感词的原始结果，未命中返回 nil |
| FindAllWordTags(target) | 待验证的字符串 | []*WordTagsDto | 返回所有敏感词及其标签 |
| FindFirstWordTags(target) | 待验证的字符串 | *WordTagsDto | 返回第一个敏感词及其标签 |
| Replace(target) | 待验证的字符串 | 字符串 | 返回脱敏后的字符串 |
| Tags(word) | 敏感词字符串 | []string | 返回敏感词的标签列表 |

### 判断是否包含敏感词

```go
text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"
fmt.Println(sensitiveword.Default().Contains(text)) // true
```

### 返回第一个敏感词

```go
text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"
fmt.Println(sensitiveword.Default().FindFirst(text)) // 五星红旗
```

### 返回所有敏感词

```go
text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"
fmt.Println(sensitiveword.Default().FindAll(text)) // [五星红旗 毛主席 天安门]
```

### 默认的替换策略

```go
text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"
fmt.Println(sensitiveword.Default().Replace(text)) // ****迎风飘扬，***的画像屹立在***前。
```

### 指定替换的内容

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordReplace(sensitiveword.NewWordReplaceChar('#')).
    Init()
text := "这是一个敏感词"
fmt.Println(bs.Replace(text)) // 这是一个##词
```

### 获取敏感词标签

```go
bs := sensitiveword.NewSensitiveWordBs().Init()
tags := bs.Tags("博彩")
fmt.Println(tags) // [3]
```

# 更多特性

## 样式处理

### 忽略大小写

```go
text := "fuCK the bad words."
fmt.Println(sensitiveword.Default().FindFirst(text)) // fuCK
```

### 忽略半角圆角

```go
text := "ｆｕｃｋ the bad words."
fmt.Println(sensitiveword.Default().FindFirst(text)) // ｆｕｃｋ
```

### 忽略数字的写法

```go
text := "这个是我的微信：9⓿二肆⁹₈③⑸⒋➃㈤㊄"
bs := sensitiveword.NewSensitiveWordBs().SetEnableNumCheck(true).Init()
fmt.Println(bs.FindAll(text)) // [9⓿二肆⁹₈③⑸⒋➃㈤㊄]
```

### 忽略繁简体

```go
text := "我爱我的祖国和五星紅旗。"
fmt.Println(sensitiveword.Default().FindAll(text)) // [五星紅旗]
```

### 忽略英文的书写格式

```go
text := "Ⓕⓤc⒦ the bad words"
fmt.Println(sensitiveword.Default().FindAll(text)) // [Ⓕⓤc⒦]
```

### 忽略重复词

```go
text := "ⒻⒻⒻfⓤuⓤ⒰cⓒ⒦ the bad words"
bs := sensitiveword.NewSensitiveWordBs().SetIgnoreRepeat(true).Init()
fmt.Println(bs.FindAll(text)) // [ⒻⒻⒻfⓤuⓤ⒰cⓒ⒦]
```

## 更多检测策略

| 方法 | 说明 | 默认值 |
|:-----|:-----|:------|
| SetWordCheckWord | 敏感词检测策略 | `WordChecks.Word()` |
| SetWordCheckNum | 数字检测策略 | `WordChecks.Num()` |
| SetWordCheckEmail | 邮箱检测策略 | `WordChecks.Email()` |
| SetWordCheckUrl | URL检测策略 | `WordChecks.Url()` |
| SetWordCheckIpv4 | ipv4检测策略 | `WordChecks.Ipv4()` |

### 邮箱检测

邮箱等个人信息，默认未启用。

```go
text := "楼主好人，邮箱 sensitiveword@xx.com"
bs := sensitiveword.NewSensitiveWordBs().SetEnableEmailCheck(true).Init()
fmt.Println(bs.FindAll(text)) // [sensitiveword@xx.com]
```

### 连续数字检测

一般用于过滤手机号/QQ等广告信息，默认未启用。

支持通过 `SetNumCheckLen(长度)` 自定义检测的长度。

```go
text := "你懂得：12345678"

// 默认检测 8 位
bs := sensitiveword.NewSensitiveWordBs().SetEnableNumCheck(true).Init()
fmt.Println(bs.FindAll(text)) // [12345678]

// 指定数字的长度，避免误杀
bs2 := sensitiveword.NewSensitiveWordBs().SetEnableNumCheck(true).SetNumCheckLen(9).Init()
fmt.Println(bs2.FindAll(text)) // []
```

### 网址检测

用于过滤常见的网址信息，默认未启用。

```go
text := "点击链接 https://www.baidu.com 查看答案"
bs := sensitiveword.NewSensitiveWordBs().SetEnableUrlCheck(true).Init()
fmt.Println(bs.FindAll(text)) // [https://www.baidu.com]
fmt.Println(bs.Replace(text)) // 点击链接 ********************* 查看答案
```

内置支持不需要 http 协议的前缀检测：

```go
text := "点击链接 https://www.baidu.com 查看答案，当然也可以是 baidu.com、www.baidu.com"
bs := sensitiveword.NewSensitiveWordBs().
    SetEnableUrlCheck(true).
    SetWordCheckUrl(sensitiveword.WordChecks.UrlNoPrefix()).
    Init()
fmt.Println(bs.FindAll(text)) // [www.baidu.com baidu.com www.baidu.com]
```

### IPV4 检测

避免用户通过 ip 绕过网址检测等，默认未启用。

```go
text := "个人网站，如果网址打不开可以访问 127.0.0.1。"
bs := sensitiveword.NewSensitiveWordBs().SetEnableIpv4Check(true).Init()
fmt.Println(bs.FindAll(text)) // [127.0.0.1]
```

# 引导类特性配置

## 配置方法

用户可以使用 `SensitiveWordBs` 进行如下定义：

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetIgnoreCase(true).
    SetIgnoreWidth(true).
    SetIgnoreNumStyle(true).
    SetIgnoreChineseStyle(true).
    SetIgnoreEnglishStyle(true).
    SetIgnoreRepeat(false).
    SetEnableNumCheck(false).
    SetEnableEmailCheck(false).
    SetEnableUrlCheck(false).
    SetEnableIpv4Check(false).
    SetEnableWordCheck(true).
    SetWordFailFast(true).
    SetWordCheckNum(sensitiveword.WordChecks.Num()).
    SetWordCheckEmail(sensitiveword.WordChecks.Email()).
    SetWordCheckUrl(sensitiveword.WordChecks.Url()).
    SetWordCheckIpv4(sensitiveword.WordChecks.Ipv4()).
    SetWordCheckWord(sensitiveword.WordChecks.Word()).
    SetNumCheckLen(8).
    SetWordTag(sensitiveword.WordTags.None()).
    SetCharIgnore(sensitiveword.SensitiveWordCharIgnores.Defaults()).
    SetWordResultCondition(sensitiveword.WordResultConditions.AlwaysTrue()).
    Init()
defer bs.Destroy()

text := "五星红旗迎风飘扬，毛主席的画像屹立在天安门前。"
fmt.Println(bs.Contains(text)) // true
```

## 配置说明

| 序号 | 方法 | 说明 | 默认值 |
|:-----|:-----|:-----|:------|
| 1 | SetIgnoreCase | 忽略大小写 | true |
| 2 | SetIgnoreWidth | 忽略半角圆角 | true |
| 3 | SetIgnoreNumStyle | 忽略数字的写法 | true |
| 4 | SetIgnoreChineseStyle | 忽略中文的书写格式 | true |
| 5 | SetIgnoreEnglishStyle | 忽略英文的书写格式 | true |
| 6 | SetIgnoreRepeat | 忽略重复词 | false |
| 7 | SetEnableNumCheck | 是否启用数字检测 | false |
| 8 | SetEnableEmailCheck | 是否启用邮箱检测 | false |
| 9 | SetEnableUrlCheck | 是否启用链接检测 | false |
| 10 | SetEnableIpv4Check | 是否启用IPv4检测 | false |
| 11 | SetEnableWordCheck | 是否启用敏感单词检测 | true |
| 12 | SetNumCheckLen | 数字检测自定义长度 | 8 |
| 13 | SetWordTag | 词对应的标签 | none |
| 14 | SetCharIgnore | 忽略的字符 | none |
| 15 | SetWordResultCondition | 匹配的敏感词额外加工 | 恒为真 |
| 16 | SetWordCheckNum | 数字检测策略 | `WordChecks.Num()` |
| 17 | SetWordCheckEmail | 邮箱检测策略 | `WordChecks.Email()` |
| 18 | SetWordCheckUrl | URL检测策略 | `WordChecks.Url()` |
| 19 | SetWordCheckIpv4 | ipv4检测策略 | `WordChecks.Ipv4()` |
| 20 | SetWordCheckWord | 敏感词检测策略 | `WordChecks.Word()` |
| 21 | SetWordReplace | 替换策略 | `WordReplaces.Defaults()` |
| 22 | SetWordFailFast | 敏感词匹配模式是否快速返回 | true |
| 23 | SetWordFormatText | 文本整体级别的格式化处理策略 | `WordFormatTexts.Defaults()` |

## SetWordFailFast 敏感词匹配快速失败模式

默认情况下，`SetWordFailFast(true)`。匹配时快速返回，性能较好，但有时不太符合人的直觉。

### failOver 模式

尽可能找到最长的匹配词。

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordFailFast(false).
    SetWordDeny(sensitiveword.NewWordDenyChain(/* 自定义黑名单 */)).
    Init()

text := "他的世界它的世界和她的世界都不是我的也不是我的世界"
fmt.Println(bs.FindAll(text)) // 尽可能匹配最长词
```

## 内存资源的释放

```go
bs := sensitiveword.NewSensitiveWordBs().Init()
// 后续因为一些原因移除了对应信息，希望释放内存。
bs.Destroy()
```

## 针对单个黑名单词的新增/删除，无需全量初始化

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordAllow(sensitiveword.WordAllows.Empty()).
    SetWordDeny(sensitiveword.WordDenys.Empty()).
    Init()

text := "测试一下新增敏感词，验证一下删除和新增对不对"
fmt.Println(bs.FindAll(text)) // []

// 新增
bs.AddWord("测试")
bs.AddWord("新增")
fmt.Println(bs.FindAll(text)) // [测试 新增 新增]

// 删除
bs.RemoveWord("新增")
fmt.Println(bs.FindAll(text)) // [测试]
```

## 针对单个白名单词的新增/删除，无需全量初始化

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordAllow(sensitiveword.WordAllows.Empty()).
    SetWordDeny(/* 包含 测试、新增 的黑名单 */).
    Init()

text := "测试一下新增敏感词白名单"
fmt.Println(bs.FindAll(text)) // [测试 新增]

// 新增白名单
bs.AddWordAllow("测试")
bs.AddWordAllow("新增")
fmt.Println(bs.FindAll(text)) // []

// 删除白名单
bs.RemoveWordAllow("测试")
fmt.Println(bs.FindAll(text)) // [测试]
```

# wordResultCondition-针对匹配词进一步判断

有时候我们可能希望对匹配的敏感词进一步限制，比如虽然我们定义了【av】作为敏感词，但是不希望【have】被匹配。

系统内置的策略在 `WordResultConditions` 工具类中：

| 实现 | 说明 |
|:-----|:-----|
| Defaults | 默认策略（英文全词匹配） |
| AlwaysTrue | 恒为真 |
| EnglishWordMatch | 英文单词全词匹配 |
| EnglishWordNumMatch | 英文单词/数字全词匹配 |
| WordTags | 满足特定标签的，比如只关注【广告】标签 |
| Chains | 支持指定多个条件，同时满足（AND 语义） |

```go
text := "I have a nice day。"
bs := sensitiveword.NewSensitiveWordBs().
    SetWordDeny(/* 包含 av 的黑名单 */).
    SetWordResultCondition(sensitiveword.WordResultConditions.EnglishWordMatch()).
    Init()
fmt.Println(bs.FindAll(text)) // []  "av" 在 "have" 中不算全词匹配
```

# 忽略字符

我们可以在敏感词中间加一些字符跳过检测，比如【傻!@#$帽】。通过指定特殊字符的跳过集合，忽略掉这些无意义的字符即可。

```go
text := "傻@冒，狗+东西"

// 默认因为有特殊字符分割，无法识别
bs := sensitiveword.NewSensitiveWordBs().Init()
fmt.Println(bs.FindAll(text)) // []

// 指定忽略的字符策略
bs2 := sensitiveword.NewSensitiveWordBs().
    SetCharIgnore(sensitiveword.SensitiveWordCharIgnores.SpecialChars()).
    Init()
fmt.Println(bs2.FindAll(text)) // [傻@冒 狗+东西]
```

# 敏感词标签

有时候我们希望对敏感词加一个分类标签：比如社情、暴力等等。

## 标签接口

```go
type WordTag interface {
    GetTag(word string) []string
}
```

## 内置实现

`WordTags` 工具类提供了常用策略：

| 方法 | 说明 |
|:-----|:-----|
| None | 空实现 |
| Map | 根据 map 初始化 |
| Lines | 根据字符串列表初始化 |
| LinesWithSplit | 字符串列表，自定义分隔符 |
| System | 系统内置实现，从内嵌 sensitive_word_tags.txt 加载（懒加载） |
| Defaults | 默认策略，目前为 System |

### 格式约定

敏感词标签的格式默认约定如下 `敏感词 tag1,tag2`，代表这 `敏感词` 的标签为 tag1 和 tag2。

```
五星红旗 政治,国家
```

### 系统内置实现（默认效果）

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordTag(sensitiveword.WordTags.System()).
    Init()
tags := bs.Tags("博彩")
fmt.Println(tags) // [3]
```

数字含义列表：

```
0 政治
1 毒品
2 色情
3 赌博
4 违法
```

### 获取敏感词的同时获取标签

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordTag(sensitiveword.WordTags.Lines([]string{"天安门 政治,国家,地址"})).
    Init()

result := bs.FindAllWordTags("天安门")
// result[0].Word() == "天安门"
// result[0].Tags() == [政治 国家 地址]
```

# 动态加载（用户自定义）

## 接口说明

### WordDeny

返回的列表表示这个词是一个敏感词（黑名单）。

```go
type WordDeny interface {
    Deny() []string
}
```

### WordAllow

返回的列表表示这个词不是一个敏感词（白名单）。

```go
type WordAllow interface {
    Allow() []string
}
```

## 配置使用

### 系统的默认配置

```go
bs := sensitiveword.NewSensitiveWordBs().
    SetWordDeny(sensitiveword.WordDenys.Defaults()).
    SetWordAllow(sensitiveword.WordAllows.Defaults()).
    Init()
```

### 同时配置多个

- 多个敏感词：`WordDenys.Chains()` 方法，将多个实现合并为同一个 WordDeny。
- 多个白名单：`WordAllows.Chains()` 方法，将多个实现合并为同一个 WordAllow。

```go
wordDeny := sensitiveword.NewWordDenyChain(sensitiveword.WordDenys.Defaults(), myWordDeny)
wordAllow := sensitiveword.NewWordAllowChain(sensitiveword.WordAllows.Defaults(), myWordAllow)

bs := sensitiveword.NewSensitiveWordBs().
    SetWordDeny(wordDeny).
    SetWordAllow(wordAllow).
    Init()
```

# 系统内置词库

| 文件 | 说明 | 默认加载类 |
|:-----|:-----|:---------|
| `sensitive_word_allow.txt` | 内置自定义白名单词库 | `WordAllowSystem` |
| `sensitive_word_deny.txt` | 内置自定义黑名单词库 | `WordDenySystem` |
| `sensitive_word_dict.txt` | 内置黑名单词库（6W+） | `WordDenySystem` |
| `sensitive_word_dict_en.txt` | 内置黑名单英文词库 | `WordDenySystem` |
| `sensitive_word_tags.txt` | 内置敏感词标签词库（4W+） | `WordTagSystem` |

词库通过 `go:embed` 内嵌到二进制文件中，无需外部数据文件。

如不希望使用内置词库，可将 `WordAllow`、`WordDeny`、`WordTag` 替换为自己的实现。

# 繁简体转换

繁简体转换基于 [OpenCC-Go](https://github.com/yanmingcao/opencc-go) 实现，使用其内嵌的 t2s 预设（纯 Go，无需外部数据文件），覆盖完整字符级繁简映射。

可通过 `RegisterT2SMapping` / `RegisterT2SMappingBatch` 在 OpenCC 结果之上追加自定义映射：

```go
// 注册额外的繁简映射
sensitiveword.RegisterT2SMapping('臺', '台')

// 批量注册
sensitiveword.RegisterT2SMappingBatch(map[rune]rune{
    '臺': '台',
    '灣': '湾',
})
```

# 拓展阅读

## 敏感词系列

[01-开源敏感词工具入门使用](https://houbb.github.io/2020/01/07/sensitive-word-00-overview)

[02-如何实现一个敏感词工具？违禁词实现思路梳理](https://houbb.github.io/2020/01/07/sensitive-word-01-intro)

[05-敏感词之 DFA 算法(Trie Tree 算法)详解](https://houbb.github.io/2020/01/07/sensitive-word-04-dfa)
