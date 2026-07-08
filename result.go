package sensitiveword

// WordLengthResult 黑白名单长度结果，对应 Java 的 WordLengthResult。
//
// 一次遍历同时统计黑名单与白名单命中长度，避免重复扫描。
type WordLengthResult struct {
	wordAllowLen int    // 白名单命中长度
	wordDenyLen  int    // 黑名单命中长度
	wordDeny     string // 黑名单实际匹配词
	wordAllow    string // 白名单实际匹配词
}

// NewWordLengthResult 创建空结果。
func NewWordLengthResult() *WordLengthResult {
	return &WordLengthResult{}
}

func (r *WordLengthResult) WordAllowLen() int     { return r.wordAllowLen }
func (r *WordLengthResult) WordDenyLen() int      { return r.wordDenyLen }
func (r *WordLengthResult) WordDeny() string      { return r.wordDeny }
func (r *WordLengthResult) WordAllow() string     { return r.wordAllow }

func (r *WordLengthResult) SetWordAllowLen(v int) *WordLengthResult { r.wordAllowLen = v; return r }
func (r *WordLengthResult) SetWordDenyLen(v int) *WordLengthResult  { r.wordDenyLen = v; return r }
func (r *WordLengthResult) SetWordDeny(s string) *WordLengthResult  { r.wordDeny = s; return r }
func (r *WordLengthResult) SetWordAllow(s string) *WordLengthResult { r.wordAllow = s; return r }

// WordCheckResult 检测结果，对应 Java 的 WordCheckResult。
type WordCheckResult struct {
	wordLengthResult *WordLengthResult
	checkType        WordType // 命中类别，对应 Java 的 type/checkClass
}

// NewWordCheckResult 创建空结果。
func NewWordCheckResult() *WordCheckResult {
	return &WordCheckResult{}
}

func (r *WordCheckResult) WordLengthResult() *WordLengthResult { return r.wordLengthResult }
func (r *WordCheckResult) CheckType() WordType                 { return r.checkType }

func (r *WordCheckResult) SetWordLengthResult(v *WordLengthResult) *WordCheckResult {
	r.wordLengthResult = v
	return r
}
func (r *WordCheckResult) SetCheckType(t WordType) *WordCheckResult {
	r.checkType = t
	return r
}

// WordResult 敏感词匹配结果，对应 Java 的 WordResult / IWordResult。
type WordResult struct {
	startIndex int
	endIndex   int
	checkType  WordType
	word       string
}

// NewWordResult 创建空结果。
func NewWordResult() *WordResult {
	return &WordResult{}
}

func (r *WordResult) StartIndex() int   { return r.startIndex }
func (r *WordResult) EndIndex() int     { return r.endIndex }
func (r *WordResult) CheckType() WordType { return r.checkType }
func (r *WordResult) Word() string      { return r.word }

func (r *WordResult) SetStartIndex(v int) *WordResult     { r.startIndex = v; return r }
func (r *WordResult) SetEndIndex(v int) *WordResult       { r.endIndex = v; return r }
func (r *WordResult) SetCheckType(t WordType) *WordResult { r.checkType = t; return r }
func (r *WordResult) SetWord(s string) *WordResult        { r.word = s; return r }
