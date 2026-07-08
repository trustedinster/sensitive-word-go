package sensitiveword

// WordFormatText 文本整体格式化策略接口，对应 Java 的 IWordFormatText。
//
// 与 WordFormat（单字符）不同，本接口针对整段文本生成
// "原始字符 -> 格式化字符" 的映射表，便于在匹配前一次性预处理，
// 减少重复的单字符格式化开销。
type WordFormatText interface {
	// Format 对 text 做整体格式化，返回发生变化的字符映射。
	// 未发生变化的字符不出现在结果中。
	Format(text string, context *WordContext) map[rune]rune
}

// wordFormatTextDefault 默认实现，对应 Java 的 WordFormatTextDefault。
//
// 遍历文本每个字符，调用上下文中的 WordFormat 进行格式化，
// 仅记录发生变化的字符映射；若 WordFormat 等价于 None，则直接返回空映射。
type wordFormatTextDefault struct{}

// WordFormatTextDefault 单例。
var WordFormatTextDefault = &wordFormatTextDefault{}

func (f *wordFormatTextDefault) Format(text string, context *WordContext) map[rune]rune {
	if text == "" {
		return nil
	}
	wordFormat := context.WordFormat()
	// 不需要处理的场景：WordFormat 为 None 单例时直接返回空映射
	if wordFormat == WordFormatNone {
		return nil
	}
	m := make(map[rune]rune)
	for _, c := range text {
		mc := wordFormat.Format(c, context)
		if c != mc {
			m[c] = mc
		}
	}
	return m
}

// WordFormatTexts 工具类，对应 Java 的 WordFormatTexts。
var WordFormatTexts = wordFormatTextsHelper{}

type wordFormatTextsHelper struct{}

// Defaults 返回默认的 WordFormatText 实现。
func (wordFormatTextsHelper) Defaults() WordFormatText {
	return WordFormatTextDefault
}
