package conventionalcommitparser

import (
	"regexp"
	"strings"
)

type Header struct {
	raw       string
	Type      string
	Scope     string
	Subject   string
	Important bool
}

func (h *Header) String() string {
	return h.raw
}

var (
	headerPattern       = regexp.MustCompile(`^(?i)([\s\w-]*)(\((.*)\))?(!?):\s+(.*)$`)
	revertHeaderPattern = regexp.MustCompile(`^(?i)revert\s(.*)$`)
)

func parseHeader(txt string) Header {
	headerMatchers := headerPattern.FindStringSubmatch(txt)
	revertHeaderMatchers := revertHeaderPattern.FindStringSubmatch(txt)
	header := Header{raw: txt}

	if len(headerMatchers) != 0 { // conventional commit
		header.Type = strings.TrimSpace(strings.ToLower(headerMatchers[1]))
		header.Scope = strings.TrimSpace(headerMatchers[3])
		header.Important = headerMatchers[4] == "!"
		header.Subject = headerMatchers[5]
	} else if len(revertHeaderMatchers) != 0 { // revert commit
		subject := strings.Trim(revertHeaderMatchers[1], "\"")
		subject = strings.Trim(subject, "'")
		header.Type = "revert"
		header.Subject = subject
	} else { // commom commit
		header.Type = ""
		header.Scope = ""
		header.Subject = txt
	}

	return header
}
