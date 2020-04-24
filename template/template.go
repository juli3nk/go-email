package template

import (
	"bytes"
	"fmt"
	templateH "html/template"
	templateT "text/template"
)

type Template struct {
	Subject  string
	BodyText string
	BodyHtml string
}

type Config struct {
	Templates map[string]map[string]*Template
}

func New() (*Config, error) {
	c := Config{
		Templates: make(map[string]map[string]*Template),
	}

	return &c, nil
}

func (c *Config) RegisterTemplate(name, language string, msgTemplate *Template) error {
	if c.Templates[name] == nil {
		c.Templates[name] = make(map[string]*Template)
	}

	if _, exist := c.Templates[name][language]; exist {
		return fmt.Errorf("Language '%s' already exists for template name '%s'", language, name)
	}

	c.Templates[name][language] = msgTemplate

	return nil
}

func (c *Config) GenerateTemplate(name, language string, dSubject, dBody map[string]string) (*Template, error) {
	// Subject
	tSubject := templateT.Must(templateT.New("emailSubject").Parse(c.Templates[name][language].Subject))

	bufSubject := new(bytes.Buffer)

	if err := tSubject.Execute(bufSubject, dSubject); err != nil {
		return nil, err
	}

	// Body text
	tBodyT := templateT.Must(templateT.New("emailBodyText").Parse(c.Templates[name][language].BodyText))

	bufBodyT := new(bytes.Buffer)

	if err := tBodyT.Execute(bufBodyT, dBody); err != nil {
		return nil, err
	}

	// Body HTML
	tBodyH := templateH.Must(templateH.New("emailBodyHtml").Parse(c.Templates[name][language].BodyHtml))

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
