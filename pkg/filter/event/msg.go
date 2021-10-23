package event

// Msq Request Action Type Enum

type MQType string

type MQAction int

const (
	MQActionPublish = 1 + iota
	MQActionSubscribe
	MQActionUnSubscribe
	MQActionConsumeAck
)

var MQActionIntToStr = map[MQAction]string{
	MQActionPublish:     "publish",
	MQActionSubscribe:   "subscribe",
	MQActionUnSubscribe: "unsubscribe",
	MQActionConsumeAck:  "consumeack",
}

var MQActionStrToInt = map[string]MQAction{
	"publish":     MQActionPublish,
	"subscribe":   MQActionSubscribe,
	"unsubscribe": MQActionUnSubscribe,
	"consumeack":  MQActionConsumeAck,
}

// MQRequest url format http://domain/publish/broker/topic
type MQRequest struct {
	Partition  int      `json:"partition"`   // for kafka
	ConsumeUrl string   `json:"consume_url"` // not empty when subscribe msg, eg: http://10.0.0.1:11451/consume
	CheckUrl   string   `json:"check_url"`   // not empty when subscribe msg, eg: http://10.0.0.1:11451/health
	Msg        []string `json:"msg"`         // not empty when publish msg, msg body
}

type MQMsgPush struct {
	Msg []string `json:"msg"`
}
