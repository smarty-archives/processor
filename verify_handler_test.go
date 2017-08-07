package processor

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HandlerFixture), t)
}

type HandlerFixture struct {
	*gunit.Fixture

	input   chan interface{}
	output  chan interface{}
	handler *VerifyHandler
}

func (this *HandlerFixture) Setup() {
	this.input = make(chan interface{}, 10)
	this.output = make(chan interface{}, 10)
	this.handler = NewVerifyHandler(this.input, this.output)
}

func (this *HandlerFixture) TestVerifierReceivesInput() {
	this.input <- 1
	close(this.input)

	this.handler.Listen()

	this.AssertEqual(1, <-this.output)
}

