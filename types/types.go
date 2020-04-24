package types

type From struct {
	Name    string
	Address string
	ReplyTo string
}

type Body struct {
	Text string
	Html string
}

type SendRequest struct {
	From        *From
	ToAddresses []string
	CcAddresses []string
	Subject     string
	Body        *Body
}
