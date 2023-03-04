package chains

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type CallbackCommand struct {
	Name string      `json:"command"`
	Data interface{} `json:"data"`
}

func NewCallbackData(command string, data interface{}) CallbackCommand {
	return CallbackCommand{
		Name: command,
		Data: data,
	}
}

func (c CallbackCommand) JSON() string {
	b, err := jsoniter.Marshal(&c)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func UnmarshalCallbackCommand(s string) (string, interface{}, error) {
	callbackCommand := &CallbackCommand{}
	err := jsoniter.Unmarshal([]byte(s), callbackCommand)
	if err != nil {
		return "", nil, err
	}

	return callbackCommand.Name, callbackCommand.Data, nil
}

func UnmarshalCallbackData(s interface{}, dst interface{}) (interface{}, error) {
	callbackDataDst := reflect.New(reflect.TypeOf(dst)).Interface()

	err := mapstructure.Decode(s, callbackDataDst)
	if err != nil {
		return nil, err
	}

	return callbackDataDst, nil
}
