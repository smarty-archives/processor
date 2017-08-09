package processor

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
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
	buffer  *ReadWriteSpyBuffer
	writer  *csv.Writer
}

func (this *WriterHandlerFixture) Setup() {
	this.buffer = NewReadWriteSpyBuffer("")
	this.input = make(chan *Envelope, 10)
	this.handler = NewWriterHandler(this.input, this.buffer)
}

var recordMatchingHeader = AddressOutput{
	Status:        "Status",
	DeliveryLine1: "DeliveryLine1",
	LastLine:      "LastLine",
	City:          "City",
	State:         "State",
	ZIPCode:       "ZIPCode",
}

func (this *WriterHandlerFixture) TestHeaderMatchesRecords() {
	this.input <- &Envelope{Output: recordMatchingHeader}
	close(this.input)

	this.handler.Handle()

	this.assertHeaderMatchesRecord()
}
func (this *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := this.outputLines()
	header := lines[0]
	record := lines[1]

	this.So(header, should.Equal, "Status,DeliveryLine1,LastLine,City,State,ZIPCode")
	this.So(record, should.Equal, header)
}

func (this *WriterHandlerFixture) TestAllEnvelopesWritten() {
	this.sendEnvelopes(2)
	this.handler.Handle()

	if lines := this.outputLines(); this.So(lines, should.HaveLength, 3) {
		this.So(lines[1], should.Equal, "A1,B1,C1,D1,E1,F1")
		this.So(lines[2], should.Equal, "A2,B2,C2,D2,E2,F2")
	}
}

func (this *WriterHandlerFixture) TestOutputClosed() {
	close(this.input)
	this.handler.Handle()

	this.So(this.buffer.closed, should.Equal, 1)
}

func (this *WriterHandlerFixture) sendEnvelopes(count int) {
	for x := 1; x < count+1; x++ {
		this.input <- &Envelope{
			Output: createOutput(strconv.Itoa(x)),
		}
	}
	close(this.input)
}
func createOutput(index string) AddressOutput {
	return AddressOutput{
		Status:        "A" + index,
		DeliveryLine1: "B" + index,
		LastLine:      "C" + index,
		City:          "D" + index,
		State:         "E" + index,
		ZIPCode:       "F" + index,
	}
}

func (this *WriterHandlerFixture) outputLines() []string {
	outputFile := strings.TrimSpace(this.buffer.String())
	return strings.Split(outputFile, "\n")
}

func (this *WriterHandlerFixture) TestEndOfFileEnvelopeC() {

}

//////////////////////////////////////////////////////////

type ReadWriteSpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewReadWriteSpyBuffer(value string) *ReadWriteSpyBuffer {
	return &ReadWriteSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *ReadWriteSpyBuffer) Close() error {
	this.closed++
	return nil
}
