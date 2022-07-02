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

type Attachment struct {
	Data        []byte
	Filename    string
	ContentType string
	Disposition string
}

type SendRequest struct {
	From        *From
	ToAddresses []string
	CcAddresses []string
	Subject     string
	Body        *Body
	Attachments []*Attachment
}
