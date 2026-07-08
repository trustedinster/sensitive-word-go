package sensitiveword

// numStyleOne 数字样式源字符集合，对应 Java WordFormatIgnoreNumStyleC2C.NUM_ONE。
//
// 每一类数字字符（共 17 组）按顺序排列，与 numStyleTwo 一一对应：
//
//  1. ⓪０零º₀⓿○          -> 0000000
//  2. １２３４５６７８９    -> 123456789
//  3. 一二三四五六七八九    -> 123456789
//  4. 壹贰叁肆伍陆柒捌玖    -> 123456789
//  5. ¹²³⁴⁵⁶⁷⁸⁹          -> 123456789
//  6. ₁₂₃₄₅₆₇₈₉          -> 123456789
//  7. ①②③④⑤⑥⑦⑧⑨        -> 123456789
//  8. ⑴⑵⑶⑷⑸⑹⑺⑻⑼        -> 123456789
//  9. ⒈⒉⒊⒋⒌⒍⒎⒏⒐        -> 123456789
// 10. ❶❷❸❹❺❻❼❽❾          -> 123456789
// 11. ➀➁➂➃➄➅➆➇➈        -> 123456789
// 12. ➊➋➌➍➎➏➐➑➒        -> 123456789
// 13. ㈠㈡㈢㈣㈤㈥㈦㈧㈨    -> 123456789
// 14. ⓵⓶⓷⓸⓹⓺⓻⓼⓽        -> 123456789
// 15. ㊀㊁㊂㊃㊄㊅㊆㊇㊈    -> 123456789
// 16. ⅰⅱⅲⅳⅴⅵⅶⅷⅸ        -> 123456789
// 17. ⅠⅡⅢⅣⅤⅥⅦⅧⅨ        -> 123456789
const numStyleOne = "⓪０零º₀⓿○" +
	"１２３４５６７８９" +
	"一二三四五六七八九" +
	"壹贰叁肆伍陆柒捌玖" +
	"¹²³⁴⁵⁶⁷⁸⁹" +
	"₁₂₃₄₅₆₇₈₉" +
	"①②③④⑤⑥⑦⑧⑨" +
	"⑴⑵⑶⑷⑸⑹⑺⑻⑼" +
	"⒈⒉⒊⒋⒌⒍⒎⒏⒐" +
	"❶❷❸❹❺❻❼❽❾" +
	"➀➁➂➃➄➅➆➇➈" +
	"➊➋➌➍➎➏➐➑➒" +
	"㈠㈡㈢㈣㈤㈥㈦㈧㈨" +
	"⓵⓶⓷⓸⓹⓺⓻⓼⓽" +
	"㊀㊁㊂㊃㊄㊅㊆㊇㊈" +
	"ⅰⅱⅲⅳⅴⅵⅶⅷⅸ" +
	"ⅠⅡⅢⅣⅤⅥⅦⅧⅨ"

// numStyleTwo 数字样式目标字符集合，对应 Java WordFormatIgnoreNumStyleC2C.NUM_TWO。
// 与 numStyleOne 等长且一一对应。
const numStyleTwo = "0000000" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789" +
	"123456789"

// numStyleMap 数字样式映射表，对应 Java 的 NUMBER_MAP。
var numStyleMap = func() map[rune]rune {
	m := make(map[rune]rune, len(numStyleOne))
	src := []rune(numStyleOne)
	dst := []rune(numStyleTwo)
	for i := 0; i < len(src) && i < len(dst); i++ {
		m[src[i]] = dst[i]
	}
	return m
}()

// wordFormatIgnoreNumStyle 忽略数字样式，将各类数字符号统一为 0-9，
// 对应 Java 的 WordFormatIgnoreNumStyleC2C。
type wordFormatIgnoreNumStyle struct{}

// WordFormatIgnoreNumStyle 单例。
var WordFormatIgnoreNumStyle = &wordFormatIgnoreNumStyle{}

func (f *wordFormatIgnoreNumStyle) Format(original rune, context *WordContext) rune {
	return getMappingChar(numStyleMap, original)
}
