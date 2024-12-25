package pusher

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/mritd/logger"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
)

type APNS struct {
	clients chan *apns2.Client
	topic   string
}

func NewAPNS(
	maxClientCount int,
	apnsPrivateKey, topic, keyID, teamID string,
	apnsCAs []string,
) (*APNS, error) {
	if maxClientCount < 1 {
		return nil, fmt.Errorf("invalid number of clients")
	}

	authKey, err := token.AuthKeyFromBytes([]byte(apnsPrivateKey))
	if err != nil {
		logger.Fatalf("failed to create APNS auth key: %v", err)
	}

	var rootCAs *x509.CertPool
	if runtime.GOOS == "windows" {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			logger.Fatalf("failed to get rootCAs: %v", err)
		}
	}

	for _, ca := range apnsCAs {
		rootCAs.AppendCertsFromPEM([]byte(ca))
	}

	clients := make(chan *apns2.Client, maxClientCount)

	for i := 0; i < min(runtime.NumCPU(), maxClientCount); i++ {
		client := &apns2.Client{
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
		logger.Infof("create apns client: %d", i)
		clients <- client
	}

	logger.Info("init apns client success...")
	return &APNS{
		clients: clients,
		topic:   topic,
	}, nil
}

func (a *APNS) Push(message Message) error {
	var (
		msg *ApnsMessage
		ok  bool
	)
	if msg, ok = message.(*ApnsMessage); !ok {
		return fmt.Errorf("invalid message type")
	}
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertSubtitle(msg.Subtitle).
		AlertBody(msg.Body).
		Sound(msg.Sound).
		Category(msg.Category).
		ThreadID(msg.Group)

	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}

	client := <-a.clients // grab a client from the pool
	a.clients <- client   // add the client back to the pool

	resp, err := client.Push(&apns2.Notification{
		DeviceToken: msg.DeviceToken,
		Topic:       a.topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNS push failed: %s", resp.Reason)
	}
	return nil
}
