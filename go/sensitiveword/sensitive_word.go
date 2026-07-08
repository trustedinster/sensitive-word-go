package sensitiveword

import "strings"

// SensitiveWord 核心敏感词引擎，对应 Java 的 SensitiveWord / AbstractSensitiveWord。
//
// 提供 findAll / findFirst / replace / contains 四个核心方法。
// 所有方法均为无状态方法（状态由 WordContext 持有），可并发调用。
type SensitiveWord struct{}

// NewSensitiveWord 创建引擎实例。
func NewSensitiveWord() *SensitiveWord {
	return &SensitiveWord{}
}

// DefaultSensitiveWord 默认单例。
var DefaultSensitiveWord = NewSensitiveWord()

// FindAll 返回文本中所有敏感词结果，对应 ISensitiveWord.findAll。
// 采用 FailOver（全量遍历）模式。
func (s *SensitiveWord) FindAll(text string, context *WordContext) []*WordResult {
	if text == "" {
		return nil
	}
	return s.innerSensitiveWords(text, WordValidModeFailOver, context)
}

// FindFirst 返回第一个敏感词结果，对应 ISensitiveWord.findFirst。
// 采用 FailFast（快速返回）模式。未命中返回 nil。
func (s *SensitiveWord) FindFirst(text string, context *WordContext) *WordResult {
	if text == "" {
		return nil
	}
	results := s.innerSensitiveWords(text, WordValidModeFailFast, context)
	if len(results) == 0 {
		return nil
	}
	return results[0]
}

// Contains 判断文本是否包含敏感词，对应 ISensitiveWord.contains。
func (s *SensitiveWord) Contains(text string, context *WordContext) bool {
	return s.FindFirst(text, context) != nil
}

// Replace 将文本中所有敏感词替换为指定字符，对应 ISensitiveWord.replace。
func (s *SensitiveWord) Replace(text string, context *WordContext) string {
	if text == "" {
		return text
	}
	allList := s.FindAll(text, context)
	if len(allList) == 0 {
		return text
	}
	return s.doReplace(text, allList, context)
}

// doReplace 根据 allList 中的下标信息，将敏感词片段替换后拼接。
func (s *SensitiveWord) doReplace(target string, allList []*WordResult, context *WordContext) string {
	replace := context.WordReplace()
	var sb strings.Builder
	runes := []rune(target)
	startIndex := 0

	for _, wordResult := range allList {
		itemStartIx := wordResult.StartIndex()
		itemEndIx := wordResult.EndIndex()

		// 脱敏的左边
		if startIndex < itemStartIx {
			sb.WriteString(string(runes[startIndex:itemStartIx]))
		}

		// 脱敏部分
		replace.Replace(&sb, target, wordResult, context)

		// 更新结尾
		if itemEndIx > startIndex {
			startIndex = itemEndIx
		}
	}

	// 最后部分
	if startIndex < len(runes) {
		sb.WriteString(string(runes[startIndex:]))
	}

	return sb.String()
}

// innerSensitiveWords 核心匹配逻辑，对应 Java SensitiveWord.innerSensitiveWords。
func (s *SensitiveWord) innerSensitiveWords(text string, modeEnum WordValidMode, context *WordContext) []*WordResult {
	sensitiveCheck := context.WordCheck()
	var resultList []*WordResult

	// 预格式化整段文本，得到字符映射表
	characterCharacterMap := context.WordFormatText().Format(text, context)
	checkContext := &InnerSensitiveWordContext{
		originalText:      text,
		wordContext:       context,
		modeEnum:          WordValidModeFailOver,
		formatCharMapping: characterCharacterMap,
	}
	wordResultCondition := context.WordResultCondition()

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		checkResult := sensitiveCheck.SensitiveCheck(i, checkContext)
		wordLengthAllow := checkResult.WordLengthResult().WordAllowLen()
		wordLengthDeny := checkResult.WordLengthResult().WordDenyLen()

		// 命中的白名单长度小于黑名单，保存敏感词
		if wordLengthAllow < wordLengthDeny {
			wordResult := &WordResult{
				startIndex: i,
				endIndex:   i + wordLengthDeny,
				checkType:  checkResult.CheckType(),
				word:       checkResult.WordLengthResult().WordDeny(),
			}

			// 结果条件过滤
			if wordResultCondition == nil ||
				wordResultCondition.Match(wordResult, text, modeEnum, context) {
				resultList = append(resultList, wordResult)
				if modeEnum == WordValidModeFailFast {
					break
				}
			}
			// 增加 i 的步长（-1 因为循环自增 1）
			i += wordLengthDeny - 1
		} else {
			// 白名单长度 >= 黑名单长度，跳过白名单个字符
			skip := wordLengthAllow - 1
			if skip < 0 {
				skip = 0
			}
			i += skip
		}
	}

	return resultList
}

// ===== 泛型便捷方法 =====

// FindAllWithHandler 返回所有敏感词并经 handler 转换后的结果列表。
func FindAllWithHandler[R any](s *SensitiveWord, text string, context *WordContext, handler WordResultHandler[R]) []R {
	wordResults := s.FindAll(text, context)
	if len(wordResults) == 0 {
		return nil
	}
	result := make([]R, 0, len(wordResults))
	for _, wr := range wordResults {
		result = append(result, handler.Handle(wr, context, text))
	}
	return result
}

// FindFirstWithHandler 返回第一个敏感词并经 handler 转换后的结果。
func FindFirstWithHandler[R any](s *SensitiveWord, text string, context *WordContext, handler WordResultHandler[R]) R {
	var zero R
	wr := s.FindFirst(text, context)
	if wr == nil {
		return zero
	}
	return handler.Handle(wr, context, text)
}
