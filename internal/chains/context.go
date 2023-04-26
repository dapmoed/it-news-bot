package chains

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Context определяет контекст выполнения цепочки обработчиков Telegram-бота.
type Context struct {
	Update tgbotapi.Update // Обновление Telegram-бота.
	Chain  *Chain          // Текущая цепочка обработчиков.
	values []KeyValue      // Произвольные значения, связанные с контекстом.
}

// KeyValue представляет пару ключ-значение и используется для хранения произвольных значений в контексте.
type KeyValue struct {
	Key   string      // Ключ значения.
	Value interface{} // Значение, которое связано с ключом.
}

// Set добавляет новое значение в контекст или обновляет существующее значение с заданным ключом.
func (c *Context) Set(key string, value interface{}) {
	// Проверяем, существует ли уже значение с таким ключом в контексте.
	if i, v := c.Get(key); v != nil {
		// Если значение уже существует, то обновляем его.
		c.values[i] = KeyValue{
			Value: value,
			Key:   key,
		}
	}

	// Добавляем новое значение в контекст.
	c.values = append(c.values, KeyValue{
		Value: value,
		Key:   key,
	})
}

// Get возвращает индекс и значение, связанные с заданным ключом.
func (c *Context) Get(key string) (int, interface{}) {
	// Проходим по всем значениям в контексте.
	for i, v := range c.values {
		if v.Key == key {
			// Если находим значение с заданным ключом, то возвращаем его индекс и значение.
			return i, v.Value
		}
	}
	// Если значение с заданным ключом не найдено, то возвращаем нулевой индекс и nil.
	return 0, nil
}
