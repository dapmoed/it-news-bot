package c_command

import "it-news-bot/internal/chain"

type Command interface {
}

type StartCommand struct {
}

func (c *StartCommand) Start(ctx chain.Context) {
	ctx.Next(c.StepOne)
}

func (c *StartCommand) StepOne(ctx chain.Context) {
	ctx.Next(c.StepTwo)
}

func (c *StartCommand) StepTwo(ctx chain.Context) {
	ctx.Close()
}

func (c *StartCommand) CallbackOne(ctx chain.Context) {
	ctx.Close()
}
