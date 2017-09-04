package jsonwire

type ClientFactoryInterface interface {
	Create(address string) ClientInterface
}

type ClientInterface interface {
	Health() (*Message, error)
	Sessions() (*Sessions, error)
	CloseSession(sessionID string) (*Message, error)
	Address() string
}

type NodeInterface interface {
	RemoveAllSessions() (int, error)
}
