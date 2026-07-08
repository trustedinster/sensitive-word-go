package sensitiveword

// SensitiveWordBs 敏感词引导类，对应 Java 的 SensitiveWordBs。
//
// 采用 fluent-api 风格：通过链式 Setter 配置各项参数，最后调用 Init() 完成初始化。
//
// 使用示例:
//
//	bs := NewSensitiveWordBs().
//	    IgnoreCase(true).
//	    EnableNumCheck(true).
//	    Init()
//	contains := bs.Contains("some text")
//	words := bs.FindAll("some text")
//	replaced := bs.Replace("some text")
type SensitiveWordBs struct {
	// 格式统一化
	ignoreCase         bool
	ignoreWidth        bool
	ignoreNumStyle     bool
	ignoreChineseStyle bool
	ignoreEnglishStyle bool
	ignoreRepeat       bool
	wordFailFast       bool

	// 开启校验
	enableNumCheck   bool
	enableEmailCheck bool
	enableUrlCheck   bool
	enableWordCheck  bool
	enableIpv4Check  bool

	// 额外配置
	numCheckLen int

	// 数据与策略
	wordData             *WordData
	wordDataAllow        *WordData
	wordDeny             WordDeny
	wordAllow            WordAllow
	wordReplace           WordReplace
	wordTag              WordTag
	charIgnore           SensitiveWordCharIgnore
	wordResultCondition  WordResultCondition
	wordFormatText       WordFormatText

	// 单项检测策略
	wordCheckWord  WordCheck
	wordCheckNum   WordCheck
	wordCheckEmail WordCheck
	wordCheckUrl   WordCheck
	wordCheckIpv4  WordCheck

	// 组合器
	wordCheckCombine    *WordCheckCombine
	wordFormatCombine   *WordFormatCombine
	wordAllowDenyCombine *WordAllowDenyCombine

	// 引擎
	sensitiveWord *SensitiveWord

	// 初始化后的上下文
	context *WordContext
}

// NewSensitiveWordBs 创建引导类实例（未初始化）。
func NewSensitiveWordBs() *SensitiveWordBs {
	return &SensitiveWordBs{
		// 格式统一化默认值
		ignoreCase:         true,
		ignoreWidth:        true,
		ignoreNumStyle:     true,
		ignoreChineseStyle: true,
		ignoreEnglishStyle: true,
		ignoreRepeat:       false,
		wordFailFast:       true,

		// 校验默认值
		enableWordCheck: true,

		// 额外配置
		numCheckLen: DefaultNumCheckLen,

		// 数据与策略默认值
		wordData:             NewWordData(),
		wordDataAllow:        NewWordData(),
		wordDeny:             WordDenys.Defaults(),
		wordAllow:            WordAllows.Defaults(),
		wordReplace:           WordReplaces.Defaults(),
		wordTag:              WordTags.Defaults(),
		charIgnore:           SensitiveWordCharIgnores.Defaults(),
		wordResultCondition:  WordResultConditions.Defaults(),
		wordFormatText:       WordFormatTexts.Defaults(),

		// 单项检测策略
		wordCheckWord:  WordChecks.Word(),
		wordCheckNum:   WordChecks.Num(),
		wordCheckEmail: WordChecks.Email(),
		wordCheckUrl:   WordChecks.Url(),
		wordCheckIpv4:  WordChecks.Ipv4(),

		// 组合器
		wordCheckCombine:     NewWordCheckCombine(),
		wordFormatCombine:    NewWordFormatCombine(),
		wordAllowDenyCombine: NewWordAllowDenyCombine(),

		// 引擎
		sensitiveWord: NewSensitiveWord(),
	}
}

// Init 执行初始化：构建上下文、格式化策略、检测策略、DFA 树。
// 必须在配置完成后调用一次。
func (b *SensitiveWordBs) Init() *SensitiveWordBs {
	// 1. 初始化 context
	context := b.initContext()

	// 2. 格式化策略
	charFormat := b.wordFormatCombine.InitWordFormat(context)
	context.SetWordFormat(charFormat)

	// 3. 检测策略
	sensitiveCheck := b.wordCheckCombine.InitWordCheck(context)
	context.SetWordCheck(sensitiveCheck)

	// 4. 初始化黑/白名单 DFA
	wordAllowList := b.wordAllow.Allow()
	wordDenyList := b.wordDeny.Deny()
	denyList := b.wordAllowDenyCombine.GetActualDenyList(wordAllowList, wordDenyList, context)
	b.wordData.InitWordData(denyList)

	// 白名单也需要格式化
	actualAllowList := formatWordList(wordAllowList, context)
	b.wordDataAllow.InitWordData(actualAllowList)

	// 5. 更新 context
	context.SetWordData(b.wordData)
	context.SetWordDataAllow(b.wordDataAllow)

	b.context = context
	return b
}

