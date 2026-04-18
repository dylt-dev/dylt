package eco

import "strings"

type KeyString string

func (s KeyString) CutPrefix(prefix string) (string, bool) {
	snew, is := strings.CutPrefix(string(s), prefix)
	return snew, is
}


func (s KeyString) ChildName(prefix string) string {
	afterPrefix, is := s.CutPrefix(prefix)
	if !is {
		return ""
	}

	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		return ""
	}

	return segments[0]
}

func (s KeyString) ElementName(prefix string) string {
	afterPrefix, is := s.CutPrefix(prefix)
	if !is {
		return ""
	}
	segments := KeyString(afterPrefix).Segments()
	if segments == nil || len(segments) == 0 {
		return ""
	}
	return segments[len(segments)-1]
}

func (s KeyString) Segments() []string {
	var segments []string = []string{}
	var iSlashes []int = []int{}

	// Make an array of all slash locations
	s2 := s.WithStartSlash().WithEndSlash()
	for i, c := range s2 {
		if c == '/' {
			iSlashes = append(iSlashes, i)
		}
	}

	for i := range len(iSlashes) - 1 {
		iStart := iSlashes[i] + 1
		iEnd := iSlashes[i+1]
		segment := s2[iStart:iEnd]
		segments = append(segments, string(segment))
	}

	return segments
}

func (s KeyString) WithEndSlash() KeyString {
	if s == "" {
		return "/"
	}

	var sb strings.Builder
	sb.WriteString(string(s))
	if s[len(s)-1] != '/' {
		sb.WriteRune('/')
	}

	return KeyString(sb.String())
}

func (s KeyString) WithoutEndSlash() KeyString {
	if len(s) == 0 {
		return s
	}

	if s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}

	return s
}


func (s KeyString) WithStartSlash() KeyString {
	if s == "" {
		return "/"
	}

	var sb strings.Builder
	if s[0] != '/' {
		sb.WriteRune('/')
	}
	sb.WriteString(string(s))

	return KeyString(sb.String())
}

func (s KeyString) WithoutStartSlash() KeyString {
	if len(s) == 0 {
		return s
	}

	if s[0] == '/' {
		return s[1:]
	}

	return s
}


func createKeyString (segments ...string) KeyString {
	sb := strings.Builder{}
	for i := range(len(segments)-1) {
		segment := segments[i]
		sb.WriteString(segment)
		sb.WriteString("/")
	}
	
	sb.WriteString(segments[len(segments)-1])
	s := sb.String()
	return KeyString(s)
}