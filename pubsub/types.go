package pubsub


type Option struct {}

type PubSub interface {
	New(opt Option) error
	PubMsg() error
	SubChannel(topic string) error
	Ack() error
}
