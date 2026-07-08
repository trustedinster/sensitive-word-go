package sensitiveword

// WordContext 敏感词上下文，对应 Java 的 IWordContext / SensitiveWordContext。
//
// 该结构同时承担配置存储与运行时策略持有职责。
// 在 Init 阶段由引导类填充；在匹配阶段只读访问。
type WordContext struct {
	// ===== 格式统一化 =====
	ignoreCase         bool
	ignoreWidth        bool
	ignoreNumStyle     bool
	ignoreChineseStyle bool
	ignoreEnglishStyle bool
	ignoreRepeat       bool
	wordFailFast       bool

	// ===== 开启校验 =====
	enableNumCheck   bool
	enableEmailCheck bool
	enableUrlCheck   bool
	enableWordCheck  bool
	enableIpv4Check  bool

	// ===== 额外配置 =====
	numCheckLen int

	// ===== 策略 =====
	wordCheck           WordCheck           // 组合后的检测策略
	wordReplace         WordReplace         // 替换策略
	wordFormat          WordFormat          // 组合后的字符格式化策略
	wordFormatText      WordFormatText      // 文本整体格式化策略
	wordData            *WordData           // 黑名单 DFA
	wordDataAllow       *WordData           // 白名单 DFA
	wordTag             WordTag             // 标签策略
	charIgnore          SensitiveWordCharIgnore // 忽略字符策略
	wordResultCondition WordResultCondition // 结果匹配条件

	// 单项检测策略(用于组合)
	wordCheckWord  WordCheck
	wordCheckNum   WordCheck
	wordCheckEmail WordCheck
	wordCheckUrl   WordCheck
	wordCheckIpv4  WordCheck
}

// NewWordContext 创建一个具备默认值的上下文。
func NewWordContext() *WordContext {
	return &WordContext{
		ignoreCase:         true,
		ignoreWidth:        true,
		ignoreNumStyle:     true,
		ignoreChineseStyle: true,
		ignoreEnglishStyle: true,
		ignoreRepeat:       false,
		wordFailFast:       true,
		enableWordCheck:    true,
		numCheckLen:        DefaultNumCheckLen,
	}
}

// ===== Getter =====

func (c *WordContext) WordFailFast() bool             { return c.wordFailFast }
func (c *WordContext) IgnoreCase() bool               { return c.ignoreCase }
func (c *WordContext) IgnoreWidth() bool              { return c.ignoreWidth }
func (c *WordContext) IgnoreNumStyle() bool           { return c.ignoreNumStyle }
func (c *WordContext) IgnoreChineseStyle() bool       { return c.ignoreChineseStyle }
func (c *WordContext) IgnoreEnglishStyle() bool       { return c.ignoreEnglishStyle }
func (c *WordContext) IgnoreRepeat() bool             { return c.ignoreRepeat }
func (c *WordContext) EnableNumCheck() bool           { return c.enableNumCheck }
func (c *WordContext) EnableEmailCheck() bool         { return c.enableEmailCheck }
func (c *WordContext) EnableUrlCheck() bool           { return c.enableUrlCheck }
func (c *WordContext) EnableWordCheck() bool          { return c.enableWordCheck }
func (c *WordContext) EnableIpv4Check() bool          { return c.enableIpv4Check }
func (c *WordContext) NumCheckLen() int               { return c.numCheckLen }
func (c *WordContext) WordCheck() WordCheck           { return c.wordCheck }
func (c *WordContext) WordReplace() WordReplace       { return c.wordReplace }
func (c *WordContext) WordFormat() WordFormat         { return c.wordFormat }
func (c *WordContext) WordFormatText() WordFormatText { return c.wordFormatText }
func (c *WordContext) WordData() *WordData            { return c.wordData }
func (c *WordContext) WordDataAllow() *WordData       { return c.wordDataAllow }
func (c *WordContext) WordTag() WordTag               { return c.wordTag }
func (c *WordContext) CharIgnore() SensitiveWordCharIgnore { return c.charIgnore }
func (c *WordContext) WordResultCondition() WordResultCondition { return c.wordResultCondition }
func (c *WordContext) WordCheckWord() WordCheck       { return c.wordCheckWord }
func (c *WordContext) WordCheckNum() WordCheck        { return c.wordCheckNum }
func (c *WordContext) WordCheckEmail() WordCheck      { return c.wordCheckEmail }
func (c *WordContext) WordCheckUrl() WordCheck        { return c.wordCheckUrl }
func (c *WordContext) WordCheckIpv4() WordCheck       { return c.wordCheckIpv4 }

// ===== Setter (链式) =====

