package chains

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	ErrEndOfScript = errors.New("end of script")
)

type Chain struct {
	steps           []Step
	current         int
	context         *Context
	ended           bool
	durationSession time.Duration
}

func NewChain() *Chain {
	return &Chain{
		durationSession: time.Second * 10,
	}
}

func (c *Chain) Clone() *Chain {
	chain := &Chain{
		steps:           c.steps,
		durationSession: c.DurationSession(),
	}
	chain.context = &Context{
		Chain: chain,
	}
	return chain
}

func (c *Chain) DurationSession() time.Duration {
	return c.durationSession
}

func (c *Chain) SetDurationSession(duration time.Duration) *Chain {
	c.durationSession = duration
	return c
}

func (c *Chain) Call(update tgbotapi.Update) {
	ctx := c.context
	ctx.Update = update
	c.steps[c.current].Call(ctx)
}

func (c *Chain) Register(f func(ctx *Context)) *Chain {
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
