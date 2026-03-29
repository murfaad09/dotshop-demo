package aws

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/harishash/dotshop-be/integration/aws/email_template"
)

type AWSService interface {
	CreateCollectionMail(recipientEmail, curatorName, collectionName string) error
	NewLookCreatedMail(recipientEmail, curatorName, lookName string) error
	LargeOrderPlacedMail(recipientEmail, curatorName, orderId string) error
	ForgotPasswordMail(recipientEmail, firstName, token string) error
}

func (c *awsConfig) CreateCollectionMail(recipientEmail, curatorName, collectionName string) error {
	tmpl, err := template.New("email").Parse(email_template.CreateCollectionTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	data := struct {
		CuratorName    string
		CollectionName string
	}{
		CuratorName:    curatorName,
		CollectionName: collectionName,
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	subject := fmt.Sprintf("New Collection Created: %s", collectionName)

	input := EmailInput{
		Sender:    "tech@dotshop.ai",
		Recipient: recipientEmail,
		Subject:   subject,
		HtmlBody:  body.String(),
		TextBody:  body.String(),
	}
	return c.SendEmail(input)
}

func (c *awsConfig) NewLookCreatedMail(recipientEmail, curatorName, lookName string) error {
	tmpl, err := template.New("email").Parse(email_template.CreateLookTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	data := struct {
		CuratorName string
		LookName    string
	}{
		CuratorName: curatorName,
		LookName:    lookName,
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	subject := fmt.Sprintf("New Look Created: %s", lookName)

	input := EmailInput{
		Sender:    "tech@dotshop.ai",
		Recipient: recipientEmail,
		Subject:   subject,
		HtmlBody:  body.String(),
		TextBody:  body.String(),
	}
	return c.SendEmail(input)
}

func (c *awsConfig) LargeOrderPlacedMail(recipientEmail, curatorName, orderId string) error {
	tmpl, err := template.New("email").Parse(email_template.LargeOrderTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	data := struct {
		CuratorName string
		OrderID     string
	}{
		CuratorName: curatorName,
		OrderID:     orderId,
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	subject := "Large order placed"

	input := EmailInput{
		Sender:    "tech@dotshop.ai",
		Recipient: recipientEmail,
		Subject:   subject,
		HtmlBody:  body.String(),
		TextBody:  body.String(),
	}
	return c.SendEmail(input)
}

func (c *awsConfig) ForgotPasswordMail(recipientEmail, firstName, token string) error {
	tmpl, err := template.New("email").Parse(email_template.ForgotPasswordTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	url := fmt.Sprintf("https://staging-dotshop.vercel.app/reset-password?token=%s", token)

	data := struct {
		FirstName string
		ResetURL  string
	}{
		FirstName: firstName,
		ResetURL:  url,
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	subject := "Password Reset Request"

	input := EmailInput{
		Sender:    "tech@dotshop.ai",
		Recipient: recipientEmail,
		Subject:   subject,
		HtmlBody:  body.String(),
		TextBody:  body.String(),
	}
	return c.SendEmail(input)
}
