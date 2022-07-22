package models

type WelcomeEvent struct {
	Name string
}

func (e *WelcomeEvent) Type() EventType {
	return EventType("WELCOME")
}

func NewWelcomeKernel() *LambdaKernel {
	return &LambdaKernel{
		Handler: func(k *LambdaKernel, e Event) {
			if e, ok := e.(*JoinEvent); ok {
				k.ParticipantChannels[e.Participant.ID] <- &WelcomeEvent{Name: e.Participant.Name}
			}
		},
	}
}
