package sensitiveword

// WordConst 敏感词相关常量。
type WordConst struct{}

// 内置常量值。保留与 Java 项目一致的命名以便对照。
const (
	// wordConstIsEnd 旧版本用于 Map 中标识关键词结束的字段，本实现采用节点字段，仅保留常量作对照。
	wordConstIsEnd = "ED"
	// MaxWebSiteLen 最长网址长度。
	MaxWebSiteLen = 70
	// MaxEmailLen 最大邮箱长度。
	MaxEmailLen = 64
	// DefaultNumCheckLen 数字检测默认长度。
	DefaultNumCheckLen = 8
	// DefaultReplaceChar 默认替换字符。
	DefaultReplaceChar = '*'
)

// WordType 敏感词命中类别。
type WordType string

const (
	// WordTypeWord 敏感词。
	WordTypeWord WordType = "WORD"
	// WordTypeEmail 邮箱。
	WordTypeEmail WordType = "EMAIL"
	// WordTypeURL 链接。
	WordTypeURL WordType = "URL"
	// WordTypeNum 数字。
	WordTypeNum WordType = "NUM"
	// WordTypeIPV4 IPv4。
	WordTypeIPV4 WordType = "IPV4"
	// WordTypeDefaults 默认。
	WordTypeDefaults WordType = "DEFAULTS"
)

// WordValidMode 校验模式。
type WordValidMode int

const (
	// WordValidModeFailFast 快速失败，匹配到第一个即返回。
	WordValidModeFailFast WordValidMode = iota
	// WordValidModeFailOver 全量遍历，返回所有匹配。
	WordValidModeFailOver
)

// WordContainsType 包含类别，对应 WordContainsTypeEnum。
type WordContainsType int

const (
	// WordContainsNotFound 不存在。
	WordContainsNotFound WordContainsType = iota
	// WordContainsPrefix 命中前缀，但不是结尾。
	WordContainsPrefix
	// WordContainsEnd 命中结尾。
	WordContainsEnd
)
