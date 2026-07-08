package sensitiveword

import "strings"

// WordResultHandler 敏感词结果处理器接口，对应 Java 的 IWordResultHandler<R>。
//
// 不同实现将 WordResult 转换为不同类型的输出：
//
//   - Raw:   原始 WordResult
//   - Word:  截取出的敏感词字符串
//   - WordTags: 敏感词 + 标签信息
type WordResultHandler[R any] interface {
	// Handle 将 wordResult 转换为目标类型 R。
	Handle(wordResult *WordResult, wordContext *WordContext, originalText string) R
}

// WordTagsDto 敏感词 + 标签信息，对应 Java 的 WordTagsDto。
type WordTagsDto struct {
	word string
	tags []string
}

func (d *WordTagsDto) Word() string   { return d.word }
func (d *WordTagsDto) Tags() []string { return d.tags }
func (d *WordTagsDto) SetWord(s string)     { d.word = s }
func (d *WordTagsDto) SetTags(t []string)   { d.tags = t }

// ===== 辅助函数 =====

// getWordString 从原始文本中按 wordResult 的下标截取子串，对应 InnerWordCharUtils.getString。
func getWordString(text string, wordResult *WordResult) string {
	return getStringByRange(text, wordResult.StartIndex(), wordResult.EndIndex())
}

// getWordTags 获取敏感词的标签，对应 InnerWordTagUtils.tags。
//
// 先用原始词查询；若为空，则用格式化后的词再查一次。
func getWordTags(word string, wordContext *WordContext) []string {
	if word == "" {
		return nil
	}
	wordTag := wordContext.WordTag()
	if wordTag == nil {
		return nil
	}
	tags := wordTag.GetTag(word)
	if len(tags) > 0 {
		return tags
	}
	formatWord := formatString(word, wordContext)
	return wordTag.GetTag(formatWord)
}

// formatString 使用上下文的 WordFormat 格式化整段字符串，对应 InnerWordFormatUtils.format。
func formatString(original string, context *WordContext) string {
	if original == "" {
		return original
	}
	wordFormat := context.WordFormat()
	var sb strings.Builder
	sb.Grow(len(original))
	for _, c := range original {
		sb.WriteRune(wordFormat.Format(c, context))
	}
	return sb.String()
}

// formatWordList 批量格式化单词列表，对应 InnerWordFormatUtils.formatWordList。
func formatWordList(list []string, context *WordContext) []string {
	if len(list) == 0 {
		return nil
	}
	result := make([]string, 0, len(list))
	for _, word := range list {
		result = append(result, formatString(word, context))
	}
	return result
}

// ===== 实现 =====

// wordResultHandlerRaw 原样返回 WordResult，对应 WordResultHandlerRaw。
type wordResultHandlerRaw struct{}

// WordResultHandlerRaw 单例。
var WordResultHandlerRaw = &wordResultHandlerRaw{}

func (h *wordResultHandlerRaw) Handle(wordResult *WordResult, wordContext *WordContext, originalText string) *WordResult {
	return wordResult
}

// wordResultHandlerWord 截取敏感词字符串，对应 WordResultHandlerWord。
type wordResultHandlerWord struct{}

// WordResultHandlerWord 单例。
var WordResultHandlerWord = &wordResultHandlerWord{}

func (h *wordResultHandlerWord) Handle(wordResult *WordResult, wordContext *WordContext, originalText string) string {
	return getWordString(originalText, wordResult)
}

// wordResultHandlerWordTags 敏感词 + 标签，对应 WordResultHandlerWordTags。
type wordResultHandlerWordTags struct{}

// WordResultHandlerWordTags 单例。
var WordResultHandlerWordTags = &wordResultHandlerWordTags{}

func (h *wordResultHandlerWordTags) Handle(wordResult *WordResult, wordContext *WordContext, originalText string) *WordTagsDto {
	dto := &WordTagsDto{}
	word := getWordString(originalText, wordResult)
	tags := getWordTags(word, wordContext)
	// 若为空，尝试用命中的敏感词匹配（v0.25.1 bug105）
	if len(tags) == 0 {
		tags = getWordTags(wordResult.Word(), wordContext)
	}
	dto.word = word
	dto.tags = tags
	return dto
}

// ===== 工具类 =====

// WordResultHandlers 工具类，对应 Java 的 WordResultHandlers。
var WordResultHandlers = wordResultHandlersHelper{}

type wordResultHandlersHelper struct{}

// Raw 返回原样处理器。
func (wordResultHandlersHelper) Raw() WordResultHandler[*WordResult] {
	return WordResultHandlerRaw
}

// Word 返回敏感词字符串处理器。
func (wordResultHandlersHelper) Word() WordResultHandler[string] {
	return WordResultHandlerWord
}

// WordTags 返回敏感词+标签处理器。
func (wordResultHandlersHelper) WordTags() WordResultHandler[*WordTagsDto] {
	return WordResultHandlerWordTags
}
