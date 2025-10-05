package session

import (
	"DHT/src/models"
	"sync"
)

type Session struct {
	Node *models.Node
}

var (
	instance *Session
	once     sync.Once
)

func GetSession() *Session {
	once.Do(func() {
		instance = &Session{
			Node: &models.Node{},
		}
	})
	return instance
}
