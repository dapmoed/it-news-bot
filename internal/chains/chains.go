package chains

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chain struct {
	command string
	steps   []Step
	current int
}

func NewChain(command string) *Chain {
	return &Chain{
		command: command,
	}
}

func (c *Chain) Call(update tgbotapi.Update) {
	c.steps[c.current].Call(Context{
		chain:  c,
		update: update,
	})
}

func (c *Chain) Register(f func(ctx Context)) *Chain {
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

type Command struct {
	bot *tgbotapi.BotAPI
}

func NewCommand(bot *tgbotapi.BotAPI) *Command {
	return &Command{
		bot: bot,
	}
}

func (c *Command) Start(ctx Context) {
	msg := tgbotapi.NewMessage(ctx.update.Message.Chat.ID, "StartOne")
	c.bot.Send(msg)
	ctx.chain.Next()
}

func (c *Command) Start2(ctx Context) {
	msg := tgbotapi.NewMessage(ctx.update.Message.Chat.ID, "StartTwo")
	c.bot.Send(msg)
	ctx.chain.Next()
}
func (c *Command) Start3(ctx Context) {
	msg := tgbotapi.NewMessage(ctx.update.Message.Chat.ID, "StartThree")
	c.bot.Send(msg)
	ctx.chain.Next()
}
func (c *Command) Start4(ctx Context) {
	msg := tgbotapi.NewMessage(ctx.update.Message.Chat.ID, "StartFour")
	c.bot.Send(msg)
}

type Context struct {
	update tgbotapi.Update
	chain  *Chain
}
