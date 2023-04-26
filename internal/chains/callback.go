package chains

import (
	// Импорт сторонних библиотек
	jsoniter "github.com/json-iterator/go" // Библиотека для работы с JSON
	"github.com/mitchellh/mapstructure"    // Библиотека для декодирования map[string]interface{} в структуры
	"reflect"                              // Пакет для работы с типами данных во время выполнения программы
)

// Определяем структуру CallbackCommand, которая будет использоваться для передачи команд и данных
type CallbackCommand struct {
	Name string      `json:"command"` // Имя команды
	Data interface{} `json:"data"`    // Данные
}

// Функция для создания новой команды с указанным именем и данными
func NewCallbackData(command string, data interface{}) CallbackCommand {
	return CallbackCommand{
		Name: command,
		Data: data,
	}
}

// Метод для сериализации объекта CallbackCommand в формат JSON
func (c CallbackCommand) JSON() string {
	b, err := jsoniter.Marshal(&c) // Преобразуем структуру в байтовый массив с помощью jsoniter
	if err != nil {
		return "{}"
	}
	return string(b) // Возвращаем сериализованные данные в формате строки
}

// Функция для десериализации JSON-строки в объект CallbackCommand
func UnmarshalCallbackCommand(s string) (string, interface{}, error) {
	callbackCommand := &CallbackCommand{}                 // Создаем объект структуры CallbackCommand
	err := jsoniter.Unmarshal([]byte(s), callbackCommand) // Декодируем JSON-строку в объект с помощью jsoniter
	if err != nil {
		return "", nil, err // Возвращаем ошибку, если декодирование не удалось
	}

	return callbackCommand.Name, callbackCommand.Data, nil // Возвращаем имя команды и данные
}

// Функция для десериализации объекта с данными в нужную структуру
func UnmarshalCallbackData(s interface{}, dst interface{}) (interface{}, error) {
	callbackDataDst := reflect.New(reflect.TypeOf(dst)).Interface() // Создаем объект типа dst с помощью рефлексии

	err := mapstructure.Decode(s, callbackDataDst) // Декодируем данные в объект с помощью mapstructure
	if err != nil {
		return nil, err // Возвращаем ошибку, если декодирование не удалось
	}

	return callbackDataDst, nil // Возвращаем десериализованные данные
}
