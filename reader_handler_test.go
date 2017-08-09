package processor

import (
	"strconv"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestReaderHandlerFixture(t *testing.T) {
	gunit.Run(new(ReaderHandlerFixture), t)
}

type ReaderHandlerFixture struct {
	*gunit.Fixture

	buffer *ReadWriteSpyBuffer
	output chan *Envelope
	reader *ReaderHandler
}

func (this *ReaderHandlerFixture) Setup() {
	this.buffer = NewReadWriteSpyBuffer("")
	this.output = make(chan *Envelope, 10)
	this.reader = NewReaderHandler(this.buffer, this.output)

	const header = "Street1,City,State,ZIPCode"
	this.writeLine(header)
}

func (this *ReaderHandlerFixture) TestAllCSVRecordsSentToOutput() {
	this.writeLine("A1,B1,C1,D1")
	this.writeLine("A2,B2,C2,D2")

	this.reader.Handle()

	this.assertRecordsSent()
	this.assertCleanup()
}
func (this *ReaderHandlerFixture) assertRecordsSent() {
	this.So(<-this.output, should.Resemble, buildEnvelope(initialSequenceValue))
	this.So(<-this.output, should.Resemble, buildEnvelope(initialSequenceValue+1))
}
func (this *ReaderHandlerFixture) assertCleanup() {
	this.So(<-this.output, should.Resemble, &Envelope{Sequence: initialSequenceValue+2, EOF: true})
	this.So(<-this.output, should.BeNil)
	this.So(this.buffer.closed, should.Equal, 1)
}

func (this *ReaderHandlerFixture) writeLine(line string) {
	this.buffer.WriteString(line + "\n")
}

func buildEnvelope(index int) *Envelope {
	suffix := strconv.Itoa(index + 1)

	return &Envelope{
		Sequence: index,
		Input: AddressInput{
			Street1: "A" + suffix,
			City:    "B" + suffix,
			State:   "C" + suffix,
			ZIPCode: "D" + suffix,
		},
	}
}

func (this *ReaderHandlerFixture) TestMalformedInputReturnsError() {
	malformedRecord := "A1"
	this.writeLine(malformedRecord)

	err := this.reader.Handle()

	if this.So(err, should.NotBeNil) {
		this.So(err.Error(), should.Equal, "Malformed input")
	}
}
