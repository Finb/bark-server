package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Message interface {
	Check() error
	SetDeviceToken(deviceToken string)
	LoadMessages(params interface{}) error
}

type ApnsMessage struct {
	DeviceToken string `form:"device_token,omitempty" json:"device_token,omitempty" xml:"device_token,omitempty" query:"device_token,omitempty"`

	Title    string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty" query:"title,omitempty"`
	Subtitle string `form:"subtitle,omitempty" json:"subtitle,omitempty" xml:"subtitle,omitempty" query:"subtitle,omitempty"`
	Body     string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty" query:"body,omitempty"`

	// ios pusher sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	Sound string `form:"sound,omitempty" json:"sound,omitempty" xml:"sound,omitempty" query:"sound,omitempty"`
	Group string `form:"group,omitempty" json:"group,omitempty" xml:"group,omitempty" query:"group,omitempty"`

	Category  string                 `form:"category,omitempty" json:"category,omitempty" xml:"category,omitempty" query:"category,omitempty"`
	ExtParams map[string]interface{} `form:"ext_params,omitempty" json:"ext_params,omitempty" xml:"ext_params,omitempty" query:"ext_params,omitempty"`
}

func (apns *ApnsMessage) SetDeviceToken(deviceToken string) {
	apns.DeviceToken = deviceToken
}
func (apns *ApnsMessage) Check() error {
	if apns.Body == "" && apns.Title == "" && apns.Subtitle == "" {
		return errors.New("Empty message")
	}
	return nil
}
func (apns *ApnsMessage) LoadMessages(params interface{}) error {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal params: %v", err)
		}
		if err := json.Unmarshal(jsonBytes, &paramsMap); err != nil {
			return fmt.Errorf("failed to unmarshal params to map: %v", err)
		}
	}
	for key, value := range paramsMap {
		switch key {
		case "device_token":
			apns.DeviceToken, _ = value.(string)
		case "title":
			apns.Title, _ = value.(string)
		case "subtitle":
			apns.Subtitle, _ = value.(string)
		case "body":
			apns.Body, _ = value.(string)
		case "sound":
			apns.Sound, _ = value.(string)
		case "group":
			apns.Group, _ = value.(string)
		default:
			if apns.ExtParams == nil {
				apns.ExtParams = make(map[string]interface{})
			}
			apns.ExtParams[key] = value
		}
	}
	fmt.Println(apns)

	return nil
}

const (
	DefaultApnsCategory = "myNotificationCategory"
	DefaultApnsSound    = "1107"
)

func NewApnsMessage() Message {
	return &ApnsMessage{
		Category:  DefaultApnsCategory,
		Sound:     DefaultApnsSound,
		ExtParams: make(map[string]interface{}),
	}
}
