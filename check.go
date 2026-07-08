package sensitiveword

import (
	"strings"
	"unicode"
)

// WordCheck 敏感信息检测策略接口，对应 Java 的 IWordCheck。
//
// 实现方负责从指定起始下标开始检测，返回命中的黑白名单长度信息。
// 可使用责任链模式组合多个策略（WordCheckArray）。
//
// 实现应为线程安全（无状态）。
type WordCheck interface {
	// SensitiveCheck 从 beginIndex 开始检测敏感信息。
	// 返回检测结果，包含命中长度、类别等信息。
	SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult
}

// ===== 辅助：模板方法 =====

// sensitiveCheckCommon 实现 AbstractWordCheck.sensitiveCheck 的模板逻辑：
// 空文本返回零长度结果；否则调用 actualLengthFn 获取实际长度。
func sensitiveCheckCommon(beginIndex int, ctx *InnerSensitiveWordContext, checkType WordType,
	actualLengthFn func(int, *InnerSensitiveWordContext) *WordLengthResult) *WordCheckResult {
	result := &WordCheckResult{
		wordLengthResult: NewWordLengthResult(),
		checkType:        checkType,
	}
	if ctx.OriginalText() == "" {
		return result
	}
	result.wordLengthResult = actualLengthFn(beginIndex, ctx)
	return result
}

// conditionActualLength 实现 AbstractConditionWordCheck.getActualLength 的通用逻辑：
// 从 beginIndex 遍历文本，跳过忽略字符，收集满足 isCharCondition 的字符，
// 最终由 isStringCondition 判定是否命中。
func conditionActualLength(
	beginIndex int, ctx *InnerSensitiveWordContext,
	isCharCondition func(mappingChar rune, index int, ctx *InnerSensitiveWordContext) bool,
	isStringCondition func(index int, sb []rune, ctx *InnerSensitiveWordContext) bool,
) *WordLengthResult {
	charIgnore := ctx.WordContext().CharIgnore()
	txt := ctx.OriginalText()
	runes := []rune(txt)
	formatCharMapping := ctx.FormatCharMapping()

	var sb []rune
	actualLength := 0
	tempIgnoreLen := 0
	currentIx := beginIndex

	for i := beginIndex; i < len(runes); i++ {
		currentIx = i
		if charIgnore.Ignore(i, txt, ctx) {
			tempIgnoreLen++
			continue
		}
		mappingChar := getMappingChar(formatCharMapping, runes[i])
		if isCharCondition(mappingChar, i, ctx) {
			sb = append(sb, runes[i])
		} else {
			break
		}
	}

	if isStringCondition(currentIx, sb, ctx) {
		actualLength = len(sb) + tempIgnoreLen
	}

	return NewWordLengthResult().
		SetWordDenyLen(actualLength).
		SetWordAllowLen(0)
}

// ===== WordCheckNone =====

// wordCheckNone 未匹配实现，对应 Java 的 WordCheckNone。
type wordCheckNone struct{}

// WordCheckNone 单例。
var WordCheckNone = &wordCheckNone{}

// noneCheckResult 预构建的空结果。
var noneCheckResult = &WordCheckResult{
	wordLengthResult: NewWordLengthResult(),
	checkType:        WordTypeDefaults,
}

// GetNoneResult 返回空结果，对应 WordCheckNone.getNoneResult。
func GetNoneResult() *WordCheckResult { return noneCheckResult }

func (c *wordCheckNone) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return noneCheckResult
}

// ===== WordCheckWord 敏感词 DFA 检测 =====

// wordCheckWord 基于黑/白名单 DFA 树的敏感词检测，对应 Java 的 WordCheckWord。
type wordCheckWord struct{}

// WordCheckWord 单例。
var WordCheckWord = &wordCheckWord{}

func (c *wordCheckWord) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return sensitiveCheckCommon(beginIndex, ctx, WordTypeWord, c.getActualLength)
}

func (c *wordCheckWord) getActualLength(beginIndex int, ctx *InnerSensitiveWordContext) *WordLengthResult {
	runes := []rune(ctx.OriginalText())
	formatCharMapping := ctx.FormatCharMapping()
	context := ctx.WordContext()
	wordData := context.WordData()
	wordDataAllow := context.WordDataAllow()
	wordCharIgnore := context.CharIgnore()
	failFast := context.WordFailFast()

	var stringBuilder []rune
	tempLen := 0
	maxWhite := 0
	maxBlack := 0
	skipLen := 0

	for i := beginIndex; i < len(runes); i++ {
		if wordCharIgnore.Ignore(i, ctx.OriginalText(), ctx) && tempLen != 0 {
			tempLen++
			skipLen++
			continue
		}
		mappingChar := getMappingChar(formatCharMapping, runes[i])
		stringBuilder = append(stringBuilder, mappingChar)
		tempLen++

		containsAllow := wordDataAllow.Contains(stringBuilder, ctx)
		containsDeny := wordData.Contains(stringBuilder, ctx)

		if containsAllow == WordContainsEnd {
			maxWhite = tempLen
			if failFast {
				containsAllow = WordContainsNotFound
			}
		}
		if containsDeny == WordContainsEnd {
			maxBlack = tempLen
			if failFast {
				containsDeny = WordContainsNotFound
			}
		}
		if containsAllow == WordContainsNotFound && containsDeny == WordContainsNotFound {
			break
		}
	}

	wordAllow := ""
	wordDeny := ""
	if len(stringBuilder) > 0 {
		whiteEnd := maxWhite - skipLen
		if whiteEnd < 0 {
			whiteEnd = 0
		}
		denyEnd := maxBlack - skipLen
		if denyEnd < 0 {
			denyEnd = 0
		}
		wordAllow = string(stringBuilder[:whiteEnd])
		wordDeny = string(stringBuilder[:denyEnd])
	}

	return NewWordLengthResult().
		SetWordAllowLen(maxWhite).
		SetWordDenyLen(maxBlack).
		SetWordAllow(wordAllow).
		SetWordDeny(wordDeny)
}

