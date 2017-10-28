package jsonwire

// ClientFactoryInterface - is an abstract http-client factory for node implementations.
type ClientFactoryInterface interface {
	Create(address string) ClientInterface
}

// ClientInterface - is an abstract http-client for node implementations.
type ClientInterface interface {
	Health() (*Message, error)
	Sessions() (*Sessions, error)
	CloseSession(sessionID string) (*Message, error)
	Address() string
}
