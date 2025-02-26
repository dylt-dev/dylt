package cmd

type Cmdline []string

func (o *Cmdline) Args () []string {
	return (*o)[1:]
}

func (o *Cmdline) Command () string {
	if len(*o) <= 0 {
		return ""
	}
	return (*o)[0]
}

func (o *Cmdline) HasCommand () bool {
	return len(*o) > 0
}

type Command interface {
	Run (args []string) error
}
