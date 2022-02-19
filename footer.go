package conventionalcommitparser

import (
	"regexp"
	"strings"
)

type Footer struct {
	Tag     string
	Title   string
	Content string
}

var (
	footerTagPattern            = regexp.MustCompile(`(?i)^([a-z]+(-[a-z]+)*):\s?(.*)$`)
	footerHashPattern           = regexp.MustCompile(`^(?i)^([\w\-]+)\s+(#.*)`)
	footerBreakingChangePattern = regexp.MustCompile(`^(BREAKING\sCHANGE):\s*(.*)$`)
)

func paseFooterParagraph(txt string) Footer {
	footer := Footer{}

	tagMatcher := footerTagPattern.FindStringSubmatch(txt)
	breakingChangeMatcher := footerBreakingChangePattern.FindStringSubmatch(txt)
	hashTagMatcher := footerHashPattern.FindStringSubmatch(txt)

	if len(breakingChangeMatcher) != 0 {
		footer.Tag = strings.TrimSpace(breakingChangeMatcher[1])
		footer.Title = strings.TrimSpace(breakingChangeMatcher[2])
	} else if len(tagMatcher) != 0 {
		footer.Tag = strings.TrimSpace(tagMatcher[1])
		footer.Title = strings.TrimSpace(tagMatcher[3])
	} else if len(hashTagMatcher) != 0 {
		footer.Tag = strings.TrimSpace(hashTagMatcher[1])
		footer.Title = strings.TrimSpace(hashTagMatcher[2])
	} else {
		footer.Tag = ""
		footer.Title = txt
	}

	return footer
}

func isFooterParagraph(txt string) bool {
	if footerBreakingChangePattern.MatchString(txt) {
		return true
	}

	if footerTagPattern.MatchString(txt) {
		return true
	}

	if footerHashPattern.MatchString(txt) {
		return true
	}

	return false
}

func parseFooter(txt string) Footer {
	lines := splitToLines(txt)

	footer := Footer{}

	contents := make([]string, 0)

lineLoop:
	for index, line := range lines {
		if index == 0 {
			footer = paseFooterParagraph(line)
			continue lineLoop
		} else {
			contents = append(contents, line)
		}
	}

	footer.Content = strings.TrimSpace(strings.Join(contents, "\n"))

	return footer
}
