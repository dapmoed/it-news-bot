package chains

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type Pool struct {
	chains map[string]Chain
}

func NewPool() *Pool {
	return &Pool{
		chains: make(map[string]Chain),
	}
}

func (p *Pool) Command(c string, chain Chain) *Pool {
	p.chains[c] = chain
	return p
}

func (p *Pool) GetChain(command string) (*Chain, error) {
	if val, ok := p.chains[command]; ok {
		return val.Clone(), nil
	}
	return nil, ErrNotFound
}
