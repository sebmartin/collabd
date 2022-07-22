package models

// A game kernel that simply runs a lambda handler when executed. This is mostly useful
// for writing unit test fixtures.
type LambdaKernel struct {
	Events              []Event
	ParticipantChannels map[uint]chan Event
	Handler             func(*LambdaKernel, Event)
}

func (k *LambdaKernel) Run(c chan Event) {
	if k.ParticipantChannels == nil {
		k.ParticipantChannels = make(map[uint]chan Event)
	}

	for {
		event := <-c
		k.Events = append(k.Events, event)
		if event.Type() == JoinEventType {
			event := event.(*JoinEvent)
			k.ParticipantChannels[event.Participant.ID] = event.Channel
		}
		if k.Handler != nil {
			k.Handler(k, event)
		}
	}
}
