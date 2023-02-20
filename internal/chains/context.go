package chains

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	Update tgbotapi.Update
	Chain  *Chain
	values []KeyValue
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func (c *Context) Set(key string, value interface{}) {
	if i, v := c.Get(key); v != nil {
		c.values[i] = KeyValue{
			Value: value,
			Key:   key,
		}
	}

	c.values = append(c.values, KeyValue{
		Value: value,
		Key:   key,
	})
}

func (c *Context) Get(key string) (int, interface{}) {
	for i, v := range c.values {
		if v.Key == key {
			return i, v.Value
		}
	}
	return 0, nil
}