// initContext 构建默认上下文并填充配置。
func (b *SensitiveWordBs) initContext() *WordContext {
	context := NewWordContext()

	// 格式统一化
	context.SetIgnoreCase(b.ignoreCase).
		SetIgnoreWidth(b.ignoreWidth).
		SetIgnoreNumStyle(b.ignoreNumStyle).
		SetIgnoreChineseStyle(b.ignoreChineseStyle).
		SetIgnoreEnglishStyle(b.ignoreEnglishStyle).
		SetIgnoreRepeat(b.ignoreRepeat).
		SetWordFailFast(b.wordFailFast).
		SetWordFormatText(b.wordFormatText)

	// 开启校验
	context.SetEnableNumCheck(b.enableNumCheck).
		SetEnableEmailCheck(b.enableEmailCheck).
		SetEnableUrlCheck(b.enableUrlCheck).
		SetEnableWordCheck(b.enableWordCheck).
		SetEnableIpv4Check(b.enableIpv4Check)

	// 校验策略实现
	context.SetWordCheckWord(b.wordCheckWord).
		SetWordCheckEmail(b.wordCheckEmail).
		SetWordCheckNum(b.wordCheckNum).
		SetWordCheckUrl(b.wordCheckUrl).
		SetWordCheckIpv4(b.wordCheckIpv4)

	// 额外配置
	context.SetNumCheckLen(b.numCheckLen).
		SetWordReplace(b.wordReplace).
		SetWordData(b.wordData).
		SetWordTag(b.wordTag).
		SetCharIgnore(b.charIgnore).
		SetWordResultCondition(b.wordResultCondition)

	return context
}

// Context 返回初始化后的上下文（Init 后可用）。
func (b *SensitiveWordBs) Context() *WordContext { return b.context }

// ===== 格式统一化 Setter =====

func (b *SensitiveWordBs) SetIgnoreCase(v bool) *SensitiveWordBs                 { b.ignoreCase = v; return b }
func (b *SensitiveWordBs) SetIgnoreWidth(v bool) *SensitiveWordBs                { b.ignoreWidth = v; return b }
func (b *SensitiveWordBs) SetIgnoreNumStyle(v bool) *SensitiveWordBs             { b.ignoreNumStyle = v; return b }
func (b *SensitiveWordBs) SetIgnoreChineseStyle(v bool) *SensitiveWordBs         { b.ignoreChineseStyle = v; return b }
func (b *SensitiveWordBs) SetIgnoreEnglishStyle(v bool) *SensitiveWordBs         { b.ignoreEnglishStyle = v; return b }
func (b *SensitiveWordBs) SetIgnoreRepeat(v bool) *SensitiveWordBs               { b.ignoreRepeat = v; return b }
func (b *SensitiveWordBs) SetWordFailFast(v bool) *SensitiveWordBs               { b.wordFailFast = v; return b }

// ===== 开启校验 Setter =====

func (b *SensitiveWordBs) SetEnableNumCheck(v bool) *SensitiveWordBs             { b.enableNumCheck = v; return b }
func (b *SensitiveWordBs) SetEnableEmailCheck(v bool) *SensitiveWordBs           { b.enableEmailCheck = v; return b }
func (b *SensitiveWordBs) SetEnableUrlCheck(v bool) *SensitiveWordBs             { b.enableUrlCheck = v; return b }
func (b *SensitiveWordBs) SetEnableWordCheck(v bool) *SensitiveWordBs            { b.enableWordCheck = v; return b }
func (b *SensitiveWordBs) SetEnableIpv4Check(v bool) *SensitiveWordBs            { b.enableIpv4Check = v; return b }

// ===== 额外配置 Setter =====

func (b *SensitiveWordBs) SetNumCheckLen(v int) *SensitiveWordBs                 { b.numCheckLen = v; return b }

// ===== 数据与策略 Setter =====

func (b *SensitiveWordBs) SetWordData(v *WordData) *SensitiveWordBs              { b.wordData = v; return b }
func (b *SensitiveWordBs) SetWordDataAllow(v *WordData) *SensitiveWordBs         { b.wordDataAllow = v; return b }
func (b *SensitiveWordBs) SetWordDeny(v WordDeny) *SensitiveWordBs               { b.wordDeny = v; return b }
func (b *SensitiveWordBs) SetWordAllow(v WordAllow) *SensitiveWordBs             { b.wordAllow = v; return b }
func (b *SensitiveWordBs) SetWordReplace(v WordReplace) *SensitiveWordBs         { b.wordReplace = v; return b }
func (b *SensitiveWordBs) SetWordTag(v WordTag) *SensitiveWordBs                 { b.wordTag = v; return b }
func (b *SensitiveWordBs) SetCharIgnore(v SensitiveWordCharIgnore) *SensitiveWordBs { b.charIgnore = v; return b }
func (b *SensitiveWordBs) SetWordResultCondition(v WordResultCondition) *SensitiveWordBs { b.wordResultCondition = v; return b }
func (b *SensitiveWordBs) SetWordFormatText(v WordFormatText) *SensitiveWordBs   { b.wordFormatText = v; return b }

