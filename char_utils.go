package sensitiveword

import (
	"strings"
	"unicode"
)

// toHalfWidth 将全角字符转换为半角字符，对应 InnerCharUtils.toHalfWidth。
func toHalfWidth(original rune) rune {
	// 全角空格
	if original == '\u3000' {
		return ' '
	}
	// 其他可转换的全角字符
	if original >= '\uFF01' && original <= '\uFF5E' {
		return original - 0xFEE0
	}
	return original
}

// parseInt 快速解析纯数字字符串，对应 InnerCharUtils.parseInt。
// 输入应保证仅包含 '0'-'9'，否则返回 0。
func parseInt(text string) int {
	sum := 0
	for _, c := range text {
		if c < '0' || c > '9' {
			return 0
		}
		sum = sum*10 + int(c-'0')
	}
	return sum
}

// isEnglish 判断字符是否为英文字母，对应 CharUtil.isEnglish。
func isEnglish(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isSpace 判断字符是否为空白字符，对应 CharUtil.isSpace。
// 与 Java Character.isWhitespace 行为类似。
func isSpace(c rune) bool {
	return unicode.IsSpace(c)
}

// isDigit 判断字符是否为数字 0-9，对应 CharUtil.isNumber。
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// isDigitOrLetter 判断字符是否为字母或数字，对应 CharUtil.isDigitOrLetter。
func isDigitOrLetter(c rune) bool {
	return isDigit(c) || isEnglish(c)
}

// isEmailChar 邮箱合法字符，对应 CharUtil.isEmilChar。
// 包含字母、数字、下划线、点、@、减号、加号。
func isEmailChar(c rune) bool {
	if isDigitOrLetter(c) {
		return true
	}
	switch c {
	case '_', '.', '@', '-', '+':
		return true
	}
	return false
}

// isWebSiteChar 网址合法字符，对应 CharUtil.isWebSiteChar。
// 包含字母、数字、点、减号、下划线。
func isWebSiteChar(c rune) bool {
	if isDigitOrLetter(c) {
		return true
	}
	switch c {
	case '.', '-', '_':
		return true
	}
	return false
}

// isEmail 校验字符串是否为合法邮箱，对应 RegexUtil.isEmail。
// 长度限制：[6, MaxEmailLen]
func isEmail(s string) bool {
	if len(s) < 6 || len(s) > MaxEmailLen {
		return false
	}
	atIdx := strings.Index(s, "@")
	if atIdx <= 0 || atIdx == len(s)-1 {
		return false
	}
	// 仅允许一个 @
	if strings.Index(s[atIdx+1:], "@") >= 0 {
		return false
	}
	// 校验本地部分
	for _, c := range s[:atIdx] {
		if !isEmailChar(c) {
			return false
		}
	}
	// 校验域名部分
	domain := s[atIdx+1:]
	if !isWebSite(domain) {
		return false
	}
	return true
}

// isUrl 校验是否为 http(s):// 开头的 URL，对应 RegexUtil.isUrl。
func isUrl(s string) bool {
	if len(s) < 4 || len(s) > MaxWebSiteLen {
		return false
	}
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "http://") {
		return isValidHost(s[len("http://"):])
	}
	if strings.HasPrefix(lower, "https://") {
		return isValidHost(s[len("https://"):])
	}
	return false
}

// isWebSite 校验是否为合法的网址(无需 http 前缀)，对应 RegexUtil.isWebSite。
// 例如：www.baidu.com、baidu.com
func isWebSite(s string) bool {
	if len(s) < 4 || len(s) > MaxWebSiteLen {
		return false
	}
	return isValidHost(s)
}

// isValidHost 校验 host 部分是否合法：
// 至少包含一个点号；字符只允许字母、数字、点、减号、下划线、冒号(端口)。
func isValidHost(s string) bool {
	if s == "" {
		return false
	}
	hasDot := false
	for _, c := range s {
		if c == '.' {
			hasDot = true
			continue
		}
		if c == ':' {
			// 端口分隔符
			continue
		}
		if !isWebSiteChar(c) {
			return false
		}
	}
	return hasDot
}

// getMappingChar 从格式化映射表获取字符，对应 InnerWordFormatUtils.getMappingChar。
// 不存在映射时返回原始字符。
func getMappingChar(m map[rune]rune, c rune) rune {
	if mc, ok := m[c]; ok {
		return mc
	}
	return c
}

// getStringByRange 截取字符串，对应 InnerWordCharUtils.getString(text, start, end)。
// 基于 rune 截取，与 Java substring 基于 char 一致。
func getStringByRange(text string, startIndex, endIndex int) string {
	if startIndex < 0 {
		startIndex = 0
	}
	runes := []rune(text)
	if endIndex > len(runes) {
		endIndex = len(runes)
	}
	if startIndex >= endIndex {
		return ""
	}
	return string(runes[startIndex:endIndex])
}

// runesOf 将字符串转为 rune 切片，便于按下标访问。
func runesOf(s string) []rune {
	return []rune(s)
}
