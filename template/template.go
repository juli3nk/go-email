package template

import (
	"bytes"
	templateH "html/template"
	templateT "text/template"
)

type Template struct {
	Subject  string
	BodyText string
	BodyHtml string
}

type Config struct {
	Template *Template
}

func New(msgTemplate *Template) (*Config, error) {
	c := Config{
		Template: msgTemplate,
	}

	return &c, nil
}

func (c *Config) GenerateTemplate(dSubject, dBody map[string]string) (*Template, error) {
	// Subject
	tSubject := templateT.Must(templateT.New("emailSubject").Parse(c.Template.Subject))

	bufSubject := new(bytes.Buffer)

	if err := tSubject.Execute(bufSubject, dSubject); err != nil {
		return nil, err
	}

	// Body text
	tBodyT := templateT.Must(templateT.New("emailBodyText").Parse(c.Template.BodyText))

	bufBodyT := new(bytes.Buffer)

	if err := tBodyT.Execute(bufBodyT, dBody); err != nil {
		return nil, err
	}

	// Body HTML
	tBodyH := templateH.Must(templateH.New("emailBodyHtml").Parse(c.Template.BodyHtml))

	bufBodyH := new(bytes.Buffer)

	if err := tBodyH.Execute(bufBodyH, dBody); err != nil {
		return nil, err
	}

	return &Template{
		Subject:  bufSubject.String(),
		BodyText: bufBodyT.String(),
		BodyHtml: bufBodyH.String(),
	}, nil
}