// ===== WordCheckNum 数字检测 =====

// wordCheckNum 连续数字检测，对应 Java 的 WordCheckNum。
type wordCheckNum struct{}

// WordCheckNum 单例。
var WordCheckNum = &wordCheckNum{}

func (c *wordCheckNum) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return sensitiveCheckCommon(beginIndex, ctx, WordTypeNum, c.getActualLength)
}

func (c *wordCheckNum) getActualLength(beginIndex int, ctx *InnerSensitiveWordContext) *WordLengthResult {
	return conditionActualLength(beginIndex, ctx,
		func(mc rune, _ int, _ *InnerSensitiveWordContext) bool {
			return unicode.IsDigit(mc)
		},
		func(_ int, sb []rune, ic *InnerSensitiveWordContext) bool {
			return len(sb) >= ic.WordContext().NumCheckLen()
		},
	)
}

// ===== WordCheckEmail 邮箱检测 =====

// wordCheckEmail 邮箱检测，对应 Java 的 WordCheckEmail。
type wordCheckEmail struct{}

// WordCheckEmail 单例。
var WordCheckEmail = &wordCheckEmail{}

func (c *wordCheckEmail) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return sensitiveCheckCommon(beginIndex, ctx, WordTypeEmail, c.getActualLength)
}

func (c *wordCheckEmail) getActualLength(beginIndex int, ctx *InnerSensitiveWordContext) *WordLengthResult {
	return conditionActualLength(beginIndex, ctx,
		func(mc rune, _ int, _ *InnerSensitiveWordContext) bool {
			return isEmailChar(mc)
		},
		func(_ int, sb []rune, _ *InnerSensitiveWordContext) bool {
			bufferLen := len(sb)
			if bufferLen < 6 || bufferLen > MaxEmailLen {
				return false
			}
			return isEmail(string(sb))
		},
	)
}

// ===== WordCheckUrl URL 检测 =====

// wordCheckUrl URL 检测（需 http(s):// 前缀），对应 Java 的 WordCheckUrl。
type wordCheckUrl struct{}

// WordCheckUrl 单例。
var WordCheckUrl = &wordCheckUrl{}

func (c *wordCheckUrl) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return sensitiveCheckCommon(beginIndex, ctx, WordTypeURL, c.getActualLength)
}

func (c *wordCheckUrl) getActualLength(beginIndex int, ctx *InnerSensitiveWordContext) *WordLengthResult {
	return conditionActualLength(beginIndex, ctx,
		func(mc rune, _ int, _ *InnerSensitiveWordContext) bool {
			return isWebSiteChar(mc) || mc == ':' || mc == '/'
		},
		func(_ int, sb []rune, _ *InnerSensitiveWordContext) bool {
			bufferLen := len(sb)
			if bufferLen < 4 || bufferLen > MaxWebSiteLen {
				return false
			}
			return c.isUrl(string(sb))
		},
	)
}

// isUrl 判断是否为合法 URL，可被子策略覆盖。
func (c *wordCheckUrl) isUrl(text string) bool {
	return isUrl(text)
}

// ===== WordCheckUrlNoPrefix 无前缀 URL 检测 =====

// wordCheckUrlNoPrefix 无需 http 前缀的网址检测，对应 Java 的 WordCheckUrlNoPrefix。
type wordCheckUrlNoPrefix struct {
	wordCheckUrl
}

// WordCheckUrlNoPrefix 单例。
var WordCheckUrlNoPrefix = &wordCheckUrlNoPrefix{}

func (c *wordCheckUrlNoPrefix) isUrl(text string) bool {
	return isWebSite(text)
}

// ===== WordCheckIPV4 IPv4 检测 =====

// wordCheckIPV4 IPv4 地址检测，对应 Java 的 WordCheckIPV4。
type wordCheckIPV4 struct{}

// WordCheckIPV4 单例。
var WordCheckIPV4 = &wordCheckIPV4{}

func (c *wordCheckIPV4) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	return sensitiveCheckCommon(beginIndex, ctx, WordTypeIPV4, c.getActualLength)
}

