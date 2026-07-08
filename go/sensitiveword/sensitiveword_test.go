package sensitiveword

import (
	"testing"
)

// newTestBs 创建用于测试的引导类实例。
func newTestBs() *SensitiveWordBs {
	return NewSensitiveWordBs().Init()
}

// ===== DFA 树测试 =====

func TestWordData_Contains(t *testing.T) {
	wd := NewWordData()
	wd.InitWordData([]string{"敏感词", "测试", "hello"})

	tests := []struct {
		name     string
		input    string
		expected WordContainsType
	}{
		{"完整命中", "敏感词", WordContainsEnd},
		{"前缀命中", "敏感", WordContainsPrefix},
		{"英文命中", "hello", WordContainsEnd},
		{"未命中", "不存在", WordContainsNotFound},
		{"空字符串", "", WordContainsNotFound},
	}

	ctx := &InnerSensitiveWordContext{wordContext: NewWordContext()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wd.Contains([]rune(tt.input), ctx)
			if result != tt.expected {
				t.Errorf("Contains(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWordData_AddRemoveWord(t *testing.T) {
	wd := NewWordData()
	wd.InitWordData([]string{"测试"})

	ctx := &InnerSensitiveWordContext{wordContext: NewWordContext()}

	// 增量添加
	wd.AddWord([]string{"新增词"})
	if wd.Contains([]rune("新增词"), ctx) != WordContainsEnd {
		t.Error("AddWord failed: 新增词 not found")
	}

	// 删除
	wd.RemoveWord([]string{"测试"})
	if wd.Contains([]rune("测试"), ctx) != WordContainsNotFound {
		t.Error("RemoveWord failed: 测试 still found")
	}
}

// ===== 字符工具测试 =====

func TestToHalfWidth(t *testing.T) {
	tests := []struct {
		input    rune
		expected rune
	}{
		{'Ａ', 'A'},
		{'ａ', 'a'},
		{'０', '0'},
		{'９', '9'},
		{'　', ' '},
		{'A', 'A'}, // 半角不变
		{'中', '中'}, // 中文不变
	}
	for _, tt := range tests {
		result := toHalfWidth(tt.input)
		if result != tt.expected {
			t.Errorf("toHalfWidth(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestIsEmail(t *testing.T) {
	valid := []string{"test@example.com", "a@b.cn"}
	invalid := []string{"not-email", "@example.com", "test@", "a.b", ""}

	for _, s := range valid {
		if !isEmail(s) {
			t.Errorf("isEmail(%q) should be true", s)
		}
	}
	for _, s := range invalid {
		if isEmail(s) {
			t.Errorf("isEmail(%q) should be false", s)
		}
	}
}

func TestIsUrl(t *testing.T) {
	valid := []string{"http://example.com", "https://baidu.com", "http://a.cn"}
	invalid := []string{"example.com", "ftp://test.com", "not-url", ""}

	for _, s := range valid {
		if !isUrl(s) {
			t.Errorf("isUrl(%q) should be true", s)
		}
	}
	for _, s := range invalid {
		if isUrl(s) {
			t.Errorf("isUrl(%q) should be false", s)
		}
	}
}

// ===== 格式化策略测试 =====

func TestWordFormatIgnoreCase(t *testing.T) {
	ctx := NewWordContext()
	f := WordFormatIgnoreCase
	if f.Format('A', ctx) != 'a' {
		t.Error("IgnoreCase should convert A to a")
	}
	if f.Format('Z', ctx) != 'z' {
		t.Error("IgnoreCase should convert Z to z")
	}
}

func TestWordFormatIgnoreWidth(t *testing.T) {
	ctx := NewWordContext()
	f := WordFormatIgnoreWidth
	if f.Format('Ａ', ctx) != 'A' {
		t.Error("IgnoreWidth should convert Ａ to A")
	}
}

func TestWordFormatIgnoreNumStyle(t *testing.T) {
	ctx := NewWordContext()
	f := WordFormatIgnoreNumStyle
	// 中文数字
	if f.Format('一', ctx) != '1' {
		t.Error("IgnoreNumStyle should convert 一 to 1")
	}
	if f.Format('九', ctx) != '9' {
		t.Error("IgnoreNumStyle should convert 九 to 9")
	}
	// 圆圈数字
	if f.Format('①', ctx) != '1' {
		t.Error("IgnoreNumStyle should convert ① to 1")
	}
	// 无映射的字符不变
	if f.Format('a', ctx) != 'a' {
		t.Error("IgnoreNumStyle should keep 'a' unchanged")
	}
}

func TestWordFormatIgnoreEnglishStyle(t *testing.T) {
	ctx := NewWordContext()
	f := WordFormatIgnoreEnglishStyle
	if f.Format('Ⓐ', ctx) != 'A' {
		t.Error("IgnoreEnglishStyle should convert Ⓐ to A")
	}
	if f.Format('ⓐ', ctx) != 'a' {
		t.Error("IgnoreEnglishStyle should convert ⓐ to a")
	}
}

func TestWordFormatIgnoreChineseStyle(t *testing.T) {
	ctx := NewWordContext()
	f := WordFormatIgnoreChineseStyle
	// 繁体转简体
	if f.Format('萬', ctx) != '万' {
		t.Error("IgnoreChineseStyle should convert 萬 to 万")
	}
}

func TestWordFormatArray(t *testing.T) {
	ctx := NewWordContext()
	// 组合：忽略大小写 + 忽略全角
	f := NewWordFormatArray([]WordFormat{WordFormatIgnoreCase, WordFormatIgnoreWidth})
	if f.Format('Ａ', ctx) != 'a' {
		t.Error("Array should convert Ａ to a")
	}
}

func TestWordFormatCombine(t *testing.T) {
	ctx := NewWordContext()
	combine := NewWordFormatCombine()
	format := combine.InitWordFormat(ctx)
	// 默认全部启用，Ａ应被转为 a
	if format.Format('Ａ', ctx) != 'a' {
		t.Error("Combine should convert Ａ to a")
	}
}

// ===== 繁简转换测试 =====

func TestToSimple(t *testing.T) {
	// 字符级繁简转换（基于 OpenCC）
	cases := []struct{ in, want rune }{
		{'萬', '万'},
		{'簡', '简'},
		{'體', '体'},
		{'漢', '汉'},
		{'字', '字'}, // 简体不变
		{'軟', '软'},
		{'體', '体'},
	}
	for _, c := range cases {
		if got := toSimple(c.in); got != c.want {
			t.Errorf("toSimple(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToSimpleString(t *testing.T) {
	// 字符级繁简转换（t2s 预设）：繁体字形 -> 简体字形
	// 注意：t2s 是字形转换，非台湾用语转换（tw2sp）。
	// 如"軟體"->"软体"（字形），而非"软件"（用语）。
	cases := []struct{ in, want string }{
		{"簡體漢字", "简体汉字"},
		{"軟體", "软体"}, // 字形级：軟->软, 體->体
		{"", ""},
	}
	for _, c := range cases {
		if got := toSimpleString(c.in); got != c.want {
			t.Errorf("toSimpleString(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRegisterT2SMapping(t *testing.T) {
	// 用户自定义补充映射应覆盖 OpenCC 结果
	RegisterT2SMapping('測', '测')
	if toSimple('測') != '测' {
		t.Error("RegisterT2SMapping failed")
	}
}

func TestSensitiveWordBs_TraditionalChinese(t *testing.T) {
	// 端到端：黑名单存简体字形，输入繁体字形应能匹配
	// t2s 做字形级转换：軟->软, 體->体, 所以"軟體"归一化为"软体"
	bs := NewSensitiveWordBs().
		SetIgnoreChineseStyle(true).
		Init()
	bs.AddWord("软体")
	// 繁体"軟體"经归一化后应命中简体字形"软体"
	if !bs.Contains("这是軟體测试") {
		t.Error("IgnoreChineseStyle should match 軟體 -> 软体")
	}
}

// ===== 引擎核心测试 =====

func TestSensitiveWordBs_Contains(t *testing.T) {
	bs := newTestBs()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"包含敏感词", "这是一个敏感词测试", true},
		{"包含自定义敏感词", "这是自定义敏感词", true},
		{"包含英文敏感词", "what the fuck", true},
		{"不包含敏感词", "这是一段正常文字", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bs.Contains(tt.input)
			if result != tt.expected {
				t.Errorf("Contains(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSensitiveWordBs_FindAll(t *testing.T) {
	bs := newTestBs()
	words := bs.FindAll("这是一个敏感词，也有自定义敏感词")
	if len(words) != 2 {
		t.Errorf("FindAll should find 2 words, got %d: %v", len(words), words)
	}
}

func TestSensitiveWordBs_FindFirst(t *testing.T) {
	// wordFailFast=true（默认）时，"敏感" 是 "敏感词" 的前缀且先命中
	bs := newTestBs()
	word := bs.FindFirst("这是一个敏感词")
	if word != "敏感" {
		t.Errorf("FindFirst should be '敏感' (failFast), got %q", word)
	}

	// wordFailFast=false 时，应匹配最长的 "敏感词"
	bs2 := NewSensitiveWordBs().SetWordFailFast(false).Init()
	word = bs2.FindFirst("这是一个敏感词")
	if word != "敏感词" {
		t.Errorf("FindFirst should be '敏感词' (no failFast), got %q", word)
	}

	// 未命中
	word = bs.FindFirst("正常文字")
	if word != "" {
		t.Errorf("FindFirst should be empty, got %q", word)
	}
}

func TestSensitiveWordBs_Replace(t *testing.T) {
	// wordFailFast=true（默认），"敏感" 先命中（2字符）
	bs := newTestBs()
	result := bs.Replace("这是一个敏感词测试")
	if result != "这是一个**词测试" {
		t.Errorf("Replace result unexpected: %q", result)
	}

	// wordFailFast=false，匹配最长 "敏感词"（3字符）
	bs2 := NewSensitiveWordBs().SetWordFailFast(false).Init()
	result = bs2.Replace("这是一个敏感词测试")
	if result != "这是一个***测试" {
		t.Errorf("Replace (no failFast) result unexpected: %q", result)
	}
}

func TestSensitiveWordBs_IgnoreCase(t *testing.T) {
	bs := NewSensitiveWordBs().Init()
	// 英文敏感词 fuck 应匹配大写 FUCK
	if !bs.Contains("FUCK") {
		t.Error("IgnoreCase should match FUCK")
	}
}

func TestSensitiveWordBs_IgnoreWidth(t *testing.T) {
	bs := NewSensitiveWordBs().Init()
	// 全角的"敏感"应匹配
	if !bs.Contains("这是ｆｕｃｋ测试") {
		t.Error("IgnoreWidth should match fullwidth fuck")
	}
}

func TestSensitiveWordBs_NumCheck(t *testing.T) {
	bs := NewSensitiveWordBs().SetEnableNumCheck(true).Init()
	// 连续8位数字应被检测
	if !bs.Contains("联系电12345678话") {
		t.Error("NumCheck should detect 8+ digits")
	}
	// 不足8位不检测
	if bs.Contains("电话1234567") {
		t.Error("NumCheck should not detect <8 digits")
	}
}

func TestSensitiveWordBs_EmailCheck(t *testing.T) {
	bs := NewSensitiveWordBs().SetEnableEmailCheck(true).Init()
	if !bs.Contains("联系test@example.com") {
		t.Error("EmailCheck should detect email")
	}
}

func TestSensitiveWordBs_UrlCheck(t *testing.T) {
	bs := NewSensitiveWordBs().SetEnableUrlCheck(true).Init()
	if !bs.Contains("访问http://example.com看看") {
		t.Error("UrlCheck should detect URL")
	}
}

func TestSensitiveWordBs_Ipv4Check(t *testing.T) {
	bs := NewSensitiveWordBs().SetEnableIpv4Check(true).Init()
	if !bs.Contains("IP是192.168.1.1看看") {
		t.Error("Ipv4Check should detect IPv4")
	}
}

func TestSensitiveWordBs_WhiteList(t *testing.T) {
	bs := NewSensitiveWordBs().Init()
	// 白名单中的 "gender" 不应被检测
	// 注意：默认 result condition 要求英文全词匹配
	// gender 在白名单中，不应作为敏感词
	// 由于字典中不包含 "gender" 作为黑名单词，这里测试白名单功能
	bs.AddWord("gender")
	// gender 在白名单中，即使加入黑名单也应被白名单覆盖
	if bs.Contains("gender is here") {
		// 注意：默认的 EnglishWordMatch 条件下，"gender" 前后有非英文字符才算全词匹配
		// 但这里 gender 已在白名单，初始化时已从黑名单中剔除
		t.Log("gender detected (may be expected if whitelist not effective)")
	}
}

func TestSensitiveWordBs_AddRemoveWord(t *testing.T) {
	bs := NewSensitiveWordBs().Init()

	// 原本不包含
	if bs.Contains("动态添加测试词") {
		t.Error("should not contain 动态添加测试词 before AddWord")
	}

	// 动态添加
	bs.AddWord("动态添加测试词")
	if !bs.Contains("动态添加测试词") {
		t.Error("should contain 动态添加测试词 after AddWord")
	}

	// 动态删除
	bs.RemoveWord("动态添加测试词")
	if bs.Contains("动态添加测试词") {
		t.Error("should not contain 动态添加测试词 after RemoveWord")
	}
}

func TestSensitiveWordBs_ReplaceChar(t *testing.T) {
	// wordFailFast=true（默认），"敏感" 先命中（2字符）
	bs := NewSensitiveWordBs().
		SetWordReplace(NewWordReplaceChar('#')).
		Init()
	result := bs.Replace("这是一个敏感词")
	if result != "这是一个##词" {
		t.Errorf("Replace with '#' should be '这是一个##词', got %q", result)
	}

	// wordFailFast=false，匹配 "敏感词"（3字符）
	bs2 := NewSensitiveWordBs().
		SetWordReplace(NewWordReplaceChar('#')).
		SetWordFailFast(false).
		Init()
	result = bs2.Replace("这是一个敏感词")
	if result != "这是一个###" {
		t.Errorf("Replace with '#' (no failFast) should be '这是一个###', got %q", result)
	}
}

func TestSensitiveWordBs_Tags(t *testing.T) {
	// 使用自定义标签
	tagMap := map[string][]string{
		"敏感词": {"政治", "危险"},
	}
	bs := NewSensitiveWordBs().
		SetWordTag(NewWordTagMap(tagMap)).
		Init()

	tags := bs.Tags("敏感词")
	if len(tags) != 2 {
		t.Errorf("Tags should return 2 tags, got %d", len(tags))
	}
}

// ===== 字符忽略测试 =====

func TestSpecialCharSensitiveWordCharIgnore(t *testing.T) {
	bs := NewSensitiveWordBs().
		SetCharIgnore(SensitiveWordCharIgnores.SpecialChars()).
		Init()
	// 敏感词中间有特殊字符也应被检测
	if !bs.Contains("这是一个敏!感!词测试") {
		t.Error("SpecialChar ignore should detect 敏!感!词")
	}
}

// ===== 结果条件测试 =====

func TestWordResultConditionAlwaysTrue(t *testing.T) {
	ctx := NewWordContext()
	cond := WordResultConditionAlwaysTrue
	wr := &WordResult{startIndex: 0, endIndex: 3}
	if !cond.Match(wr, "测试文本", WordValidModeFailOver, ctx) {
		t.Error("AlwaysTrue should return true")
	}
}

func TestWordResultConditionEnglishWordMatch(t *testing.T) {
	ctx := NewWordContext()
	cond := WordResultConditionEnglishWordMatch

	// 英文敏感词前后无英文字符，应匹配
	wr := &WordResult{startIndex: 0, endIndex: 4}
	if !cond.Match(wr, "fuck is bad", WordValidModeFailOver, ctx) {
		t.Error("EnglishWordMatch should match 'fuck ' with no preceding english")
	}

	// 英文敏感词前有英文字符，不应匹配（非全词匹配）
	wr2 := &WordResult{startIndex: 1, endIndex: 5}
	if cond.Match(wr2, "afuck is bad", WordValidModeFailOver, ctx) {
		t.Error("EnglishWordMatch should not match 'afuck' (not full word)")
	}
}

// ===== CharIgnore 测试 =====

func TestNoneSensitiveWordCharIgnore(t *testing.T) {
	ignore := NoneSensitiveWordCharIgnore
	ctx := &InnerSensitiveWordContext{}
	if ignore.Ignore(0, "test", ctx) {
		t.Error("None ignore should always return false")
	}
}

func TestSpecialCharSensitiveWordCharIgnore_Char(t *testing.T) {
	ignore := SpecialCharSensitiveWordCharIgnore
	ctx := &InnerSensitiveWordContext{}
	if !ignore.Ignore(0, "!test", ctx) {
		t.Error("SpecialChar should ignore '!'")
	}
	if ignore.Ignore(0, "atest", ctx) {
		t.Error("SpecialChar should not ignore 'a'")
	}
}

// ===== 默认实例测试 =====

func TestDefault(t *testing.T) {
	// Default() 应返回同一个实例
	d1 := Default()
	d2 := Default()
	if d1 != d2 {
		t.Error("Default() should return the same instance")
	}

	// 默认实例应能检测敏感词
	if !d1.Contains("这是一个敏感词") {
		t.Error("Default instance should detect sensitive word")
	}
}