func (c *WordContext) SetWordFailFast(v bool) *WordContext             { c.wordFailFast = v; return c }
func (c *WordContext) SetIgnoreCase(v bool) *WordContext               { c.ignoreCase = v; return c }
func (c *WordContext) SetIgnoreWidth(v bool) *WordContext              { c.ignoreWidth = v; return c }
func (c *WordContext) SetIgnoreNumStyle(v bool) *WordContext           { c.ignoreNumStyle = v; return c }
func (c *WordContext) SetIgnoreChineseStyle(v bool) *WordContext       { c.ignoreChineseStyle = v; return c }
func (c *WordContext) SetIgnoreEnglishStyle(v bool) *WordContext       { c.ignoreEnglishStyle = v; return c }
func (c *WordContext) SetIgnoreRepeat(v bool) *WordContext             { c.ignoreRepeat = v; return c }
func (c *WordContext) SetEnableNumCheck(v bool) *WordContext           { c.enableNumCheck = v; return c }
func (c *WordContext) SetEnableEmailCheck(v bool) *WordContext         { c.enableEmailCheck = v; return c }
func (c *WordContext) SetEnableUrlCheck(v bool) *WordContext           { c.enableUrlCheck = v; return c }
func (c *WordContext) SetEnableWordCheck(v bool) *WordContext          { c.enableWordCheck = v; return c }
func (c *WordContext) SetEnableIpv4Check(v bool) *WordContext          { c.enableIpv4Check = v; return c }
func (c *WordContext) SetNumCheckLen(v int) *WordContext               { c.numCheckLen = v; return c }
func (c *WordContext) SetWordCheck(v WordCheck) *WordContext           { c.wordCheck = v; return c }
func (c *WordContext) SetWordReplace(v WordReplace) *WordContext       { c.wordReplace = v; return c }
func (c *WordContext) SetWordFormat(v WordFormat) *WordContext         { c.wordFormat = v; return c }
func (c *WordContext) SetWordFormatText(v WordFormatText) *WordContext { c.wordFormatText = v; return c }
func (c *WordContext) SetWordData(v *WordData) *WordContext            { c.wordData = v; return c }
func (c *WordContext) SetWordDataAllow(v *WordData) *WordContext       { c.wordDataAllow = v; return c }
func (c *WordContext) SetWordTag(v WordTag) *WordContext               { c.wordTag = v; return c }
func (c *WordContext) SetCharIgnore(v SensitiveWordCharIgnore) *WordContext { c.charIgnore = v; return c }
func (c *WordContext) SetWordResultCondition(v WordResultCondition) *WordContext {
	c.wordResultCondition = v
	return c
}
func (c *WordContext) SetWordCheckWord(v WordCheck) *WordContext  { c.wordCheckWord = v; return c }
func (c *WordContext) SetWordCheckNum(v WordCheck) *WordContext   { c.wordCheckNum = v; return c }
func (c *WordContext) SetWordCheckEmail(v WordCheck) *WordContext { c.wordCheckEmail = v; return c }
func (c *WordContext) SetWordCheckUrl(v WordCheck) *WordContext   { c.wordCheckUrl = v; return c }
func (c *WordContext) SetWordCheckIpv4(v WordCheck) *WordContext  { c.wordCheckIpv4 = v; return c }

// InnerSensitiveWordContext 内部运行时上下文，对应 Java 的 InnerSensitiveWordContext。
type InnerSensitiveWordContext struct {
	originalText      string
	formatCharMapping map[rune]rune
	modeEnum          WordValidMode
	wordContext       *WordContext
}

// NewInnerSensitiveWordContext 创建内部上下文。
func NewInnerSensitiveWordContext() *InnerSensitiveWordContext {
	return &InnerSensitiveWordContext{}
}

func (i *InnerSensitiveWordContext) OriginalText() string          { return i.originalText }
func (i *InnerSensitiveWordContext) FormatCharMapping() map[rune]rune { return i.formatCharMapping }
func (i *InnerSensitiveWordContext) ModeEnum() WordValidMode       { return i.modeEnum }
func (i *InnerSensitiveWordContext) WordContext() *WordContext     { return i.wordContext }

func (i *InnerSensitiveWordContext) SetOriginalText(s string) *InnerSensitiveWordContext {
	i.originalText = s
	return i
}
func (i *InnerSensitiveWordContext) SetFormatCharMapping(m map[rune]rune) *InnerSensitiveWordContext {
	i.formatCharMapping = m
	return i
}
func (i *InnerSensitiveWordContext) SetModeEnum(m WordValidMode) *InnerSensitiveWordContext {
	i.modeEnum = m
	return i
}
func (i *InnerSensitiveWordContext) SetWordContext(c *WordContext) *InnerSensitiveWordContext {
	i.wordContext = c
	return i
}
