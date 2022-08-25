package models

// A test stage that simply runs a lambda handler when executed
type LambdaStage struct {
	Events         []Event
	PlayerChannels map[uint]chan<- ServerEvent
	Handler        func(*LambdaStage, PlayerEventEnvelope)
}

func (stage *LambdaStage) Run(c <-chan PlayerEventEnvelope) StageRunner {
	if stage.PlayerChannels == nil {
		stage.PlayerChannels = make(map[uint]chan<- ServerEvent)
	}

	for {
		event, ok := <-c
		if !ok {
			return nil
		}
		stage.Events = append(stage.Events, event.PlayerEvent)
		if event.Type() == JoinEventType {
			event := event.PlayerEvent.(*JoinEvent)
			stage.PlayerChannels[event.Sender().ID] = event.Channel
		}
		if stage.Handler != nil {
			stage.Handler(stage, event)
		}
	}
}

type WelcomeEvent struct {
	Name string
}

func (e *WelcomeEvent) Type() EventType {
	return EventType("WELCOME")
}

// A test stage that sends a welcome event to players when they join
func NewWelcomeStage() *LambdaStage {
	return &LambdaStage{
		Handler: func(k *LambdaStage, e PlayerEventEnvelope) {
			if e, ok := e.PlayerEvent.(*JoinEvent); ok {
				k.PlayerChannels[e.Sender().ID] <- &WelcomeEvent{Name: e.Sender().Name}
			}
		},
	}
}
