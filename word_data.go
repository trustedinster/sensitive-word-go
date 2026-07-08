package sensitiveword

import "sync"

// wordDataTreeNode DFA 树节点，对应 Java 的 WordDataTreeNode。
type wordDataTreeNode struct {
	end       bool                   // 关键词结束标识
	subNodeMap map[rune]*wordDataTreeNode // 子节点
}

func newWordDataTreeNode() *wordDataTreeNode {
	return &wordDataTreeNode{}
}

func (n *wordDataTreeNode) End() bool { return n.end }

func (n *wordDataTreeNode) SetEnd(end bool) { n.end = end }

func (n *wordDataTreeNode) GetSubNode(c rune) *wordDataTreeNode {
	if n.subNodeMap == nil {
		return nil
	}
	return n.subNodeMap[c]
}

// NodeSize 子节点数量。
func (n *wordDataTreeNode) NodeSize() int {
	if n.subNodeMap == nil {
		return 0
	}
	return len(n.subNodeMap)
}

func (n *wordDataTreeNode) ClearNode() { n.subNodeMap = nil }

func (n *wordDataTreeNode) RemoveNode(c rune) {
	if n.subNodeMap == nil {
		return
	}
	delete(n.subNodeMap, c)
}

func (n *wordDataTreeNode) AddSubNode(c rune, sub *wordDataTreeNode) {
	if n.subNodeMap == nil {
		n.subNodeMap = make(map[rune]*wordDataTreeNode)
	}
	n.subNodeMap[c] = sub
}

func (n *wordDataTreeNode) Destroy() {
	if n.subNodeMap != nil {
		n.subNodeMap = nil
	}
}

// WordData 敏感词 DFA 树数据结构，对应 Java 的 WordDataTree。
// 线程安全：InitWordData/AddWord/RemoveWord 均加锁。
type WordData struct {
	mu   sync.RWMutex
	root *wordDataTreeNode
}

// NewWordData 创建空的 DFA 数据结构。
func NewWordData() *WordData {
	return &WordData{root: newWordDataTreeNode()}
}

// InitWordData 用给定词集合全量重建 DFA。
// 替换是原子的：先构建新树再替换根。
func (w *WordData) InitWordData(words []string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	newRoot := newWordDataTreeNode()
	for _, word := range words {
		if word == "" {
			continue
		}
		w.addWordLocked(newRoot, word)
	}
	w.root = newRoot
}

// AddWord 增量新增敏感词。
func (w *WordData) AddWord(words []string) {
	if len(words) == 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, word := range words {
		if word == "" {
			continue
		}
		w.addWordLocked(w.root, word)
	}
}

// RemoveWord 增量删除敏感词。
func (w *WordData) RemoveWord(words []string) {
	if len(words) == 0 {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, word := range words {
		if word == "" {
			continue
		}
		w.removeWordLocked(w.root, word)
	}
}

// Contains 判断 builder 中的字符序列是否在 DFA 中命中。
// 返回包含类别（NotFound / Prefix / End）。
func (w *WordData) Contains(builder []rune, innerCtx *InnerSensitiveWordContext) WordContainsType {
	if len(builder) == 0 {
		return WordContainsNotFound
	}
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.root == nil {
		return WordContainsNotFound
	}
	nowNode := w.root
	for i := 0; i < len(builder); i++ {
		nowNode = w.getNowMap(nowNode, i, builder, innerCtx)
		if nowNode == nil {
			return WordContainsNotFound
		}
	}
	if nowNode.end {
		return WordContainsEnd
	}
	return WordContainsPrefix
}

// Destroy 释放内存。
func (w *WordData) Destroy() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.root != nil {
		w.root.Destroy()
		w.root = nil
	}
}

// getNowMap 获取当前下标对应的子节点。
// 处理 ignoreRepeat：若启用且与上一字符相同，则停留在当前节点。
func (w *WordData) getNowMap(nowNode *wordDataTreeNode, index int, builder []rune, innerCtx *InnerSensitiveWordContext) *wordDataTreeNode {
	ctx := innerCtx.wordContext
	mappingChar := builder[index]

	current := nowNode.GetSubNode(mappingChar)
	if ctx.IgnoreRepeat() && index > 0 {
		pre := builder[index-1]
		if pre == mappingChar {
			current = nowNode
		}
	}
	return current
}

// addWordLocked 在指定根节点下新增词。
func (w *WordData) addWordLocked(root *wordDataTreeNode, word string) {
	temp := root
	for _, c := range word {
		sub := temp.GetSubNode(c)
		if sub == nil {
			sub = newWordDataTreeNode()
			temp.AddSubNode(c, sub)
		}
		temp = sub
	}
	temp.SetEnd(true)
}

// removeWordLocked 删除指定词。
func (w *WordData) removeWordLocked(root *wordDataTreeNode, word string) {
	temp := root
	// 需要删除的节点路径：字符 -> 父节点
	type pair struct {
		ch   rune
		parent *wordDataTreeNode
	}
	path := make([]pair, 0, len(word))
	runes := []rune(word)
	length := len(runes)
	for i, c := range runes {
		sub := temp.GetSubNode(c)
		if sub == nil {
			return
		}
		if i == length-1 {
			if !sub.end {
				return
			}
			if sub.NodeSize() > 0 {
				sub.SetEnd(false)
				return
			}
		}
		if sub.end {
			// 遇到结束节点，清空之前的路径(无法安全删除)
			path = path[:0]
		}
		path = append(path, pair{ch: c, parent: temp})
		temp = sub
	}
	for _, p := range path {
		if p.parent.NodeSize() == 1 {
			p.parent.ClearNode()
			return
		}
		p.parent.RemoveNode(p.ch)
	}
}
