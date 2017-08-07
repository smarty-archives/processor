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

	input       chan *Envelope
	output      chan *Envelope
	application *FakeVerifier
	handler     *VerifyHandler
}

func (this *HandlerFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	this.application = NewFakeVerifier()
	this.handler = NewVerifyHandler(this.input, this.output, this.application)
}

func (this *HandlerFixture) TestVerifierReceivesInput() {
	envelope := &Envelope{
		Input: AddressInput{
			Street1: "42",
		},
	}
	this.input <- envelope
	close(this.input)

	this.handler.Handle()

	this.AssertEqual(envelope, <-this.output)
	this.AssertEqual(envelope.Input, this.application.input)
}

/////////////////////////////////////////////////////////////////

type FakeVerifier struct {
	input AddressInput
}

func NewFakeVerifier() *FakeVerifier {
	return &FakeVerifier{}
}

func (this *FakeVerifier) Verify(value AddressInput) {
	this.input = value
}
