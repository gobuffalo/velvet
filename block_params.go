package velvet

type blockParams struct {
	current []string
	stack   [][]string
}

func newBlockParams() *blockParams {
	return &blockParams{
		current: []string{},
		stack:   [][]string{},
	}
}

func (bp *blockParams) push(params []string) {
	bp.current = params
	bp.stack = append(bp.stack, params)
}

func (bp *blockParams) pop() []string {
	l := len(bp.stack)
	if l == 0 {
		return bp.current
	}
	p := bp.stack[l-1]
	bp.stack = bp.stack[0:(l - 1)]
	l = len(bp.stack)
	if l == 0 {
		bp.current = []string{}
	} else {
		bp.current = bp.stack[l-1]
	}
	return p
}
