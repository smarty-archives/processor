package processor

import (
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

func (this *ReaderHandlerFixture) TestCSVRecordSentInEnvelope() {
	this.writeLine("A,B,C,D")

	this.reader.Handle()

	this.So(<-this.output, should.Resemble, &Envelope{Input: AddressInput{
		Street1: "A",
		City:    "B",
		State:   "C",
		ZIPCode: "D",
	}})
}

func (this *ReaderHandlerFixture) TestAllCSVRecordsWrittenToOutput() {
	this.writeLine("A1,B1,C1,D1")
	this.writeLine("A2,B2,C2,D2")

	this.reader.Handle()

	this.So(<-this.output, should.Resemble, &Envelope{Input: AddressInput{
		Street1: "A1",
		City:    "B1",
		State:   "C1",
		ZIPCode: "D1",
	}})

	this.So(<-this.output, should.Resemble, &Envelope{Input: AddressInput{
		Street1: "A2",
		City:    "B2",
		State:   "C2",
		ZIPCode: "D2",
	}})

	this.So(<-this.output, should.Resemble, &Envelope{EOF: true})
	this.So(<-this.output, should.BeNil)
	this.So(this.buffer.closed, should.Equal, 1)
}

func (this *ReaderHandlerFixture) writeLine(line string) {
	this.buffer.WriteString(line + "\n")
}
