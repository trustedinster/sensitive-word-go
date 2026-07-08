package sensitiveword

import "unicode"

// WordFormat 单字符格式化策略接口，对应 Java 的 IWordFormat。
//
// 实现方负责将原始字符按某种规则（忽略大小写、全角半角、数字样式等）
// 映射为统一形式的字符，以便 DFA 树以规范化字符存储与匹配。
//
// 实现应为线程安全（无状态）。
type WordFormat interface {
	// Format 将 original 字符格式化后返回。
	Format(original rune, context *WordContext) rune
}

// ===== 基础实现 =====

// wordFormatNone 无处理，原样返回，对应 WordFormatNone。
type wordFormatNone struct{}

// WordFormatNone 单例。
var WordFormatNone = &wordFormatNone{}

func (f *wordFormatNone) Format(original rune, context *WordContext) rune {
	return original
}

// wordFormatIgnoreCase 忽略大小写，对应 WordFormatIgnoreCase。
type wordFormatIgnoreCase struct{}

// WordFormatIgnoreCase 单例。
var WordFormatIgnoreCase = &wordFormatIgnoreCase{}

func (f *wordFormatIgnoreCase) Format(original rune, context *WordContext) rune {
	// 等价于 Java Character.toLowerCase
	return unicode.ToLower(original)
}

// wordFormatIgnoreWidth 全角转半角，对应 WordFormatIgnoreWidth。
type wordFormatIgnoreWidth struct{}

// WordFormatIgnoreWidth 单例。
var WordFormatIgnoreWidth = &wordFormatIgnoreWidth{}

func (f *wordFormatIgnoreWidth) Format(original rune, context *WordContext) rune {
	return toHalfWidth(original)
}

// wordFormatIgnoreChineseStyle 繁体转简体，对应 WordFormatIgnoreChineseStyle。
type wordFormatIgnoreChineseStyle struct{}

// WordFormatIgnoreChineseStyle 单例。
var WordFormatIgnoreChineseStyle = &wordFormatIgnoreChineseStyle{}

func (f *wordFormatIgnoreChineseStyle) Format(original rune, context *WordContext) rune {
	return toSimple(original)
}

// ===== 组合实现 =====

// wordFormatArray 按顺序依次调用一组格式化策略，对应 WordFormatArray。
type wordFormatArray struct {
	formats []WordFormat
}

// NewWordFormatArray 创建一个组合格式化策略。
// formats 不能为空。
func NewWordFormatArray(formats []WordFormat) WordFormat {
	if len(formats) == 0 {
		return WordFormatNone
	}
	// 拷贝避免外部修改
	cp := make([]WordFormat, len(formats))
	copy(cp, formats)
	return &wordFormatArray{formats: cp}
}

func (f *wordFormatArray) Format(original rune, context *WordContext) rune {
	c := original
	for _, fmt := range f.formats {
		c = fmt.Format(c, context)
	}
	return c
}

// ===== 工具类 WordFormats =====

// WordFormats 对应 Java 的 WordFormats 工具类，提供常用策略获取与链式组合。
var WordFormats = wordFormatsHelper{}

type wordFormatsHelper struct{}

// None 返回无处理策略。
func (wordFormatsHelper) None() WordFormat { return WordFormatNone }

// IgnoreCase 返回忽略大小写策略。
func (wordFormatsHelper) IgnoreCase() WordFormat { return WordFormatIgnoreCase }

// IgnoreWidth 返回忽略全角半角策略。
func (wordFormatsHelper) IgnoreWidth() WordFormat { return WordFormatIgnoreWidth }

// IgnoreChineseStyle 返回忽略繁简体策略。
func (wordFormatsHelper) IgnoreChineseStyle() WordFormat { return WordFormatIgnoreChineseStyle }

// IgnoreNumStyle 返回忽略数字样式策略。
func (wordFormatsHelper) IgnoreNumStyle() WordFormat { return WordFormatIgnoreNumStyle }

// IgnoreEnglishStyle 返回忽略英文字母样式策略。
func (wordFormatsHelper) IgnoreEnglishStyle() WordFormat { return WordFormatIgnoreEnglishStyle }

// Chains 将多个策略组合为链式调用，对应 WordFormats.chains。
// 传入空切片时返回 None。
func (wordFormatsHelper) Chains(formats ...WordFormat) WordFormat {
	if len(formats) == 0 {
		return WordFormatNone
	}
	return NewWordFormatArray(formats)
}

// ChainsList 基于切片组合，对应 WordFormats.chains(List)。
func (h wordFormatsHelper) ChainsList(formats []WordFormat) WordFormat {
	if len(formats) == 0 {
		return WordFormatNone
	}
	return NewWordFormatArray(formats)
}

// ===== 组合器 WordFormatCombine =====

// WordFormatCombine 根据上下文配置构建组合格式化策略，对应 Java 的 WordFormatCombine。
//
// 组合顺序与 Java 保持一致：
//
//	ignoreEnglishStyle -> ignoreCase -> ignoreWidth -> ignoreNumStyle -> ignoreChineseStyle
type WordFormatCombine struct{}

// NewWordFormatCombine 创建组合器。
func NewWordFormatCombine() *WordFormatCombine {
	return &WordFormatCombine{}
}

// InitWordFormat 依据上下文构造组合后的 WordFormat。
func (c *WordFormatCombine) InitWordFormat(context *WordContext) WordFormat {
	formats := make([]WordFormat, 0, 5)
	if context.IgnoreEnglishStyle() {
		formats = append(formats, WordFormats.IgnoreEnglishStyle())
	}
	if context.IgnoreCase() {
		formats = append(formats, WordFormats.IgnoreCase())
	}
	if context.IgnoreWidth() {
		formats = append(formats, WordFormats.IgnoreWidth())
	}
	if context.IgnoreNumStyle() {
		formats = append(formats, WordFormats.IgnoreNumStyle())
	}
	if context.IgnoreChineseStyle() {
		formats = append(formats, WordFormats.IgnoreChineseStyle())
	}
	return WordFormats.ChainsList(formats)
}
