package storage

type Storage interface {
	StoreSquidForUserId(userId, squid string) error
	GetSquidForUserId(userId string) (string, error)
	IsReactableMessage(messageId string) (bool, error)
	StoreReactableMessage(messageId string) error
	Close()
}
