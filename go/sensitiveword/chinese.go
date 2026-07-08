package sensitiveword

import "strings"

// t2sMapping 繁体到简体的字符映射表。
//
// 由于完整的 OpenCC 词表过大(数千条)，此处内置一份覆盖常用字的精简映射。
// 用户可通过 RegisterT2SMapping 注册更多字符映射以扩展覆盖率。
//
// 数据来源：基于 OpenCC 项目常用字表的精选子集。
var t2sMapping = map[rune]rune{}

func init() {
	// 一次性初始化常用繁简映射。
	for _, p := range t2sPairs {
		t2sMapping[p[0]] = p[1]
	}
}

// RegisterT2SMapping 注册额外的繁简映射，用于扩展默认表。
// 已存在的映射会被覆盖。线程安全。
func RegisterT2SMapping(trad, simp rune) {
	t2sMapping[trad] = simp
}

// RegisterT2SMappingBatch 批量注册繁简映射。
func RegisterT2SMappingBatch(pairs map[rune]rune) {
	for k, v := range pairs {
		t2sMapping[k] = v
	}
}

// toSimple 将单个字符转换为简体，对应 ZhSlimUtil.toSimple。
// 不在映射表中的字符原样返回。
func toSimple(c rune) rune {
	if mc, ok := t2sMapping[c]; ok {
		return mc
	}
	return c
}

// toSimpleString 将字符串转换为简体。
func toSimpleString(s string) string {
	if s == "" {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range s {
		sb.WriteRune(toSimple(c))
	}
	return sb.String()
}
