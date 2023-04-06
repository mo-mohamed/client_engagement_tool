package aws

type sms struct {
	mdn  string
	body string
}

func NewSms(mdn string, body string) *sms {
	return &sms{
		mdn:  mdn,
		body: body,
	}
}

func (s *sms) Process() error {

	// simulate work
	return nil
}
