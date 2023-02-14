package eventloop_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sghaida/go-stuff/src/eventloop"
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
					EventType: eventloop.RequireFeedback,
					Channel:   ch,
				}

				el.Emmit(event, func() eventloop.EventReply {
					return eventloop.EventReply{
						Payload: fmt.Sprintf("event - %d", index),
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
					// no feedback is not being checked for the time being
					EventType: eventloop.NoFeedback,
					Channel:   ch,
				}

				el.Emmit(event, func() eventloop.EventReply {
					return eventloop.EventReply{
						Payload: struct {
							Name string
						}{fmt.Sprintf("name: %d", index)},
						Error: nil,
					}
				})
				go func(ev eventloop.Event) {
					t := time.Duration(rand.Intn(100))
					time.Sleep(t * time.Millisecond)
					result := <-ev.Channel
					log.Info(result.Payload.(struct{ Name string }))
					wg.Done()
				}(event)
			}(i)
		}
		wg.Wait()
	})
}
