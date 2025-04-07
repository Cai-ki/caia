package actor

type Envelope struct {
	Sender Ref
	Msg    interface{}
}
