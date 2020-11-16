package web

import (
	"fmt"
	"net/url"
)

type Linker struct {
	BaseURL *url.URL
}

func NewLinker(base string) (*Linker, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	return &Linker{
		BaseURL: baseURL,
	}, nil
}

func (l *Linker) Origin() string {
	return l.BaseURL.Scheme + "://" + l.BaseURL.Host
}

func (l *Linker) Domain() string {
	return l.BaseURL.Hostname()
}

func (l *Linker) RPID() string {
	return l.Domain()
}

func (l *Linker) EnrollmentPage(token string) string {
	return l.BaseURL.String() + "#/enroll/" + token
}

func (l *Linker) LoginPage() string {
	return "#/login"
}
