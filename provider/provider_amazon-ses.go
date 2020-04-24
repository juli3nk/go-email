package provider

import (
	"fmt"

	"github.com/juli3nk/go-email/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func init() {
	RegisterDriver("amazon-ses", NewAmazonSESProvider)
}

type AmazonSESProvider struct {
	Config map[string]string
}

func NewAmazonSESProvider(config map[string]string) (Provider, error) {
	return &AmazonSESProvider{Config: config}, nil
}

func (p *AmazonSESProvider) Send(req *types.SendRequest) error {
	charset := "UTF-8"

	// start a new aws session
	sess, err := session.NewSession()
	if err != nil {
		fmt.Errorf("failed to create session,", err)
	}

	// Create an SES session
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(req.Body.Text),
				},
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(req.Body.Html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(req.Subject),
			},
		},
		Source: aws.String(req.From.Address),
	}

	/*
	if len(req.GetFrom().GetName()) > 0 {
		input.ConfigurationSetName = aws.String(req.GetFrom().GetName())
	}
	*/

	if len(req.From.ReplyTo) > 0 {
		input.ReplyToAddresses = []*string{
			aws.String(req.From.ReplyTo),
		}
	}

	var toAddrs []*string
	var ccAddrs []*string

	for _, a := range req.ToAddresses {
		toAddrs = append(toAddrs, aws.String(a))
	}

	if len(req.CcAddresses) > 0 {
		for _, a := range req.CcAddresses {
			ccAddrs = append(ccAddrs, aws.String(a))
		}
	}

	destination := &ses.Destination{
		ToAddresses: toAddrs,
		CcAddresses: ccAddrs,
	}

	input.Destination = destination

	// Attempt to send the email.
	_, err = svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		}

		return err
	}

	return nil
}
