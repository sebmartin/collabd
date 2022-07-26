package models

// A test state that simply runs a lambda handler when executed
type LambdaState struct {
	Events         []Event
	PlayerChannels map[uint]chan<- ServerEvent
	Handler        func(*LambdaState, PlayerEvent)
}

func (k *LambdaState) Run(c <-chan PlayerEvent) {
	if k.PlayerChannels == nil {
		k.PlayerChannels = make(map[uint]chan<- ServerEvent)
	}

	for {
		event := <-c
		k.Events = append(k.Events, event)
		if event.Type() == JoinEventType {
			event := event.(*JoinEvent)
			k.PlayerChannels[event.Sender().ID] = event.Channel
		}
		if k.Handler != nil {
			k.Handler(k, event)
		}
	}
}

// A test state that sends a welcome event to players when they join
type WelcomeEvent struct {
	Name string
}

func (e *WelcomeEvent) Type() EventType {
	return EventType("WELCOME")
}

func NewWelcomeKernel() *LambdaState {
	return &LambdaState{
		Handler: func(k *LambdaState, e PlayerEvent) {
			if e, ok := e.(*JoinEvent); ok {
				k.PlayerChannels[e.Sender().ID] <- &WelcomeEvent{Name: e.Sender().Name}
			}
		},
	}
}
