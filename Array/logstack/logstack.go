package logstack

type LogStack interface {
	Insert(key int)
	Delete(key int)
	Search(key int) bool
}

type logStack struct {
	cache layer
	limit uint
	stack []layer
}

func NewLogStack(limit uint) LogStack {
	if limit < 4 {
		limit = 4
	}
	var ls = new(logStack)
	ls.limit = limit
	return ls
}

func (ls *logStack) change(key int, mark bool) {
	ls.cache.change(key, mark)
	if ls.cache.size() == ls.limit {
		ls.stack = append(ls.stack, ls.cache)
		ls.cache.reset()
	}
}
func (ls *logStack) Insert(key int) {
	ls.change(key, true)
}
func (ls *logStack) Delete(key int) {
	ls.change(key, false)
}

func (ls *logStack) Search(key int) bool {
	var found = ls.cache.search(key)
	if found == 0 && len(ls.stack) != 0 {
		found = ls.stack[len(ls.stack)-1].search(key)
	}
	for found == 0 && len(ls.stack) > 1 {
		found = ls.stack[len(ls.stack)-2].search(key)
		ls.stack[len(ls.stack)-2].merge(&ls.stack[len(ls.stack)-1])
		ls.stack = ls.stack[:len(ls.stack)-1]
	}
	if len(ls.stack) == 1 {
		ls.stack[0].compact()
	}
	return found > 0
}

func (ls *logStack) debug() {
	ls.cache.print()
	for i := len(ls.stack) - 1; i >= 0; i-- {
		ls.stack[i].print()
	}
}
