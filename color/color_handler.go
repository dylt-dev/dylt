package color

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/dylt-dev/seq"
)

type ColorHandler struct {
	options ColorOptions
	w io.Writer
	mutex *sync.Mutex
	meta metaarglist
}

func NewColorHandler(w io.Writer, options ColorOptions) *ColorHandler {
	return &ColorHandler{w: w, options: options, mutex: &sync.Mutex{}, meta: []any{}}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.options.Level == nil {
		return true
	}

	return level >= h.options.Level.Level()
}

func (h *ColorHandler) Handle(ctx context.Context, rec slog.Record) error {
	var attrMap map[string]string
	var err error

	var meta metaarglist = make(metaarglist, len(h.meta) + rec.NumAttrs())
	copy(meta, h.meta)
	rec.Attrs(func (a slog.Attr) bool {
			meta = append(h.meta, a)
			return true
	})

	attrMap, err = createAttrMap(meta)
	if err != nil { return err }
	var sOut Styledstring
	if len(attrMap) == 0 {
		sOut = Styledstring(fmt.Sprintf("%s\n", rec.Message))
	} else {
		sMap := fmt.Sprintf("%#v", attrMap)
		sOut = Styledstring(fmt.Sprintf("%s: %s\n", rec.Message, sMap))
	}
	
	if rec.Level == slog.LevelDebug {
		sOut = Styledstring(sOut).Fg(X11.Gray50)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()
	_, err = h.w.Write([]byte(sOut))

	return err
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hNew := h
	for _, attr := range attrs {
		hNew.meta = append(h.meta, attr)
	}
	return hNew
}

func (h *ColorHandler) WithGroup(grp string) slog.Handler {
	hNew := h
	hNew.meta = append(h.meta, groupName(grp))
	return hNew
}

type ColorOptions struct {
	Level slog.Leveler
}

func createAttrMap (l metaarglist) (map[string]string, error) {
	var attr slog.Attr
	var attrMap = map[string]string{}
	var val string
	var gn groupName
	var groupNames = []groupName{}
	var sq seq.Seq[any] = newArraySeq(l)
	var arg any
	var err error
	var is bool
	var key string

	arg, err = sq.Next()
	// fmt.Printf("arg=%v type(arg)=%s err=%s\n", arg, reflect.TypeOf(arg), err)
	for true {
		// Loop  + add group names until a non-groupname is found	
		for true {
			if errors.Is(err, io.EOF) { break }
			if err != nil { return nil, err }

			gn, is = arg.(groupName)
			// fmt.Printf("gn=%s is=%t\n", gn, is)
			if !is { break }
			groupNames = append(groupNames, gn)

			arg, err = sq.Next()
			// fmt.Printf("arg=%v type(arg)=%s err=%s\n", arg, reflect.TypeOf(arg), err)
		}

		// Loop + add qualified attrs to map until a non-attr is found
		for true {
			if errors.Is(err, io.EOF) { break }
			if err != nil { return nil, err }
			
			attr, is = arg.(slog.Attr)
			// fmt.Printf("attr=%s is=%t\n", attr, is)
			if !is { break }
			
			key = fmt.Sprintf("%s.%s", join(groupNames...), attr.Key)
			val = attr.Value.String()
			// fmt.Printf("key=%s val=%s", key, val)
			attrMap[key] = val
		
			arg, err = sq.Next()
			// fmt.Printf("arg=%v type(arg)=%s err=%s\n", arg, reflect.TypeOf(arg), err)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, err
			}
		}
	}

	return attrMap, nil
}