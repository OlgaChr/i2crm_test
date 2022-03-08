package messages

import (
	"fmt"
	"sync"
	"testing"
)

type TestUpsertCase struct {
	Key     string
	Message *Message
	Result  bool
}

func TestUpsert(t *testing.T) {
	storage := NewChatStorage(3)

	first := &Message{
		Key:           "first",
		UnixTimestamp: 1,
	}
	second := &Message{
		Key:           "second",
		UnixTimestamp: 2,
	}
	third := &Message{
		Key:           "third",
		UnixTimestamp: 3,
	}

	wg := &sync.WaitGroup{}
	wg.Add(4)
	// добавлены три сообщения
	go func(wg *sync.WaitGroup, storage *ChatStorage, second *Message) {
		storage.Upsert(second)
		wg.Done()
	}(wg, storage, second)
	go func(wg *sync.WaitGroup, storage *ChatStorage, first *Message) {
		storage.Upsert(first)
		wg.Done()
	}(wg, storage, first)
	go func(wg *sync.WaitGroup, storage *ChatStorage, third *Message) {
		storage.Upsert(third)
		wg.Done()
	}(wg, storage, third)

	go func(wg *sync.WaitGroup, storage *ChatStorage) {
		// после добавления должно удалиться первое сообщение, т.к. переполнение хранилища
		storage.Upsert(&Message{
			Key:           "fourth",
			UnixTimestamp: 3,
		})
		wg.Done()
	}(wg, storage)

	wg.Wait()

	if len(storage.storage) != 3 {
		fmt.Errorf("Error insert elements. Expected len = 3, get %d", len(storage.storage))
	}
	_, exist := storage.storage["first"]
	if exist {
		fmt.Errorf("Error insert 4th message. Expected deleted 'first', but first exist")
	}

	// это не должно добавиться - если сообщение уже существует, то ничего не должно произойти
	storage.Upsert(&Message{
		Key:           "second",
		UnixTimestamp: 4,
	})
	s := storage.storage["second"]
	if s.message != second {
		fmt.Errorf("Error insert duplicate message")
	}
}

func TestLast(t *testing.T) {
	storage := NewChatStorage(3)

	first := &Message{
		Key:           "first",
		UnixTimestamp: 1,
	}
	second := &Message{
		Key:           "second",
		UnixTimestamp: 2,
	}
	third := &Message{
		Key:           "third",
		UnixTimestamp: 2,
	}
	fourth := &Message{
		Key:           "fourth",
		UnixTimestamp: 3,
	}
	fifth := &Message{
		Key:           "fifth",
		UnixTimestamp: 2,
	}

	// добавлены три сообщения
	storage.Upsert(first)
	storage.Upsert(second)
	storage.Upsert(third)

	last := storage.Last()
	if last != third {
		fmt.Errorf("Wrong last element. Expected 'third'")
	}

	// 4ое новее остальных
	storage.Upsert(fourth)

	last = storage.Last()
	if last != fourth {
		fmt.Errorf("Wrong last element. Expected 'fourth'")
	}

	// 5ое старше 4ого
	storage.Upsert(fifth)

	last = storage.Last()
	if last != fourth {
		fmt.Errorf("Wrong last element. Expected 'fourth'")
	}
}

func TestDelete(t *testing.T) {
	storage := NewChatStorage(3)

	first := &Message{
		Key:           "first",
		UnixTimestamp: 1,
	}
	second := &Message{
		Key:           "second",
		UnixTimestamp: 2,
	}
	third := &Message{
		Key:           "third",
		UnixTimestamp: 3,
	}
	storage.Upsert(first)
	storage.Upsert(second)
	storage.Upsert(third)

	storage.Delete("second")
	storage.Delete("fourth")

	_, exist := storage.storage["second"]
	if exist {
		fmt.Errorf("Error delete element")
	}
}
