package sensitiveword

import (
	"strings"
	"sync"
)

// WordTag 敏感词标签策略接口，对应 Java 的 IWordTag。
//
// 用于查询某个敏感词所属的标签集合，便于分类管理与过滤。
type WordTag interface {
	// GetTag 返回 word 对应的标签集合；无标签时返回 nil。
	GetTag(word string) []string
}

// ===== 空实现 =====

// noneWordTag 空标签实现，对应 Java 的 NoneWordTag。
type noneWordTag struct{}

// NoneWordTag 单例。
var NoneWordTag = &noneWordTag{}

func (t *noneWordTag) GetTag(word string) []string {
	return nil
}

// ===== 基于内存映射的实现 =====

// wordTagMap 基于内存映射的标签实现，对应 Java 的 WordTagMap。
type wordTagMap struct {
	tagMap map[string][]string
}

// NewWordTagMap 基于给定映射创建标签策略。
func NewWordTagMap(tagMap map[string][]string) WordTag {
	return &wordTagMap{tagMap: tagMap}
}

func (t *wordTagMap) GetTag(word string) []string {
	if word == "" {
		return nil
	}
	if tags, ok := t.tagMap[word]; ok {
		return tags
	}
	return nil
}

// ===== 基于行的实现 =====

// wordTagLines 按标准行解析的标签实现，对应 Java 的 WordTagLines。
//
// 行格式："单词 标签1,标签2"（默认 wordSplit=" ", tagSplit=","）。
// 构建时一次性解析所有行，后续查询走内存 map，O(1)。
type wordTagLines struct {
	tagMap map[string][]string
}

// NewWordTagLines 按指定分隔符从行构建标签策略，对应 Java 的 WordTagLines(lines, wordSplit, tagSplit)。
func NewWordTagLines(lines []string, wordSplit, tagSplit string) WordTag {
	return &wordTagLines{tagMap: buildWordTagMap(lines, wordSplit, tagSplit)}
}

// NewWordTagLinesDefault 使用默认分隔符（空格、逗号）从行构建标签策略，对应 Java 的 WordTagLines(lines)。
func NewWordTagLinesDefault(lines []string) WordTag {
	return NewWordTagLines(lines, " ", ",")
}

func (t *wordTagLines) GetTag(word string) []string {
	if word == "" {
		return nil
	}
	if tags, ok := t.tagMap[word]; ok {
		return tags
	}
	return nil
}

// buildWordTagMap 解析行集合，构建 word -> tags 映射。
// 行格式："单词<wordSplit>标签1<tagSplit>标签2<tagSplit>..."。
// 使用 SplitN 仅按 wordSplit 分割一次，保证单词内部不会因包含分隔符而被错误切断。
func buildWordTagMap(lines []string, wordSplit, tagSplit string) map[string][]string {
	tagMap := make(map[string][]string, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, wordSplit, 2)
		if len(parts) < 2 {
			continue
		}
		word := strings.TrimSpace(parts[0])
		if word == "" {
			continue
		}
		var tags []string
		for _, t := range strings.Split(parts[1], tagSplit) {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
		if len(tags) > 0 {
			tagMap[word] = tags
		}
	}
	return tagMap
}

// ===== 系统默认实现 =====

// wordTagSystem 系统默认标签策略，对应 Java 的 WordTagSystem。
//
// 从内嵌的 sensitive_word_tags.txt 加载（go:embed），懒加载（sync.Once）：
// 首次调用 GetTag 时才解析 4 万余行标签数据并构建 map，避免无标签需求时的开销。
type wordTagSystem struct {
	once    sync.Once
	wordTag WordTag
}

// WordTagSystem 单例。
var WordTagSystem = &wordTagSystem{}

func (t *wordTagSystem) GetTag(word string) []string {
	t.once.Do(func() {
		lines := readDictLines("sensitive_word_tags.txt")
		t.wordTag = NewWordTagLinesDefault(lines)
	})
	return t.wordTag.GetTag(word)
}

// ===== 工具类 =====

// WordTags 工具类，对应 Java 的 WordTags。
var WordTags = wordTagsHelper{}

type wordTagsHelper struct{}

// None 返回空标签策略。
func (wordTagsHelper) None() WordTag { return NoneWordTag }

// System 返回系统默认标签策略（从内嵌 sensitive_word_tags.txt 加载，懒加载）。
func (wordTagsHelper) System() WordTag { return WordTagSystem }

// Defaults 返回默认标签策略（与 Java 一致，返回系统默认策略）。
func (wordTagsHelper) Defaults() WordTag { return WordTagSystem }

// Map 返回基于映射的标签策略。
func (wordTagsHelper) Map(tagMap map[string][]string) WordTag {
	return NewWordTagMap(tagMap)
}

// Lines 返回基于行的标签策略（默认分隔符：空格、逗号）。
func (wordTagsHelper) Lines(lines []string) WordTag {
	return NewWordTagLinesDefault(lines)
}

// LinesWithSplit 返回基于行的标签策略，使用自定义分隔符。
func (wordTagsHelper) LinesWithSplit(lines []string, wordSplit, tagSplit string) WordTag {
	return NewWordTagLines(lines, wordSplit, tagSplit)
}
