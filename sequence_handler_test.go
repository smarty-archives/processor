package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSequenceHandler(t *testing.T) {
	gunit.Run(new(SequenceHandlerFixture), t)
}

type SequenceHandlerFixture struct {
	*gunit.Fixture

	input   chan *Envelope
	output  chan *Envelope
	handler *SequenceHandler
}

func (this *SequenceHandlerFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	var handler *SequenceHandler = NewSequenceHandler(this.input, this.output)
}
func NewSequenceHandler(input chan *Envelope, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{}
}
