package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	CharSet = "UTF-8"
)

type EmailInput struct {
	Sender    string
	Recipient string
	Subject   string
	HtmlBody  string
	TextBody  string
}

func (c *awsConfig) SendEmail(input EmailInput) error {
	sess := c.GetSession()
	if sess == nil {
		return fmt.Errorf("AWS session not initialized")
	}

	svc := ses.New(sess)

	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(input.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(input.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(input.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(input.Subject),
			},
		},
		Source: aws.String(input.Sender),
	}

	_, err := svc.SendEmail(emailInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return fmt.Errorf("%s: %v", ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return fmt.Errorf("%s: %v", ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return fmt.Errorf("%s: %v", ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				return fmt.Errorf("unexpected error: %v", aerr.Error())
			}
		} else {
			return fmt.Errorf("error: %v", err.Error())
		}
	}

	return nil
}
