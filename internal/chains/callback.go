package chains

import jsoniter "github.com/json-iterator/go"

type CallbackData struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}

func NewCallbackData(command string, data interface{}) CallbackData {
	return CallbackData{
		Command: command,
		Data:    data,
	}
}

func (c CallbackData) JSON() string {
	b, err := jsoniter.Marshal(&c)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func UnmarshalCallbackData(s string) (string, interface{}, error) {
	c := &CallbackData{}
	err := jsoniter.Unmarshal([]byte(s), c)
	if err != nil {
		return "", nil, err
	}

	return c.Command, c.Data, nil
}
