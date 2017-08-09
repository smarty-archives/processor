package processor

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture

	handler *WriterHandler
	input   chan *Envelope
	buffer  *WriterSpyBuffer
	writer  *csv.Writer
}

func (this *WriterHandlerFixture) Setup() {
	this.buffer = NewWriterSpyBuffer("")
	this.input = make(chan *Envelope, 10)
	this.handler = NewWriterHandler(this.input, this.buffer)
}

func (this *WriterHandlerFixture) TestHeaderWritten() {
	this.handler.Handle()

	this.So(this.buffer.String(), should.Equal, "Status,DeliveryLine1,City,State,ZIPCode\n")
}

func (this *WriterHandlerFixture) TestOutputClosed() {
	this.handler.Handle()

	this.So(this.buffer.closed, should.Equal, 1)
}

func (this *WriterHandlerFixture) TestEnvelopeWritten() {
	this.input <- &Envelope{AddressOutput{
		Status: "Status",
		
	}}
	close(this.input)
}

//////////////////////////////////////////////////////////

type WriterSpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewWriterSpyBuffer(value string) *WriterSpyBuffer {
	return &WriterSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *WriterSpyBuffer) Close() error {
	this.closed++
	return nil
}
