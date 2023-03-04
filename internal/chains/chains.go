package chains

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	ErrEndOfScript      = errors.New("end of script")
	ErrNotFoundCallback = errors.New("not found callback")
)

type Chain struct {
	steps           []Step
	callbacks       map[string]Step
	current         int
	context         *Context
	ended           bool
	durationSession time.Duration
}

func NewChain() *Chain {
	return &Chain{
		durationSession: time.Second * 10,
		callbacks:       make(map[string]Step),
	}
}

func (c *Chain) Clone() *Chain {
	chain := &Chain{
		steps:           c.steps,
		callbacks:       c.callbacks,
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

func (c *Chain) CallBack(update tgbotapi.Update) error {
	ctx := c.context
	ctx.Update = update

	if update.CallbackQuery != nil {
		if c, ok := c.callbacks[update.CallbackQuery.Data]; ok {
			c.Call(ctx)
			return nil
		}
	}

	return ErrNotFoundCallback
}

func (c *Chain) RegisterCallback(data string, f func(ctx *Context)) *Chain {
	c.callbacks[data] = Step{
		f: f,
	}
	return c
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
