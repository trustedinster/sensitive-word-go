package sensitiveword

import "strings"

// WordReplace 敏感词替换策略接口，对应 Java 的 IWordReplace。
//
// 在替换场景中，每命中一个敏感词，核心引擎会调用本策略将敏感词
// 对应的原始文本替换为指定内容（如星号）。
type WordReplace interface {
	// Replace 将 wordResult 指向的敏感词片段替换后追加到 builder。
	// rawText 为原始全文；实现可按需读取上下文配置。
	Replace(builder *strings.Builder, rawText string, wordResult *WordResult, wordContext *WordContext)
}

// ===== 实现 =====

// wordReplaceChar 使用固定字符替换敏感词，对应 Java 的 WordReplaceChar。
type wordReplaceChar struct {
	replaceChar rune
}

// NewWordReplaceChar 创建以指定字符替换的策略。
func NewWordReplaceChar(replaceChar rune) WordReplace {
	return &wordReplaceChar{replaceChar: replaceChar}
}

// NewWordReplaceCharDefault 创建以默认字符 '*' 替换的策略。
func NewWordReplaceCharDefault() WordReplace {
	return NewWordReplaceChar(DefaultReplaceChar)
}

func (r *wordReplaceChar) Replace(builder *strings.Builder, rawText string, wordResult *WordResult, wordContext *WordContext) {
	wordLen := wordResult.EndIndex() - wordResult.StartIndex()
	for i := 0; i < wordLen; i++ {
		builder.WriteRune(r.replaceChar)
	}
}

// ===== 工具类 =====

// WordReplaces 工具类，对应 Java 的 WordReplaces。
var WordReplaces = wordReplacesHelper{}

type wordReplacesHelper struct{}

// Chars 返回以指定字符替换的策略。
func (wordReplacesHelper) Chars(c rune) WordReplace {
	return NewWordReplaceChar(c)
}

// CharsDefault 返回以默认字符 '*' 替换的策略。
func (wordReplacesHelper) CharsDefault() WordReplace {
	return NewWordReplaceCharDefault()
}

// Defaults 返回默认替换策略（'*'）。
func (wordReplacesHelper) Defaults() WordReplace {
	return NewWordReplaceCharDefault()
}
