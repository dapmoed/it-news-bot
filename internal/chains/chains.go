package chains

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrEndOfScript = errors.New("end of script")
)

type Chain struct {
	steps   []Step
	current int
	context *Context
	ended   bool
}

func NewChain() Chain {
	return Chain{}
}

func (c Chain) Clone() *Chain {
	chain := &Chain{
		steps: c.steps,
	}
	chain.context = &Context{
		Chain: chain,
	}
	return chain
}

func (c *Chain) Call(update tgbotapi.Update) {
	ctx := c.context
	ctx.Update = update
	c.steps[c.current].Call(ctx)
}

func (c Chain) Register(f func(ctx *Context)) Chain {
	c.steps = append(c.steps, Step{
		f: f,
	})
	return c
}

func (c *Chain) Next() error {
	if len(c.steps) == c.current+1 {
		return ErrEndOfScript
	}
	c.current++
	return nil
}

func (c *Chain) End() {
	c.ended = true
}

func (c *Chain) IsEnded() bool {
	return c.ended
}

type Step struct {
	f func(ctx *Context)
}

func (s Step) Call(ctx *Context) {
	s.f(ctx)
}
