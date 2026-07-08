package sensitiveword

import (
	"bufio"
	"embed"
	"io"
	"strings"
)

// WordAllow 白名单策略接口，对应 Java 的 IWordAllow。
// 返回的单词不会被当做敏感词。
type WordAllow interface {
	Allow() []string
}

// WordDeny 黑名单策略接口，对应 Java 的 IWordDeny。
// 返回的单词会被当做敏感词。
type WordDeny interface {
	Deny() []string
}

// ===== 空实现 =====

// wordAllowEmpty 空白名单，对应 WordAllowEmpty。
type wordAllowEmpty struct{}

// WordAllowEmpty 单例。
var WordAllowEmpty = &wordAllowEmpty{}

func (a *wordAllowEmpty) Allow() []string { return nil }

// wordDenyEmpty 空黑名单，对应 WordDenyEmpty。
type wordDenyEmpty struct{}

// WordDenyEmpty 单例。
var WordDenyEmpty = &wordDenyEmpty{}

func (d *wordDenyEmpty) Deny() []string { return nil }

// ===== 系统默认实现（从内嵌资源加载） =====

// dictFS 内嵌字典文件系统。空目录占位，实际文件在 data/ 下。
//
//go:embed data/*.txt
var dictFS embed.FS

// readDictLines 从内嵌资源读取所有非空行。
// path 相对于 data/ 目录，如 "sensitive_word_dict.txt"。
func readDictLines(path string) []string {
	fullPath := "data/" + path
	f, err := dictFS.Open(fullPath)
	if err != nil {
		return nil
	}
	defer f.Close()
	return readLinesFromReader(f)
}

// readLinesFromReader 逐行读取，跳过空行与空白行。
func readLinesFromReader(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

// wordDenySystem 系统默认黑名单，对应 WordDenySystem。
// 加载内嵌的 sensitive_word_dict.txt + sensitive_word_dict_en.txt + sensitive_word_deny.txt。
type wordDenySystem struct{}

// WordDenySystem 单例。
var WordDenySystem = &wordDenySystem{}

func (d *wordDenySystem) Deny() []string {
	result := readDictLines("sensitive_word_dict.txt")
	result = append(result, readDictLines("sensitive_word_dict_en.txt")...)
	result = append(result, readDictLines("sensitive_word_deny.txt")...)
	return result
}

// wordAllowSystem 系统默认白名单，对应 WordAllowSystem。
// 加载内嵌的 sensitive_word_allow.txt。
type wordAllowSystem struct{}

// WordAllowSystem 单例。
var WordAllowSystem = &wordAllowSystem{}

func (a *wordAllowSystem) Allow() []string {
	return readDictLines("sensitive_word_allow.txt")
}

// ===== 链式实现 =====

// wordDenyChain 责任链黑名单，对应 WordDenyInit 的链式实现。
type wordDenyChain struct {
	denies []WordDeny
}

// NewWordDenyChain 创建链式黑名单，依次合并各实现的 deny 结果。
func NewWordDenyChain(denies ...WordDeny) WordDeny {
	cp := make([]WordDeny, 0, len(denies))
	cp = append(cp, denies...)
	return &wordDenyChain{denies: cp}
}

func (d *wordDenyChain) Deny() []string {
	var result []string
	for _, deny := range d.denies {
		result = append(result, deny.Deny()...)
	}
	return result
}

// wordAllowChain 责任链白名单。
type wordAllowChain struct {
	allows []WordAllow
}

// NewWordAllowChain 创建链式白名单，依次合并各实现的 allow 结果。
func NewWordAllowChain(allows ...WordAllow) WordAllow {
	cp := make([]WordAllow, 0, len(allows))
	cp = append(cp, allows...)
	return &wordAllowChain{allows: cp}
}

func (a *wordAllowChain) Allow() []string {
	var result []string
	for _, allow := range a.allows {
		result = append(result, allow.Allow()...)
	}
	return result
}

// ===== 工具类 =====

// WordDenys 工具类，对应 Java 的 WordDenys。
var WordDenys = wordDenysHelper{}

type wordDenysHelper struct{}

// Defaults 返回系统默认黑名单。
func (wordDenysHelper) Defaults() WordDeny { return WordDenySystem }

// Empty 返回空黑名单。
func (wordDenysHelper) Empty() WordDeny { return WordDenyEmpty }

// Chains 创建链式黑名单。
func (wordDenysHelper) Chains(deny WordDeny, others ...WordDeny) WordDeny {
	all := make([]WordDeny, 0, 1+len(others))
	all = append(all, deny)
	all = append(all, others...)
	return NewWordDenyChain(all...)
}

// WordAllows 工具类，对应 Java 的 WordAllows。
var WordAllows = wordAllowsHelper{}

type wordAllowsHelper struct{}

// Defaults 返回系统默认白名单。
func (wordAllowsHelper) Defaults() WordAllow { return WordAllowSystem }

// Empty 返回空白名单。
func (wordAllowsHelper) Empty() WordAllow { return WordAllowEmpty }

// Chains 创建链式白名单。
func (wordAllowsHelper) Chains(allow WordAllow, others ...WordAllow) WordAllow {
	all := make([]WordAllow, 0, 1+len(others))
	all = append(all, allow)
	all = append(all, others...)
	return NewWordAllowChain(all...)
}

// ===== AllowDeny 组合器 =====

// WordAllowDenyCombine 黑白名单组合器，对应 Java 的 WordAllowDenyCombine。
//
// 先对黑白名单做格式化，再从黑名单中剔除白名单包含的词。
type WordAllowDenyCombine struct{}

// NewWordAllowDenyCombine 创建组合器。
func NewWordAllowDenyCombine() *WordAllowDenyCombine {
	return &WordAllowDenyCombine{}
}

// GetActualDenyList 获取最终的黑名单列表：
//  1. 格式化白名单与黑名单
//  2. 从黑名单中剔除白名单包含的词
func (c *WordAllowDenyCombine) GetActualDenyList(allowList []string, denyList []string, context *WordContext) []string {
	formatAllowList := formatWordList(allowList, context)
	formatDenyList := formatWordList(denyList, context)

	if len(formatDenyList) == 0 {
		return nil
	}
	if len(formatAllowList) == 0 {
		return formatDenyList
	}

	allowSet := make(map[string]bool, len(formatAllowList))
	for _, w := range formatAllowList {
		allowSet[w] = true
	}
	result := make([]string, 0, len(formatDenyList))
	seen := make(map[string]bool, len(formatDenyList))
	for _, deny := range formatDenyList {
		if allowSet[deny] {
			continue
		}
		if seen[deny] {
			continue
		}
		seen[deny] = true
		result = append(result, deny)
	}
	return result
}
