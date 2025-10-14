package event

import "context"

type DeadLetterResolver struct{}

func NewDeadLetterResolver() *DeadLetterResolver {
	return &DeadLetterResolver{}
}

func (dlr *DeadLetterResolver) Handle(ctx context.Context, event *Event) error {
}