func (c *wordCheckIPV4) getActualLength(beginIndex int, ctx *InnerSensitiveWordContext) *WordLengthResult {
	return conditionActualLength(beginIndex, ctx,
		func(mc rune, _ int, _ *InnerSensitiveWordContext) bool {
			return isDigit(mc) || mc == '.'
		},
		func(_ int, sb []rune, _ *InnerSensitiveWordContext) bool {
			bufferLen := len(sb)
			// 0.0.0.0 ~ 255.255.255.255
			if bufferLen < 7 || bufferLen > 15 {
				return false
			}
			parts := strings.Split(string(sb), ".")
			if len(parts) != 4 {
				return false
			}
			for _, numStr := range parts {
				num := parseInt(numStr)
				if num < 0 || num > 256 {
					return false
				}
			}
			return true
		},
	)
}

// ===== WordCheckArray 责任链 =====

// wordCheckArray 责任链组合，对应 Java 的 WordCheckArray。
// 依次调用各策略，首个命中（长度>0）即返回。
type wordCheckArray struct {
	checks []WordCheck
}

// NewWordCheckArray 创建责任链组合。
func NewWordCheckArray(checks []WordCheck) WordCheck {
	if len(checks) == 0 {
		return WordCheckNone
	}
	cp := make([]WordCheck, len(checks))
	copy(cp, checks)
	return &wordCheckArray{checks: cp}
}

func (c *wordCheckArray) SensitiveCheck(beginIndex int, ctx *InnerSensitiveWordContext) *WordCheckResult {
	for _, check := range c.checks {
		result := check.SensitiveCheck(beginIndex, ctx)
		wlr := result.WordLengthResult()
		if wlr.WordAllowLen() > 0 || wlr.WordDenyLen() > 0 {
			return result
		}
	}
	return GetNoneResult()
}

// ===== 工具类 WordChecks =====

// WordChecks 工具类，对应 Java 的 WordChecks。
var WordChecks = wordChecksHelper{}

type wordChecksHelper struct{}

func (wordChecksHelper) Word() WordCheck       { return WordCheckWord }
func (wordChecksHelper) Num() WordCheck        { return WordCheckNum }
func (wordChecksHelper) Email() WordCheck      { return WordCheckEmail }
func (wordChecksHelper) Url() WordCheck        { return WordCheckUrl }
func (wordChecksHelper) UrlNoPrefix() WordCheck { return WordCheckUrlNoPrefix }
func (wordChecksHelper) Ipv4() WordCheck       { return WordCheckIPV4 }
func (wordChecksHelper) None() WordCheck       { return WordCheckNone }

// Chains 组合多个检测策略为责任链。
func (h wordChecksHelper) Chains(checks ...WordCheck) WordCheck {
	if len(checks) == 0 {
		return WordCheckNone
	}
	return NewWordCheckArray(checks)
}

// ChainsList 基于切片组合。
func (h wordChecksHelper) ChainsList(checks []WordCheck) WordCheck {
	if len(checks) == 0 {
		return WordCheckNone
	}
	return NewWordCheckArray(checks)
}

// Array 基于切片创建责任链。
func (h wordChecksHelper) Array(checks []WordCheck) WordCheck {
	return h.ChainsList(checks)
}

// ===== 组合器 WordCheckCombine =====

// WordCheckCombine 根据上下文配置构建组合检测策略，对应 Java 的 WordCheckCombine。
//
// 组合顺序与 Java 保持一致：
//
//	word -> num -> email -> url -> ipv4
type WordCheckCombine struct{}

// NewWordCheckCombine 创建组合器。
func NewWordCheckCombine() *WordCheckCombine {
	return &WordCheckCombine{}
}

// InitWordCheck 依据上下文构造组合后的 WordCheck。
func (c *WordCheckCombine) InitWordCheck(context *WordContext) WordCheck {
	checks := make([]WordCheck, 0, 5)
	if context.EnableWordCheck() {
		if context.WordCheckWord() != nil {
			checks = append(checks, context.WordCheckWord())
		} else {
			checks = append(checks, WordChecks.Word())
		}
	}
	if context.EnableNumCheck() {
		if context.WordCheckNum() != nil {
			checks = append(checks, context.WordCheckNum())
		} else {
			checks = append(checks, WordChecks.Num())
		}
	}
	if context.EnableEmailCheck() {
		if context.WordCheckEmail() != nil {
			checks = append(checks, context.WordCheckEmail())
		} else {
			checks = append(checks, WordChecks.Email())
		}
	}
	if context.EnableUrlCheck() {
		if context.WordCheckUrl() != nil {
			checks = append(checks, context.WordCheckUrl())
		} else {
			checks = append(checks, WordChecks.Url())
		}
	}
	if context.EnableIpv4Check() {
		if context.WordCheckIpv4() != nil {
			checks = append(checks, context.WordCheckIpv4())
		} else {
			checks = append(checks, WordChecks.Ipv4())
		}
	}
	return WordChecks.ChainsList(checks)
}
