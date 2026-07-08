package sensitiveword

import (
	"sync"

	"github.com/yanmingcao/opencc-go"
	"github.com/yanmingcao/opencc-go/pkg/config"
	"github.com/yanmingcao/opencc-go/pkg/embeddata"
)

// 繁简体转换基于 OpenCC-Go (https://github.com/yanmingcao/opencc-go) 实现。
//
// 使用其内嵌的 t2s（繁体 -> 简体）预设，无需任何外部数据文件。
// 同时保留 RegisterT2SMapping / RegisterT2SMappingBatch 供用户在 OpenCC
// 结果之上追加自定义映射（OpenCC 字典为只读，无法动态扩展）。

var (
	t2sConverterOnce sync.Once
	t2sConv          *opencc.SimpleConverter
	t2sConvErr       error

	// userT2SMapping 用户自定义繁简补充映射，在 OpenCC 转换结果之上覆盖。
	// 读取时加读锁，注册时加写锁。
	userT2SMapping   = map[rune]rune{}
	userT2SMappingMu sync.RWMutex
)

// initT2SConverter 懒加载 OpenCC t2s converter。
func initT2SConverter() {
	t2sConverterOnce.Do(func() {
		data, err := embeddata.GetConfig("t2s")
		if err != nil {
			t2sConvErr = err
			return
		}
		cfg, err := config.LoadConfigFromData(data, "")
		if err != nil {
			t2sConvErr = err
			return
		}
		t2sConv, t2sConvErr = opencc.NewSimpleConverterFromConfig(cfg)
	})
}

// t2sConverter 返回已初始化的 OpenCC t2s converter。
// 初始化失败时返回 nil（调用方应处理）。
func t2sConverter() *opencc.SimpleConverter {
	initT2SConverter()
	return t2sConv
}

// RegisterT2SMapping 注册额外的繁简映射，用于在 OpenCC 结果之上扩展。
// 已存在的映射会被覆盖。线程安全。
//
// 说明：OpenCC 内嵌字典为只读，无法动态扩展；本函数提供一个轻量补充入口，
// 用于满足特殊业务字符的繁简归一化需求。
func RegisterT2SMapping(trad, simp rune) {
	userT2SMappingMu.Lock()
	userT2SMapping[trad] = simp
	userT2SMappingMu.Unlock()
}

// RegisterT2SMappingBatch 批量注册繁简映射。
func RegisterT2SMappingBatch(pairs map[rune]rune) {
	userT2SMappingMu.Lock()
	for k, v := range pairs {
		userT2SMapping[k] = v
	}
	userT2SMappingMu.Unlock()
}

// toSimple 将单个字符转换为简体，对应 Java ZhSlimUtil.toSimple。
//
// 查询顺序：
//  1. 用户自定义补充映射（RegisterT2SMapping 注册）
//  2. OpenCC t2s 字典
//  3. 无匹配则原样返回
func toSimple(c rune) rune {
	// 1. 用户自定义补充映射优先
	userT2SMappingMu.RLock()
	if mc, ok := userT2SMapping[c]; ok {
		userT2SMappingMu.RUnlock()
		return mc
	}
	userT2SMappingMu.RUnlock()

	// 2. OpenCC 转换
	conv := t2sConverter()
	if conv == nil {
		// 初始化失败，降级为原样返回
		return c
	}
	result := conv.Convert(string(c))
	rs := []rune(result)
	if len(rs) == 0 {
		return c
	}
	return rs[0]
}

// toSimpleString 将整段字符串转换为简体。
//
// 直接使用 OpenCC 整段转换，可利用其最大前缀分词，对词组级转换更准确
// （如"軟體" -> "软件"而非逐字结果）。
func toSimpleString(s string) string {
	if s == "" {
		return s
	}
	conv := t2sConverter()
	if conv == nil {
		return s
	}
	result := conv.Convert(s)

	// 应用用户自定义补充映射（覆盖 OpenCC 结果中的对应字符）
	userT2SMappingMu.RLock()
	if len(userT2SMapping) == 0 {
		userT2SMappingMu.RUnlock()
		return result
	}
	changed := false
	rs := []rune(result)
	for i, c := range rs {
		if mc, ok := userT2SMapping[c]; ok {
			rs[i] = mc
			changed = true
		}
	}
	userT2SMappingMu.RUnlock()
	if !changed {
		return result
	}
	return string(rs)
}
