package main

type Storage interface {
	StoreSquidForuserId(userId string, squid string) error
	GetSquidForUserId(userId string) (string, error)
	IsReactableMessage(messageId string) (bool, error)
	StoreReactableMessage(messageId string) error
	Close()
}
