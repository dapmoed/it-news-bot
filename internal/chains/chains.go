package chains

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chain struct {
	steps   []Step
	current int
}

func NewChain() Chain {
	return Chain{}
}

func (c Chain) Clone() *Chain {
	return &Chain{
		steps: c.steps,
	}
}

func (c *Chain) Call(update tgbotapi.Update) {
	c.steps[c.current].Call(Context{
		Chain:  c,
		Update: update,
	})
}

func (c Chain) Register(f func(ctx Context)) Chain {
	c.steps = append(c.steps, Step{
		f: f,
	})
	return c
}

func (c *Chain) Next() {
	c.current++
}

type Step struct {
	f func(ctx Context)
}

func (s Step) Call(ctx Context) {
	s.f(ctx)
}

type Context struct {
	Update tgbotapi.Update
	Chain  *Chain
}
