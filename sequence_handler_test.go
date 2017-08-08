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
	envelope := &Envelope{Sequence: 0}
	this.input <- envelope
	close(this.input)
	this.handler.Handle()
	this.So(<-this.output, should.Equal, envelope)
}

func (this *SequenceHandlerFixture) TestEnvelopesReceivedOutOfOrder_BufferedUntilContiguousBlock() {
	this.input <- &Envelope{Sequence: 4}
	this.input <- &Envelope{Sequence: 3}
	this.input <- &Envelope{Sequence: 2}
	this.input <- &Envelope{Sequence: 1}
	this.input <- &Envelope{Sequence: 0}
	close(this.input)

	this.handler.Handle()

	this.So((<-this.output).Sequence, should.Equal, 0)
	this.So((<-this.output).Sequence, should.Equal, 1)
	this.So((<-this.output).Sequence, should.Equal, 2)
	this.So((<-this.output).Sequence, should.Equal, 3)
	this.So((<-this.output).Sequence, should.Equal, 4)
}