// ===== 单项检测策略 Setter =====

func (b *SensitiveWordBs) SetWordCheckWord(v WordCheck) *SensitiveWordBs         { b.wordCheckWord = v; return b }
func (b *SensitiveWordBs) SetWordCheckNum(v WordCheck) *SensitiveWordBs          { b.wordCheckNum = v; return b }
func (b *SensitiveWordBs) SetWordCheckEmail(v WordCheck) *SensitiveWordBs        { b.wordCheckEmail = v; return b }
func (b *SensitiveWordBs) SetWordCheckUrl(v WordCheck) *SensitiveWordBs          { b.wordCheckUrl = v; return b }
func (b *SensitiveWordBs) SetWordCheckIpv4(v WordCheck) *SensitiveWordBs         { b.wordCheckIpv4 = v; return b }

// ===== 组合器 Setter =====

func (b *SensitiveWordBs) SetWordCheckCombine(v *WordCheckCombine) *SensitiveWordBs         { b.wordCheckCombine = v; return b }
func (b *SensitiveWordBs) SetWordFormatCombine(v *WordFormatCombine) *SensitiveWordBs       { b.wordFormatCombine = v; return b }
func (b *SensitiveWordBs) SetWordAllowDenyCombine(v *WordAllowDenyCombine) *SensitiveWordBs { b.wordAllowDenyCombine = v; return b }

// ===== 公开方法 =====

// Contains 是否包含敏感词。
func (b *SensitiveWordBs) Contains(target string) bool {
	return b.sensitiveWord.Contains(target, b.context)
}

// FindAll 返回所有敏感词（字符串形式，已去重保序）。
func (b *SensitiveWordBs) FindAll(target string) []string {
	results := FindAllWithHandler(b.sensitiveWord, target, b.context, WordResultHandlers.Word())
	if len(results) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(results))
	deduped := make([]string, 0, len(results))
	for _, w := range results {
		if !seen[w] {
			seen[w] = true
			deduped = append(deduped, w)
		}
	}
	return deduped
}

// FindFirst 返回第一个敏感词（字符串形式），未命中返回空字符串。
func (b *SensitiveWordBs) FindFirst(target string) string {
	return FindFirstWithHandler(b.sensitiveWord, target, b.context, WordResultHandlers.Word())
}

// FindAllRaw 返回所有敏感词的原始 WordResult 列表。
func (b *SensitiveWordBs) FindAllRaw(target string) []*WordResult {
	return b.sensitiveWord.FindAll(target, b.context)
}

// FindFirstRaw 返回第一个敏感词的原始 WordResult，未命中返回 nil。
func (b *SensitiveWordBs) FindFirstRaw(target string) *WordResult {
	return b.sensitiveWord.FindFirst(target, b.context)
}

// FindAllWordTags 返回所有敏感词及其标签。
func (b *SensitiveWordBs) FindAllWordTags(target string) []*WordTagsDto {
	return FindAllWithHandler(b.sensitiveWord, target, b.context, WordResultHandlers.WordTags())
}

// FindFirstWordTags 返回第一个敏感词及其标签。
func (b *SensitiveWordBs) FindFirstWordTags(target string) *WordTagsDto {
	return FindFirstWithHandler(b.sensitiveWord, target, b.context, WordResultHandlers.WordTags())
}

// Replace 替换所有敏感词为默认字符。
func (b *SensitiveWordBs) Replace(target string) string {
	return b.sensitiveWord.Replace(target, b.context)
}

// Tags 获取敏感词的标签。
func (b *SensitiveWordBs) Tags(word string) []string {
	return getWordTags(word, b.context)
}

// ===== 动态增删 =====

// AddWord 增量新增敏感词（黑名单）。
// 新增的词会经过与初始化一致的格式化处理。
func (b *SensitiveWordBs) AddWord(words ...string) {
	if len(words) == 0 {
		return
	}
	formatList := formatWordList(words, b.context)
	b.wordData.AddWord(formatList)
}

// RemoveWord 增量删除敏感词（黑名单）。
func (b *SensitiveWordBs) RemoveWord(words ...string) {
	if len(words) == 0 {
		return
	}
	formatList := formatWordList(words, b.context)
	b.wordData.RemoveWord(formatList)
}

// AddWordAllow 增量新增白名单。
func (b *SensitiveWordBs) AddWordAllow(words ...string) {
	if len(words) == 0 {
		return
	}
	formatList := formatWordList(words, b.context)
	b.wordDataAllow.AddWord(formatList)
}

// RemoveWordAllow 增量删除白名单。
func (b *SensitiveWordBs) RemoveWordAllow(words ...string) {
	if len(words) == 0 {
		return
	}
	formatList := formatWordList(words, b.context)
	b.wordDataAllow.RemoveWord(formatList)
}

// Destroy 释放 DFA 树内存。
func (b *SensitiveWordBs) Destroy() {
	b.wordData.Destroy()
	b.wordDataAllow.Destroy()
}
