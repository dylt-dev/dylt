package eco

import (
	"strconv"
	"strings"
)

type KeyString string


func (ks KeyString) AddSegment(segment string) KeyString {
	sb := strings.Builder{}
	sb.WriteString(string(ks.WithEndSlash()))
	sb.WriteString(segment)
	return KeyString(sb.String())
}


func (ks KeyString) Child(prefix KeyString) (KeyString, bool) {
	childName, is := ks.ChildName(prefix)
	if !is {
		return "", false
	}
	return prefix.AddSegment(childName), true
}


func (s KeyString) ChildName(prefix KeyString) (string, bool) {
	afterPrefix, is := s.CutPrefix(string(prefix))
	if !is {
		return "", false
	}

	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		return "", false
	}

	return segments[0], true
}


func (s KeyString) CutPrefix(prefix string) (KeyString, bool) {
	snew, is := strings.CutPrefix(string(s), prefix)
	return KeyString(snew), is
}


func (s KeyString) ElementName(prefix string) string {
	afterPrefix, is := s.CutPrefix(prefix)
	if !is {
		return ""
	}
	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		return ""
	}
	return segments[len(segments)-1]
}

func (s KeyString) Index() (int, bool) {
	if string(s) == "" {
		return 0, false
	}
	segments := s.Segments()
	lastSeg := segments[len(segments)-1]
	index, err := strconv.Atoi(lastSeg)
	if err != nil {
		return 0, false
	}
	return index, true
}


func (ks KeyString) IsLeaf () bool {
	segments := ks.Segments()
	return len(segments) == 1
}


func (ks KeyString) IsParent(keyString KeyString) bool {
	s := string(ks.WithoutEndSlash())
	sKeyString := string(keyString.WithoutEndSlash())
	return strings.HasPrefix(sKeyString, s) && s != sKeyString
}


func (ks KeyString) PopHead () (string, KeyString) {
	s := string(ks.WithoutStartSlash().WithoutEndSlash())
	iFirstSlash := strings.Index(s, "/") 
	if iFirstSlash != -1 {
		head := s[:iFirstSlash]
		body := KeyString(s[iFirstSlash:]) 
		return head, body
	}
	
	return string(ks.WithoutStartSlash()), KeyString("")
}


func (ks KeyString) TrimHead () KeyString {
	s := string(ks.WithoutStartSlash().WithoutEndSlash())
	iFirstSlash := strings.Index(s, "/") 
	if iFirstSlash != -1 {
		return KeyString(s[iFirstSlash:])
	}

	return KeyString("")
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