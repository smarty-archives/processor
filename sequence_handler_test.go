package processor

import (
	"testing"

	"github.com/smartystreets/assertions/should"
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
	this.handler = NewSequenceHandler(this.input, this.output)
}

func (this *SequenceHandlerFixture) TestExpectedEnvelopeSentToOutput() {
	this.sendEnvelopesInSequence(0, 1, 2, 3, 4)

	this.handler.Handle()

	this.So(this.sequenceOrder(), should.Resemble, []int{0, 1, 2, 3, 4})
	this.So(this.handler.buffer, should.BeEmpty)
}

func (this *SequenceHandlerFixture) TestEnvelopesReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	this.sendEnvelopesInSequence(4, 2, 0, 3, 1)

	this.handler.Handle()

	this.So(this.sequenceOrder(), should.Resemble, []int{0, 1, 2, 3, 4})
	this.So(this.handler.buffer, should.BeEmpty)
}
func (this *SequenceHandlerFixture) sendEnvelopesInSequence(sequences ...int) {
	for _, sequence := range sequences {
		this.input <- &Envelope{Sequence: sequence}
	}
	close(this.input)
}


func (this *SequenceHandlerFixture) sequenceOrder() (order []int) {
	close(this.output)

	for envelope := range this.output {
		order = append(order, envelope.Sequence)
	}
	return order
}
