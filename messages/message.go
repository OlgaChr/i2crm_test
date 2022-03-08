package messages

import (
	"sync"
	"time"
)

type MessageKey string

type Message struct {
	Key           MessageKey
	UnixTimestamp uint64
}

type History interface {
	Upsert(message *Message)
	Last() *Message
	Delete(key MessageKey)
}

type ChatStorage struct {
	storage  map[MessageKey]*storageMessage
	capacity int
	mu       *sync.Mutex
}

type storageMessage struct {
	message   *Message
	createdAt int64
}

func NewChatStorage(n int) *ChatStorage {
	return &ChatStorage{
		storage:  make(map[MessageKey]*storageMessage, n),
		capacity: n,
		mu:       &sync.Mutex{},
	}
}

func (cs *ChatStorage) addMessage(message *Message) {
	cs.storage[message.Key] = &storageMessage{
		message:   message,
		createdAt: time.Now().UnixNano(),
	}
}

func (cs *ChatStorage) deleteOldestMessage() {
	var message *storageMessage

	for _, m := range cs.storage {
		if message == nil {
			message = m
			continue
		}
		if m.message.UnixTimestamp < message.message.UnixTimestamp {
			message = m
			continue
		}
		if m.message.UnixTimestamp == message.message.UnixTimestamp && m.createdAt < message.createdAt {
			message = m
		}
	}

	delete(cs.storage, message.message.Key)
}

func (cs *ChatStorage) Upsert(message *Message) {
	cs.mu.Lock()
	// если такое сообщение уже существует, то ничего не произойдёт
	if _, exist := cs.storage[message.Key]; !exist {
		if len(cs.storage) == cs.capacity {
			cs.deleteOldestMessage()
			cs.addMessage(message)
		} else {
			cs.addMessage(message)
		}
	}
	cs.mu.Unlock()
}

func (cs *ChatStorage) Last() *Message {
	var message *storageMessage

	cs.mu.Lock()
	for _, m := range cs.storage {
		if message == nil {
			message = m
			continue
		}
		if m.message.UnixTimestamp > message.message.UnixTimestamp {
			message = m
			continue
		}
		if m.message.UnixTimestamp == message.message.UnixTimestamp && m.createdAt > message.createdAt {
			message = m
		}
	}
	cs.mu.Unlock()

	return message.message
}

func (cs *ChatStorage) Delete(key MessageKey) {
	cs.mu.Lock()
	delete(cs.storage, key)
	cs.mu.Unlock()
}
