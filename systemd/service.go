package systemd

import (
	"github.com/dylt-dev/dylt/template"
)
type Data struct {
	Name string
	Description string
	DisplayName string
	Username string
}

type Host struct {}

type Service struct {
	host Host
}

type Spec struct {
	Data Data
	Template template.Template
}