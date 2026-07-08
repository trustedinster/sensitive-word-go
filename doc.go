// Package sensitiveword 提供基于 DFA 算法的高性能敏感词检测能力。
//
// 本包是 Java 项目 sensitive-word (https://github.com/houbb/sensitive-word) 的 Go 重构实现。
//
// 特性：
//   - 基于 DFA 算法，高性能匹配
//   - 支持敏感词判断、返回、脱敏
//   - 支持大小写、全角半角、数字、繁简体、英文样式等多种格式化忽略
//   - 支持敏感词、邮箱、数字、URL、IPv4 等多种检测策略
//   - 支持自定义替换策略
//   - 支持用户自定义黑/白名单与动态加载
//   - 支持敏感词标签
//   - 支持跳过特殊字符
//   - 支持单个黑/白名单的增删，无需全量初始化
//   - 支持快速失败 / 全量遍历两种匹配模式
//
// 繁简体转换基于 OpenCC-Go (https://github.com/yanmingcao/opencc-go) 实现，
// 使用其内嵌的 t2s 预设（纯 Go，无需外部数据文件），覆盖完整字符级繁简映射。
// 可通过 RegisterT2SMapping / RegisterT2SMappingBatch 在 OpenCC 结果之上追加自定义映射。
//
// 快速开始：
//
//	// 方式一：使用默认实例（懒加载，首次调用自动初始化）
//	if sensitiveword.Default().Contains("敏感词测试") {
//	    fmt.Println("包含敏感词")
//	}
//	words := sensitiveword.Default().FindAll("敏感词测试")
//	replaced := sensitiveword.Default().Replace("敏感词测试")
//
//	// 方式二：自定义配置（fluent-api 风格）
//	bs := sensitiveword.NewSensitiveWordBs().
//	    SetIgnoreCase(true).
//	    SetIgnoreChineseStyle(true).
//	    SetEnableEmailCheck(true).
//	    SetEnableNumCheck(true).
//	    Init()
//	defer bs.Destroy()
//
//	if bs.Contains("联系 test@example.com") {
//	    fmt.Println("包含敏感信息")
//	}
//
//	// 动态增删敏感词
//	bs.AddWord("新敏感词")
//	bs.RemoveWord("新敏感词")
package sensitiveword

import "sync"

// defaultBs 默认引导类实例（懒加载）。
var (
	defaultBs     *SensitiveWordBs
	defaultBsOnce sync.Once
)

// Default 返回默认的 SensitiveWordBs 实例（懒加载，线程安全）。
// 首次调用时自动执行 Init()，使用系统内置字典与默认配置。
func Default() *SensitiveWordBs {
	defaultBsOnce.Do(func() {
		defaultBs = NewSensitiveWordBs().Init()
	})
	return defaultBs
}
