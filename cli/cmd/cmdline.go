package cmd

type Cmdline []string

func (o Cmdline) Args () Cmdline {
	if len(o) < 1 {
		return nil
	}
	var args Cmdline = (o)[1:]
	return args
}

func (o Cmdline) Command () string {
	if len(o) <= 0 {
		return ""
	}
	return (o)[0]
}

func (o Cmdline) HasCommand () bool {
	return len(o) > 0
}

type Command interface {
	Run () error
}
