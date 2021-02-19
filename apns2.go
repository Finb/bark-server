package main

import (
	"sync"

	"github.com/sideshow/apns2/payload"

	"github.com/mritd/logger"
	"github.com/sideshow/apns2/token"

	"github.com/sideshow/apns2"
)

type PushInfo struct {
	DeviceKey string `form:"device_key,omitempty" json:"device_key,omitempty"`
	Category  string `form:"category,omitempty,default=bark" json:"category,omitempty,default=bark"`
	Title     string `form:"title" json:"title"`
	// ios sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	Sound string `form:"sound,omitempty,default=1107" json:"sound,omitempty,default=1107"`
	Badge string `form:"badge,omitempty" json:"badge,omitempty"`
	Body  string `form:"body,omitempty,default=nobody" json:"body,omitempty,default=nobody"`

	// compatible with old req
	OldDeviceKey string `form:"key,omitempty" json:"key,omitempty"`

	DeviceToken string `form:"-" json:"-"`
}

const (
	keyID  = "LH4T9V5U4R"
	teamID = "5U8LBRXG3A"
)

const topic = "me.fin.bark"

var apns2Once sync.Once
var apns2Client *apns2.Client
var keyBs = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 77, 73, 71, 84, 65, 103, 69, 65, 77, 66, 77, 71, 66, 121, 113, 71, 83, 77, 52, 57, 65, 103, 69, 71, 67, 67, 113, 71, 83, 77, 52, 57, 65, 119, 69, 72, 66, 72, 107, 119, 100, 119, 73, 66, 65, 81, 81, 103, 52, 118, 116, 67, 51, 103, 53, 76, 53, 72, 103, 75, 71, 74, 50, 43, 10, 84, 49, 101, 65, 48, 116, 79, 105, 118, 82, 69, 118, 69, 65, 89, 50, 103, 43, 106, 117, 82, 88, 74, 107, 89, 76, 50, 103, 67, 103, 89, 73, 75, 111, 90, 73, 122, 106, 48, 68, 65, 81, 101, 104, 82, 65, 78, 67, 65, 65, 83, 109, 79, 115, 51, 74, 107, 83, 121, 111, 71, 69, 87, 90, 10, 115, 85, 71, 120, 70, 115, 47, 52, 112, 119, 49, 114, 73, 108, 83, 86, 50, 73, 67, 49, 57, 77, 56, 117, 51, 71, 53, 107, 113, 51, 54, 117, 112, 79, 119, 121, 70, 87, 106, 57, 71, 105, 51, 69, 106, 99, 57, 100, 51, 115, 67, 55, 43, 83, 72, 82, 113, 88, 114, 69, 65, 74, 111, 119, 10, 56, 47, 55, 116, 82, 112, 86, 43, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45}

func apns2Setup() {
	apns2Once.Do(func() {
		authKey, err := token.AuthKeyFromBytes(keyBs)
		if err != nil {
			logger.Fatalf("failed to create APNS auth key: %v", err)
		}

		apnsClient = apns2.NewTokenClient(&token.Token{
			AuthKey: authKey,
			KeyID:   keyID,
			TeamID:  teamID,
		}).Production()
	})
}

func apnsPush(info *PushInfo) error {
	noti := apns2.Notification{
		DeviceToken: info.DeviceToken,
		Topic:       topic,
	}

	pl := payload.NewPayload().
		Sound(info.Sound).
		Category(info.Category).
		Badge(info.Badge)
}
