package pubsub

type Option struct{}

type PubSub interface {
	PubMsg() error
	SubChannel(topic string) error
	Ack() error
}
