package chain

type chainsLink func(ctx Context)

type Context interface {
	Next(chainsLink)
	Close()
}

type TgContext struct {
}

func NewContext() Context {
	return &TgContext{}
}

func (c *TgContext) Next(chainsLink) {

}

func (c *TgContext) Close() {

}
