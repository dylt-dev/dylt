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


func (ks KeyString) ChildName(prefix KeyString) (string, bool) {
	afterPrefix, is := ks.CutPrefix(prefix)
	if !is {
		return "", false
	}

	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		return "", false
	}

	return segments[0], true
}


func (ks KeyString) CutPrefix(prefix KeyString) (KeyString, bool) {
	snew, is := strings.CutPrefix(string(ks), string(prefix))
	return KeyString(snew), is
}


func (ks KeyString) ElementName(prefix KeyString) string {
	afterPrefix, is := ks.CutPrefix(prefix)
	if !is {
		return ""
	}
	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		return ""
	}
	return segments[len(segments)-1]
}


func (ks KeyString) HasPrefix (prefix KeyString) bool {
	return strings.HasPrefix(string(ks), string(prefix))
}

func (ks KeyString) Index() (int, bool) {
	if string(ks) == "" {
		return 0, false
	}
	segments := ks.Segments()
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


func (ks KeyString) Segments() []string {
	var segments []string = []string{}
	var iSlashes []int = []int{}

	// Make an array of all slash locations
	s2 := ks.WithStartSlash().WithEndSlash()
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

func (ks KeyString) WithEndSlash() KeyString {
	if ks == "" {
		return "/"
	}

	var sb strings.Builder
	sb.WriteString(string(ks))
	if ks[len(ks)-1] != '/' {
		sb.WriteRune('/')
	}

	return KeyString(sb.String())
}

func (ks KeyString) WithoutEndSlash() KeyString {
	if len(ks) == 0 {
		return ks
	}

	if ks[len(ks)-1] == '/' {
		return ks[:len(ks)-1]
	}

	return ks
}


func (ks KeyString) WithStartSlash() KeyString {
	if ks == "" {
		return "/"
	}

	var sb strings.Builder
	if ks[0] != '/' {
		sb.WriteRune('/')
	}
	sb.WriteString(string(ks))

	return KeyString(sb.String())
}

func (ks KeyString) WithoutStartSlash() KeyString {
	if len(ks) == 0 {
		return ks
	}

	if ks[0] == '/' {
		return ks[1:]
	}

	return ks
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