package pubsub

type natsClient struct {
}

func (n *natsClient) PubMsg() error {
	//TODO implement me
	panic("implement me")
}

func (n *natsClient) SubChannel(topic string) error {
	//TODO implement me
	panic("implement me")
}

func (n *natsClient) Ack() error {
	//TODO implement me
	panic("implement me")
}

func NewNatsClient() PubSub {
	return &natsClient{}
}
