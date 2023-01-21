package sweep

import "context"

type Manager struct {
	Cancel context.CancelFunc
	ctx    context.Context
}

func (m Manager) IsDone() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}

func (m Manager) Child() Manager {
	ctx, _ := context.WithCancel(m.ctx)
	return Manager{
		Cancel: m.Cancel,
		ctx:    ctx,
	}
}

func CreateManager() Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return Manager{
		Cancel: cancel,
		ctx:    ctx,
	}
}
