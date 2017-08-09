package processor

import (
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWireupFixture(t *testing.T) {
	gunit.Run(new(WireupFixture), t)
}

type WireupFixture struct {
	*gunit.Fixture

	reader  *ReadWriteSpyBuffer
	writer  *ReadWriteSpyBuffer
	client  *FakeHTTPClient
	handler Handler
}

func (this *WireupFixture) Setup() {
	this.reader = NewReadWriteSpyBuffer("")
	this.writer = NewReadWriteSpyBuffer("")
	this.client = &FakeHTTPClient{}
	this.handler = Configure(this.reader, this.writer, this.client).Build()
}

func (this *WireupFixture) Test() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	this.reader.WriteString("Street1,City,State,ZIPCode")
	this.reader.WriteString("A,B,C,D")
	this.reader.WriteString("A,B,C,D")

	this.handler.Handle()

	this.So(this.writer.String(), should.Equal, "")
}
