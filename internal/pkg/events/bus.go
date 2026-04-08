package events

import (
	"log"
	"sync"
	"time"
)

// EventType определяет тип события
type EventType string

const (
	UserDeleted EventType = "user.deleted"
)

// Event представляет событие в системе
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// Handler функция-обработчик события
type Handler func(event Event) error

// EventBus синхронная шина событий
type EventBus struct {
	handlers map[EventType][]Handler
	mu       sync.RWMutex
}

// NewEventBus создает новый EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]Handler),
	}
}

// Subscribe подписывает обработчик на событие
func (eb *EventBus) Subscribe(eventType EventType, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish публикует событие синхронно
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			log.Printf("event handler failed for %s: %v", event.Type, err)
		}
	}
}
