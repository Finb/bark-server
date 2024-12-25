package pusher

type Pusher interface {
	Push(message Message) error
}

type PushProvider string

const PushProviderApns PushProvider = "apns"

// todo const PushProviderFcm PushProvider = "fcm"
// todo const PushProviderHarmony PushProvider = "harmony"

var pusherList = make(map[PushProvider]Pusher)

func GetPusher(pushProvider PushProvider, params ...interface{}) (pusher Pusher) {
	if pusher, ok := pusherList[pushProvider]; ok {
		return pusher
	}
	switch pushProvider {
	case PushProviderApns:
		maxClientCount := ApnsClientCount
		if len(params) > 0 {
			maxClientCount = params[0].(int)
		}
		if apns, err := NewAPNS(
			maxClientCount,
			ApnsPrivateKey,
			ApnsTopic,
			ApnsKeyID,
			ApnsTeamID,
			ApnsCAs,
		); err != nil {
			return nil
		} else {
			pusherList[pushProvider] = apns
			return apns
		}
	}

	return nil
}
