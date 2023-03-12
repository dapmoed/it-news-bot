package chains

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	keyTestCall   = "key_test_call"
	valueTestCall = "value_test_call"

	keyTestCallback   = "key_test_callback"
	valueTestCallback = "value_test_callback"
)

func GetTestChain() *Chain {
	chain := NewChain().
		RegisterStep(func(ctx *Context) {
			ctx.Set(keyTestCall, valueTestCall)
		}).
		RegisterCallback("test_command_callback", func(ctx *Context, data interface{}) {
			ctx.Set(keyTestCallback, valueTestCallback)
		})
	return chain
}

func GetTestChainMultipleStep() *Chain {
	chain := NewChain().
		RegisterStep(func(ctx *Context) {
			ctx.Set(keyTestCall, valueTestCall)
		}).
		RegisterStep(func(ctx *Context) {
			ctx.Set(keyTestCall, valueTestCall)
		}).
		RegisterStep(func(ctx *Context) {
			ctx.Set(keyTestCall, valueTestCall)
		}).
		RegisterCallback("test_command_callback", func(ctx *Context, data interface{}) {
			ctx.Set(keyTestCallback, valueTestCallback)
		})
	return chain
}

func TestNewChain(t *testing.T) {
	chain := NewChain()

	assert.NotNil(t, chain)
	assert.NotNil(t, chain.callbacks)
	assert.Equal(t, defaultDurationSession, chain.durationSession)
}

func TestChain_Clone(t *testing.T) {
	chain := GetTestChain()
	cloneChain := chain.Clone()

	assert.NotNil(t, cloneChain)
	assert.NotNil(t, cloneChain.callbacks)
	assert.Len(t, cloneChain.callbacks, 1)
	assert.NotNil(t, cloneChain.steps)
	assert.Len(t, cloneChain.steps, 1)
	assert.Equal(t, defaultDurationSession, cloneChain.durationSession)
}

func TestChain_DurationSession(t *testing.T) {
	chain := GetTestChain()
	duration := chain.DurationSession()

	assert.IsType(t, time.Second, duration)
	assert.Equal(t, defaultDurationSession, duration)
}

func TestChain_SetDurationSession(t *testing.T) {
	chain := GetTestChain()
	chain.SetDurationSession(time.Second * 666)

	assert.Equal(t, time.Second*666, chain.durationSession)
}

func TestChain_Call(t *testing.T) {
	chain := GetTestChain()

	t.Run("success", func(t *testing.T) {
		cloneChain := chain.Clone()
		err := cloneChain.Call(tgbotapi.Update{})
		_, key := cloneChain.context.Get(keyTestCall)
		keyValueString, ok := key.(string)
		assert.Equal(t, tgbotapi.Update{}, cloneChain.context.Update)
		assert.Equal(t, true, ok)
		assert.Equal(t, valueTestCall, keyValueString)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		cloneChain := chain.Clone()
		cloneChain.current = 10
		err := cloneChain.Call(tgbotapi.Update{})
		assert.ErrorIs(t, err, ErrNotFoundStep)
	})
}

func TestChain_CallCallback(t *testing.T) {
	chain := GetTestChain()
	update := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Data: `{"command":"test_command_callback","data":"test_data"}`,
		},
	}

	updateTwo := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Data: `{"command":"test_command_callback_2","data":"test_data"}`,
		},
	}

	updateErrorData := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Data: `test`,
		},
	}

	t.Run("success", func(t *testing.T) {
		cloneChain := chain.Clone()
		err := cloneChain.CallCallback(update)
		_, key := cloneChain.context.Get(keyTestCallback)
		keyValueString, ok := key.(string)
		assert.Equal(t, update, cloneChain.context.Update)
		assert.Equal(t, true, ok)
		assert.Equal(t, valueTestCallback, keyValueString)
		assert.NoError(t, err)
	})

	t.Run("error callback query is nil", func(t *testing.T) {
		cloneChain := chain.Clone()
		err := cloneChain.CallCallback(tgbotapi.Update{})
		assert.ErrorIs(t, ErrNotFoundCallback, err)
	})

	t.Run("error not fund callback", func(t *testing.T) {
		cloneChain := chain.Clone()
		err := cloneChain.CallCallback(updateTwo)
		assert.ErrorIs(t, ErrNotFoundCallback, err)
	})

	t.Run("error unmarshal callback data", func(t *testing.T) {
		cloneChain := chain.Clone()
		err := cloneChain.CallCallback(updateErrorData)
		assert.Error(t, err)
	})
}

func TestChain_RegisterStep(t *testing.T) {
	chain := GetTestChain()
	assert.Len(t, chain.steps, 1)
}

func TestChain_RegisterCallback(t *testing.T) {
	chain := GetTestChain()
	assert.Len(t, chain.callbacks, 1)
}

func TestChain_Next(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		chain := GetTestChainMultipleStep()
		cloneChain := chain.Clone()
		err := cloneChain.Next()
		assert.NoError(t, err)
	})

	t.Run("success", func(t *testing.T) {
		chain := GetTestChain()
		cloneChain := chain.Clone()
		err := cloneChain.Next()
		assert.ErrorIs(t, ErrEndOfScript, err)
	})
}

func TestChain_End(t *testing.T) {
	chain := GetTestChain()
	cloneChain := chain.Clone()
	cloneChain.ended = true
	cloneChain.End()
	assert.Equal(t, true, cloneChain.ended)
}

func TestChain_IsEnded(t *testing.T) {
	chain := GetTestChain()
	cloneChain := chain.Clone()
	cloneChain.ended = true
	result := cloneChain.IsEnded()
	assert.Equal(t, true, result)
}
