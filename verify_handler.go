package processor

type VerifyHandler struct {
	in  chan interface{}
	out chan interface{}
}

func NewVerifyHandler(in, out chan interface{}) *VerifyHandler {
	return &VerifyHandler{
		in:  in,
		out: out,
	}
}

func (this *VerifyHandler) Listen() {
	this.out <- 1
}
