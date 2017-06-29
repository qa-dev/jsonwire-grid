package jsonwire

type ClientFactoryInterface interface {
	Create(address string) ClientInterface
}

type ClientInterface interface {
	Status() (*Message, error)
	Sessions() (*Sessions, error)
	CloseSession(sessionId string) (*Message, error)
	Address() string
}

type NodeInterface interface {
	RemoveAllSessions() (int, error)
}
