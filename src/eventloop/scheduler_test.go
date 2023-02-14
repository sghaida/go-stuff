package eventloop_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

)

func Test_Scheduler(t *testing.T) {
	t.Parallel()
	el := eventloop.InitializeEventLoop()
	el.Start()

	t.Run("emit couple of events", func(t *testing.T) {
		channels := make([]chan eventloop.EventReply, 0, 100)
		lock := new(sync.RWMutex)
		wg := new(sync.WaitGroup)
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func(index int) {
				ch := make(chan eventloop.EventReply, 1)
				lock.Lock()
				channels = append(channels, ch)
				lock.Unlock()
				event := eventloop.Event{
					RequestID: uuid.New().String(),
					UserID:    uuid.New().String(),
					EventType: eventloop.ServiceRequestSubmission,
					Channel:   ch,
				}

				el.Emmit(event, func() eventloop.EventReply {
					return eventloop.EventReply{
						Payload: event.UserID,
						Error:   nil,
					}
				})
				go func(ev eventloop.Event) {
					t := time.Duration(rand.Intn(100))
					time.Sleep(t * time.Millisecond)
					result := <-ev.Channel
					log.Infof(result.Payload.(string))
					wg.Done()
				}(event)
			}(i)
		}
		wg.Wait()
	})

	t.Run("emit another couple of events", func(t *testing.T) {
		channels := make([]chan eventloop.EventReply, 0, 100)
		lock := new(sync.RWMutex)
		wg := new(sync.WaitGroup)
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func(index int) {
				ch := make(chan eventloop.EventReply, 1)
				lock.Lock()
				channels = append(channels, ch)
				lock.Unlock()
				event := eventloop.Event{
					RequestID: uuid.New().String(),
					UserID:    uuid.New().String(),
					EventType: eventloop.ServiceRequestSubmission,
					Channel:   ch,
				}

				el.Emmit(event, func() eventloop.EventReply {
					return eventloop.EventReply{
						Payload: event.UserID,
						Error:   nil,
					}
				})
				go func(ev eventloop.Event) {
					t := time.Duration(rand.Intn(100))
					time.Sleep(t * time.Millisecond)
					result := <-ev.Channel
					log.Infof(result.Payload.(string))
					wg.Done()
				}(event)
			}(i)
		}
		wg.Wait()
	})
}
