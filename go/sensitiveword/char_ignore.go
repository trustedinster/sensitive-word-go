package sensitiveword

// SensitiveWordCharIgnore 字符忽略策略接口，对应 Java 的 ISensitiveWordCharIgnore。
//
// 在敏感词匹配过程中，若某字符被判定为忽略，则跳过该字符继续匹配，
// 典型场景：忽略标点、空格等干扰字符。
type SensitiveWordCharIgnore interface {
	// Ignore 判断 text 中下标 ix 处的字符是否应被忽略。
	Ignore(ix int, text string, innerCtx *InnerSensitiveWordContext) bool
}

// ===== 实现 =====

// noneSensitiveWordCharIgnore 不忽略任何字符，对应 NoneSensitiveWordCharIgnore。
type noneSensitiveWordCharIgnore struct{}

// NoneSensitiveWordCharIgnore 单例。
var NoneSensitiveWordCharIgnore = &noneSensitiveWordCharIgnore{}

func (i *noneSensitiveWordCharIgnore) Ignore(ix int, text string, innerCtx *InnerSensitiveWordContext) bool {
	return false
}

// specialCharSensitiveWordCharIgnore 忽略常见特殊字符，对应 SpecialCharSensitiveWordCharIgnore。
type specialCharSensitiveWordCharIgnore struct{}

// SpecialCharSensitiveWordCharIgnore 单例。
var SpecialCharSensitiveWordCharIgnore = &specialCharSensitiveWordCharIgnore{}

// specialCharSet 需要忽略的特殊字符集合。
var specialCharSet = func() map[rune]bool {
	const special = "`-=~!@#$%^&*()_+[]{}\\|;:'\",./<>?"
	m := make(map[rune]bool, len(special))
	for _, c := range special {
		m[c] = true
	}
	return m
}()

func (i *specialCharSensitiveWordCharIgnore) Ignore(ix int, text string, innerCtx *InnerSensitiveWordContext) bool {
	runes := []rune(text)
	if ix < 0 || ix >= len(runes) {
		return false
	}
	return specialCharSet[runes[ix]]
}

// ===== 工具类 =====

// SensitiveWordCharIgnores 工具类，对应 Java 的 SensitiveWordCharIgnores。
var SensitiveWordCharIgnores = sensitiveWordCharIgnoresHelper{}

type sensitiveWordCharIgnoresHelper struct{}

// SpecialChars 返回忽略特殊字符的策略。
func (sensitiveWordCharIgnoresHelper) SpecialChars() SensitiveWordCharIgnore {
	return SpecialCharSensitiveWordCharIgnore
}

// None 返回不忽略任何字符的策略。
func (sensitiveWordCharIgnoresHelper) None() SensitiveWordCharIgnore {
	return NoneSensitiveWordCharIgnore
}

// Defaults 返回默认策略（不忽略）。
func (sensitiveWordCharIgnoresHelper) Defaults() SensitiveWordCharIgnore {
	return NoneSensitiveWordCharIgnore
}
