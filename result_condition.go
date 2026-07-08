package sensitiveword

// WordResultCondition 结果匹配条件接口，对应 Java 的 IWordResultCondition。
//
// 用于在命中敏感词后做二次过滤：只有满足条件的命中才会被保留。
// 典型场景：英文敏感词需全词匹配、按标签过滤等。
type WordResultCondition interface {
	// Match 判断 wordResult 是否满足保留条件。
	Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool
}

// ===== 实现 =====

// wordResultConditionAlwaysTrue 恒为真，对应 WordResultConditionAlwaysTrue。
type wordResultConditionAlwaysTrue struct{}

// WordResultConditionAlwaysTrue 单例。
var WordResultConditionAlwaysTrue = &wordResultConditionAlwaysTrue{}

func (c *wordResultConditionAlwaysTrue) Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool {
	return true
}

// wordResultConditionEnglishWordMatch 英文单词需全词匹配，对应 WordResultConditionEnglishWordMatch。
//
// 若命中片段含非英文非空格字符（如中文、数字），直接通过；
// 若全为英文（可含空格），则检查前后字符是否也为英文，是则视为非全词匹配而拒绝。
type wordResultConditionEnglishWordMatch struct{}

// WordResultConditionEnglishWordMatch 单例。
var WordResultConditionEnglishWordMatch = &wordResultConditionEnglishWordMatch{}

func (c *wordResultConditionEnglishWordMatch) Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool {
	startIndex := wordResult.StartIndex()
	endIndex := wordResult.EndIndex()
	runes := []rune(text)
	length := len(runes)

	// 检查匹配片段是否全为英文（可含空格）
	for i := startIndex; i < endIndex; i++ {
		c := runes[i]
		if !isEnglish(c) && !isSpace(c) {
			// 含非英文字符（如中文、数字），直接通过
			return true
		}
	}
	// 全英文：检查前后字符是否为英文
	if startIndex > 0 && isEnglish(runes[startIndex-1]) {
		return false
	}
	if endIndex < length && isEnglish(runes[endIndex]) {
		return false
	}
	return true
}

// wordResultConditionEnglishWordNumMatch 英文单词和数字需全词匹配，对应 WordResultConditionEnglishWordNumMatch。
type wordResultConditionEnglishWordNumMatch struct{}

// WordResultConditionEnglishWordNumMatch 单例。
var WordResultConditionEnglishWordNumMatch = &wordResultConditionEnglishWordNumMatch{}

func (c *wordResultConditionEnglishWordNumMatch) Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool {
	startIndex := wordResult.StartIndex()
	endIndex := wordResult.EndIndex()
	runes := []rune(text)
	length := len(runes)

	// 前一字符为字母或数字则拒绝
	if startIndex > 0 && isDigitOrLetter(runes[startIndex-1]) {
		return false
	}
	// 后一字符为字母或数字则拒绝
	if endIndex < length && isDigitOrLetter(runes[endIndex]) {
		return false
	}
	// 判断当前片段是否为纯字母数字；含其他字符则直接通过
	for i := startIndex; i < endIndex; i++ {
		if !isDigitOrLetter(runes[i]) {
			return true
		}
	}
	return true
}

// wordResultConditionWordTagsMatch 按标签过滤，对应 WordResultConditionWordTagsMatch。
//
// 仅当敏感词的标签与期望标签集合有交集时才保留。
type wordResultConditionWordTagsMatch struct {
	tags map[string]bool
	wordTag WordTag
}

// NewWordResultConditionWordTagsMatch 创建按标签过滤的条件。
// tags 为期望的标签列表；wordTag 从上下文获取（可为 nil，此时取 context.WordTag()）。
func NewWordResultConditionWordTagsMatch(tags []string, wordTag WordTag) WordResultCondition {
	m := make(map[string]bool, len(tags))
	for _, t := range tags {
		m[t] = true
	}
	return &wordResultConditionWordTagsMatch{tags: m, wordTag: wordTag}
}

func (c *wordResultConditionWordTagsMatch) Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool {
	if wordResult.Word() == "" {
		return false
	}
	wt := c.wordTag
	if wt == nil {
		wt = context.WordTag()
	}
	if wt == nil {
		return false
	}
	wordTags := wt.GetTag(wordResult.Word())
	for _, t := range wordTags {
		if c.tags[t] {
			return true
		}
	}
	return false
}

// wordResultConditionChains 链式条件，需同时满足所有条件，对应 WordResultConditionInit 的链式实现。
type wordResultConditionChains struct {
	conditions []WordResultCondition
}

// NewWordResultConditionChains 创建链式条件（AND 语义）。
func NewWordResultConditionChains(conditions ...WordResultCondition) WordResultCondition {
	if len(conditions) == 0 {
		return WordResultConditionAlwaysTrue
	}
	cp := make([]WordResultCondition, len(conditions))
	copy(cp, conditions)
	return &wordResultConditionChains{conditions: cp}
}

func (c *wordResultConditionChains) Match(wordResult *WordResult, text string, modeEnum WordValidMode, context *WordContext) bool {
	for _, cond := range c.conditions {
		if !cond.Match(wordResult, text, modeEnum, context) {
			return false
		}
	}
	return true
}

// ===== 工具类 =====

// WordResultConditions 工具类，对应 Java 的 WordResultConditions。
var WordResultConditions = wordResultConditionsHelper{}

type wordResultConditionsHelper struct{}

// Defaults 返回默认条件（英文全词匹配）。
func (wordResultConditionsHelper) Defaults() WordResultCondition {
	return WordResultConditions.EnglishWordMatch()
}

// AlwaysTrue 返回恒为真的条件。
func (wordResultConditionsHelper) AlwaysTrue() WordResultCondition {
	return WordResultConditionAlwaysTrue
}

// EnglishWordMatch 返回英文全词匹配条件。
func (wordResultConditionsHelper) EnglishWordMatch() WordResultCondition {
	return WordResultConditionEnglishWordMatch
}

// EnglishWordNumMatch 返回英文/数字全词匹配条件。
func (wordResultConditionsHelper) EnglishWordNumMatch() WordResultCondition {
	return WordResultConditionEnglishWordNumMatch
}

// WordTags 返回按标签过滤的条件。
func (wordResultConditionsHelper) WordTags(tags []string, wordTag WordTag) WordResultCondition {
	return NewWordResultConditionWordTagsMatch(tags, wordTag)
}

// Chains 创建链式条件（AND 语义）。
func (wordResultConditionsHelper) Chains(condition WordResultCondition, others ...WordResultCondition) WordResultCondition {
	all := make([]WordResultCondition, 0, 1+len(others))
	all = append(all, condition)
	all = append(all, others...)
	return NewWordResultConditionChains(all...)
}
