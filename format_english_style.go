package sensitiveword

// englishStyleOne 英文字母样式源字符集合，对应 Java WordFormatIgnoreEnglishStyleC2C.LETTERS_ONE。
//
// 共 3 组，每组 26 个字符，与 englishStyleTwo 一一对应：
//
//  1. ⒶⒷⒸ...Ⓩ  (实心圆圈大写) -> ABC...XYZ
//  2. ⓐⓑⓒ...ⓩ  (实心圆圈小写) -> abc...xyz
//  3. ⒜⒝⒞...⒵  (圆括号小写)   -> abc...xyz
const englishStyleOne = "ⒶⒷⒸⒹⒺⒻⒼⒽⒾⒿⓀⓁⓂⓃⓄⓅⓆⓇⓈⓉⓊⓋⓌⓍⓎⓏ" +
	"ⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ" +
	"⒜⒝⒞⒟⒠⒡⒢⒣⒤⒥⒦⒧⒨⒩⒪⒫⒬⒭⒮⒯⒰⒱⒲⒳⒴⒵"

// englishStyleTwo 英文字母样式目标字符集合，对应 Java WordFormatIgnoreEnglishStyleC2C.LETTERS_TWO。
const englishStyleTwo = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"abcdefghijklmnopqrstuvwxyz"

// englishStyleMap 英文字母样式映射表，对应 Java 的 LETTER_MAP。
var englishStyleMap = func() map[rune]rune {
	m := make(map[rune]rune, len(englishStyleOne))
	src := []rune(englishStyleOne)
	dst := []rune(englishStyleTwo)
	for i := 0; i < len(src) && i < len(dst); i++ {
		m[src[i]] = dst[i]
	}
	return m
}()

// wordFormatIgnoreEnglishStyle 忽略英文字母样式，将各类带圈/带括号字母统一为普通字母，
// 对应 Java 的 WordFormatIgnoreEnglishStyleC2C。
type wordFormatIgnoreEnglishStyle struct{}

// WordFormatIgnoreEnglishStyle 单例。
var WordFormatIgnoreEnglishStyle = &wordFormatIgnoreEnglishStyle{}

func (f *wordFormatIgnoreEnglishStyle) Format(original rune, context *WordContext) rune {
	return getMappingChar(englishStyleMap, original)
}
