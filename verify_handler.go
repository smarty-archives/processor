package processor

type VerifyHandler struct {
	in          chan *Envelope
	out         chan *Envelope
	application Verifier
}

type Verifier interface {
	Verify(AddressInput)
}

func NewVerifyHandler(in, out chan *Envelope, application Verifier) *VerifyHandler {
	return &VerifyHandler{
		in:          in,
		out:         out,
		application: application,
	}
}

func (this *VerifyHandler) Handle() {
	received := <-this.in

	this.application.Verify(received.Input)

	this.out <- received
}
