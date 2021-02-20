package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/http2"

	"github.com/sideshow/apns2/payload"

	"github.com/mritd/logger"
	"github.com/sideshow/apns2/token"

	"github.com/sideshow/apns2"
)

type PushMessage struct {
	DeviceKey string `form:"device_key,omitempty" json:"device_key,omitempty" xml:"device_key,omitempty"`
	Category  string `form:"category,omitempty" json:"category,omitempty" xml:"category,omitempty"`
	Title     string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
	Body      string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
	// ios notification sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	Sound     string            `form:"sound,omitempty" json:"sound,omitempty" xml:"sound,omitempty"`
	ExtParams map[string]string `form:"ext_params,omitempty" json:"ext_params,omitempty" xml:"ext_params,omitempty"`

	DeviceToken string `form:"-" json:"-" xml:"-"`
}

const (
	topic  = "me.fin.bark"
	keyID  = "LH4T9V5U4R"
	teamID = "5U8LBRXG3A"
)

var apnsOnce sync.Once
var apnsClient *apns2.Client

func apnsSetup() {
	apnsOnce.Do(func() {
		authKey, err := token.AuthKeyFromBytes([]byte(apnsPrivateKey))
		if err != nil {
			logger.Fatalf("failed to create APNS auth key: %v", err)
		}

		rootCAs, _ := x509.SystemCertPool()
		for _, ca := range apnsCAs {
			rootCAs.AppendCertsFromPEM([]byte(ca))
		}

		apnsClient = &apns2.Client{
			Token: &token.Token{
				AuthKey: authKey,
				KeyID:   keyID,
				TeamID:  teamID,
			},
			HTTPClient: &http.Client{
				Transport: &http2.Transport{
					DialTLS: apns2.DialTLS,
					TLSClientConfig: &tls.Config{
						RootCAs: rootCAs,
					},
				},
				Timeout: apns2.HTTPClientTimeout,
			},
			Host: apns2.HostProduction,
		}

	})
}

func apnsPush(msg *PushMessage) error {
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Sound(msg.Sound).
		Category(msg.Category)

	for k, v := range msg.ExtParams {
		pl.Custom(k, v)
	}

	resp, err := apnsClient.Push(&apns2.Notification{
		DeviceToken: msg.DeviceToken,
		Topic:       topic,
		Payload:     pl.MutableContent(),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNS push failed: %s", resp.Reason)
	}
	return nil
}
