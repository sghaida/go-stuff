package eventloop

import (
	"errors"
	"os"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// this event loop implementation executes events FIFO bases
// to run this. below is the example
// 			el := eventloop.InitializeEventLoop()
// 			el.Start()
// to emmit events
//   		event := eventloop.Event{
//					RequestID:	uuid.New().String(),
//					UserID:    uuid.New().String(),
//					EventType: eventloop.ServiceRequestSubmission,
//					Channel:   ch,
//				}
//
//				el.Emmit(event, func() eventloop.EventReply {
//					return eventloop.EventReply{
//						Payload: event.UserID,
//						Error:   nil,
//					}
//				})

type EventID = string

type EventType = int

// holds the event types that is used internally
const (
	RequireFeedback EventType = iota
	NoFeedback
)

// Event which to be emitted into the queue
type Event struct {
	EventType EventType
	Channel   chan EventReply
	eventID   EventID
}

// EventReply  is the event response which is returned after executing the event
type EventReply struct {
	Payload interface{}
	Error   error
}
type Action func() EventReply

// Queue queue data structure that holds the events meta data and the events functions to be executed
type Queue struct {
	events  []Event
	actions []Action
	mutex   *sync.RWMutex
}

type baseEvent struct {
	e  Event
	f  func() EventReply
	ex error
}

var (
	pushChan = make(chan baseEvent, 100)
	popChan  = make(chan baseEvent, 100)
	pushWG   = new(sync.WaitGroup)
	pullWG   = new(sync.WaitGroup)
)

// InitializeEventLoop creates the event loop
func InitializeEventLoop() *Queue {
	return &Queue{
		events:  make([]Event, 0),
		actions: make([]Action, 0),
		mutex:   new(sync.RWMutex),
	}
}

// emmit event into the queue, uses channel for synchronization
func (q *Queue) Emmit(e Event, f func() EventReply) EventID {
	eventID := uuid.New().String()
	e.eventID = eventID
	event := baseEvent{e, f, nil}
	pushChan <- event
	log.Infof("EventID: %v has been emitted", e.eventID)
	return eventID
}

// pop get an event out of the queue
// if the queue is empty return error
func (q *Queue) pop() {
	// lock
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.events) == 0 {
		popChan <- baseEvent{
			ex: errors.New("empty queue"),
		}
		return
	}
	// get related event
	event := q.events[len(q.events)-1]
	q.events = q.events[:len(q.events)-1]
	// get related action
	action := q.actions[len(q.actions)-1]
	q.actions = q.actions[:len(q.actions)-1]

	popChan <- baseEvent{e: event, f: action}
	return
}

// push push event into the queue to be executed
func (q *Queue) push(event baseEvent) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.events = append(q.events, event.e)
	q.actions = append(q.actions, event.f)
}

// start starts the event loop and executes or emmit some events to be executed
func (q *Queue) Start() {
	var internal func()
	internal = func() {
		select {
		case event := <-pushChan:
			pushWG.Add(1)
			go func() {
				q.push(event)
				pushWG.Done()
			}()
			pushWG.Wait()
		default:
			if len(q.events) == 0 {
				return
			}
			// pop one item
			q.pop()
			ev := <-popChan
			if ev.ex != nil || ev.e.Channel == nil {
				return
			}
			pullWG.Add(1)
			go func(channel chan<- EventReply) {
				defer close(ev.e.Channel)
				ev.e.Channel <- ev.f()
				pullWG.Done()
			}(ev.e.Channel)
			pullWG.Wait()
		}
	}
	// start the main el on its own goroutine
	go func() {
		for {
			internal()
		}
	}()
}

func (q *Queue) Stop() {
	defer close(pushChan)
	defer close(popChan)
	q.events = nil
	q.actions = nil
	os.Exit(0)
}
