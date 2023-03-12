package chains

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	ErrEndOfScript      = errors.New("end of script")
	ErrNotFoundCallback = errors.New("not found callback")
	ErrNotFoundStep     = errors.New("not found step")
)

const (
	defaultDurationSession = 10 * time.Second
)

// Chain implements the ability to organize a dialogue with the bot
// from multiple related steps
// also implemented the ability to respond to the buttons under the messages
type Chain struct {
	ended           bool
	current         int
	durationSession time.Duration
	context         *Context
	steps           []Step
	callbacks       map[string]StepCallback
}

// NewChain initializes the Chain structure
func NewChain() *Chain {
	return &Chain{
		durationSession: defaultDurationSession,
		callbacks:       make(map[string]StepCallback),
	}
}

// Clone creates an instance for further processing of messages from the user by script
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

// DurationSession returns the duration of possible user inactivity in the script
func (c *Chain) DurationSession() time.Duration {
	return c.durationSession
}

// SetDurationSession sets the duration of user inactivity in the script
func (c *Chain) SetDurationSession(duration time.Duration) *Chain {
	c.durationSession = duration
	return c
}

// Call causes the corresponding script step to be executed
// c.current - current script step
func (c *Chain) Call(update tgbotapi.Update) error {
	ctx := c.context
	ctx.Update = update
	if c.current < len(c.steps) {
		c.steps[c.current].Call(ctx)
		return nil
	}
	return ErrNotFoundStep
}

// CallCallback calling message button command handler
func (c *Chain) CallCallback(update tgbotapi.Update) error {
	ctx := c.context
	ctx.Update = update

	if update.CallbackQuery != nil {
		command, data, err := UnmarshalCallbackCommand(update.CallbackQuery.Data)
		if err != nil {
			return err
		}
		if c, ok := c.callbacks[command]; ok {
			c.Call(ctx, data)
			return nil
		}
	}

	return ErrNotFoundCallback
}

// RegisterCallback registers a command handler for clicking buttons below the message
func (c *Chain) RegisterCallback(command string, f func(ctx *Context, data interface{})) *Chain {
	c.callbacks[command] = StepCallback{
		f: f,
	}
	return c
}

// RegisterStep registers the step in the script
func (c *Chain) RegisterStep(f func(ctx *Context)) *Chain {
	c.steps = append(c.steps, Step{
		f: f,
	})
	return c
}

// Next tells the chain that the user's next message,
// you need to use the next handler in the queue
func (c *Chain) Next() error {
	if len(c.steps) <= c.current+1 {
		return ErrEndOfScript
	}
	c.current++
	return nil
}

// End tells the script to end script processing
func (c *Chain) End() {
	c.ended = true
}

// IsEnded returns whether the script is complete
func (c *Chain) IsEnded() bool {
	return c.ended
}

// Step implements a step in the script
type Step struct {
	f func(ctx *Context)
}

// Call the handler corresponding to the step
func (s Step) Call(ctx *Context) {
	s.f(ctx)
}

// StepCallback stores a function to handle an event on a button click below the message
type StepCallback struct {
	f func(ctx *Context, data interface{})
}

// Call alls a function to handle the action, passing call parameters
func (c StepCallback) Call(ctx *Context, data interface{}) {
	c.f(ctx, data)
}
