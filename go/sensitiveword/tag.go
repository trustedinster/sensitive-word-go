package sensitiveword

// WordTag 敏感词标签策略接口，对应 Java 的 IWordTag。
//
// 用于查询某个敏感词所属的标签集合，便于分类管理与过滤。
type WordTag interface {
	// GetTag 返回 word 对应的标签集合；无标签时返回 nil。
	GetTag(word string) []string
}

// ===== 实现 =====

// noneWordTag 空标签实现，对应 Java 的 NoneWordTag。
type noneWordTag struct{}

// NoneWordTag 单例。
var NoneWordTag = &noneWordTag{}

func (t *noneWordTag) GetTag(word string) []string {
	return nil
}

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

// ===== 工具类 =====

// WordTags 工具类，对应 Java 的 WordTags。
var WordTags = wordTagsHelper{}

type wordTagsHelper struct{}

// None 返回空标签策略。
func (wordTagsHelper) None() WordTag { return NoneWordTag }

// Defaults 返回默认标签策略（空标签）。
func (wordTagsHelper) Defaults() WordTag { return NoneWordTag }

// Map 返回基于映射的标签策略。
func (wordTagsHelper) Map(tagMap map[string][]string) WordTag {
	return NewWordTagMap(tagMap)
}
