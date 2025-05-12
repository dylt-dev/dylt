package color

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
)

// There are a ton of rules to properly handling attributes for slog Handlers.
// The rules are rules. They apply to all slog Handlers.
// But slog doesn't provide any support for following the rules at all.
// It just makes rules.
//
// SAD.
//
// It would be nice to have some help.
//
// On to the rules
// - logger.WithGroup() adds a group to the handler
// - logger.With() adds one or more attributes to the handler
// - slog.With() adds one or more attributes to the default logger (! asymmetric; no WithGroup())

type groupName string
type metaarg interface { slog.Attr | groupName }
type metaarglist []any

func join[T ~string] (l ...T) string {
	return strings.Join(toStrings(l...), ".")
}

func toStrings[T ~string] (l ...T) []string {
	if l == nil { return nil}
	var ss = make([]string, 0, len(l))
	for _, el := range l {
		var s string = string(el)
		ss = append(ss, s)
	}

	return ss
}

// Lotta compromises here
// - It would be nice to make this a method on metaarglist. But go does not allow
//   type constraings on methods.
// - It would be nice if this function were variadic. But generics don't work that way.
//   If it were variadic all of the arguments would need to be of the same type. That type
//   would need to a member of metaarg, but groupnames and attrs could not be included in the same
//   call, which is a problem. 
func addMetaarg[T metaarg] (l metaarglist, arg T) metaarglist {
	l = append(l, arg)
	return l
}

func addMetaargs (l metaarglist, args ...any) metaarglist {
	for _, arg := range args {
		switch m := arg.(type) {
		case groupName: l = addMetaarg(l, m)
		case slog.Attr: l = addMetaarg(l, m)
		default: panic(fmt.Sprintf("invalid metaarg type: %s", reflect.TypeOf(m)))
		}
	}
	return l
}

type arraySeq[T any] struct {
	a []T
	iCur int
} 

func newArraySeq[T any] (a []T) *arraySeq[T] {
	return &arraySeq[T]{a: a, iCur: 0}
}

func (sq *arraySeq[T]) Next () (T, error) {
	if sq.iCur >= len(sq.a) {
		return *new(T), io.EOF
	}
	var el T = sq.a[sq.iCur]
	sq.iCur++
	return el, nil
}