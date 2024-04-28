package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {

}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher EventDispatcher
}

func (e *EventDispatcherTestSuite) SetupTest() {
	e.eventDispatcher = *NewEventDispatcher()
	e.handler = TestEventHandler{
		ID: 1,
	}
	e.handler2 = TestEventHandler{
		ID: 2,
	}
	e.handler3 = TestEventHandler{
		ID: 3,
	}
	e.event = TestEvent{Name: "Event", Payload: "Test1"}
	e.event2 = TestEvent{Name: "Event2", Payload: "Test2"}
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))

	err = e.eventDispatcher.Register(e.event.GetName(), &e.handler2)
	e.Nil(err)
	e.Equal(2, len(e.eventDispatcher.handlers[e.event.GetName()]))

	assert.Equal(e.T(), &e.handler, e.eventDispatcher.handlers[e.event.GetName()][0])
	assert.Equal(e.T(), &e.handler2, e.eventDispatcher.handlers[e.event.GetName()][1])
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameEvent() {
	err := e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))

	err = e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Equal(ErrHandlerAlreadyRegistered, err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// Event 1
	err := e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))

	err = e.eventDispatcher.Register(e.event.GetName(), &e.handler2)
	e.Nil(err)
	e.Equal(2, len(e.eventDispatcher.handlers[e.event.GetName()]))

	// Event 2
	err = e.eventDispatcher.Register(e.event2.GetName(), &e.handler3)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event2.GetName()]))

	e.eventDispatcher.Clear()
	e.Equal(0, len(e.eventDispatcher.handlers))
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	// Event 1
	err := e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))

	err = e.eventDispatcher.Register(e.event.GetName(), &e.handler2)
	e.Nil(err)
	e.Equal(2, len(e.eventDispatcher.handlers[e.event.GetName()]))

	assert.True(e.T(), e.eventDispatcher.Has(e.event.GetName(), &e.handler))
	assert.True(e.T(), e.eventDispatcher.Has(e.event.GetName(), &e.handler2))
	assert.False(e.T(), e.eventDispatcher.Has(e.event.GetName(), &e.handler3))
}

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Dispatcher() {
	eh := &MockEventHandler{}
	eh.On("Handle", &e.event)

	eh2 := &MockEventHandler{}

	eh2.On("Handle", &e.event)

	e.eventDispatcher.Register(e.event.GetName(), eh)
	e.eventDispatcher.Register(e.event.GetName(), eh2)

	e.eventDispatcher.Dispatch(&e.event)
	eh.AssertExpectations(e.T())
	eh2.AssertExpectations(e.T())
	eh.AssertNumberOfCalls(e.T(), "Handle", 1)
	eh2.AssertNumberOfCalls(e.T(), "Handle", 1)
}

func (e *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	// Event 1
	err := e.eventDispatcher.Register(e.event.GetName(), &e.handler)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))

	err = e.eventDispatcher.Register(e.event.GetName(), &e.handler2)
	e.Nil(err)
	e.Equal(2, len(e.eventDispatcher.handlers[e.event.GetName()]))

	// Event 2
	err = e.eventDispatcher.Register(e.event2.GetName(), &e.handler3)
	e.Nil(err)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event2.GetName()]))

	e.eventDispatcher.Remove(e.event.GetName(), &e.handler)
	e.Equal(1, len(e.eventDispatcher.handlers[e.event.GetName()]))
	assert.Equal(e.T(), &e.handler2, e.eventDispatcher.handlers[e.event.GetName()][0])

	e.eventDispatcher.Remove(e.event.GetName(), &e.handler2)
	e.Equal(0, len(e.eventDispatcher.handlers[e.event.GetName()]))

	e.eventDispatcher.Remove(e.event2.GetName(), &e.handler3)
	e.Equal(0, len(e.eventDispatcher.handlers[e.event2.GetName()]))

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
