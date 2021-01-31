package cli

import (
	"fmt"
	"index/suffixarray"
	"regexp"
	"sort"
)

// SearchingCommand interface to describe a command that performs a search operation
type SearchingCommand interface {
  GetSearchParams() SearchParameters
}

type SearchParameters struct {
	Search      string
	Replacement string
	Mode        KeyValueMode
  doRegexp bool
}

// Match structure to keep indices of matched terms
type Match struct {
	path  string
	term  string
	key   string
	value string
	// sorted slices of indices of match starts and length
	keyIndex   [][]int
	valueIndex [][]int
}

type Searcher struct {
	cmd    SearchingCommand
	regexp *regexp.Regexp
}

// NewSearcher creates a new Searcher container for performing search and optionally replace
func NewSearcher(cmd SearchingCommand) (*Searcher, error) {
	var re *regexp.Regexp
	var err error
	params := cmd.GetSearchParams()

	if params.doRegexp {
		re, err = regexp.Compile(params.Search)
		if err != nil {
			return nil, fmt.Errorf("cannot parse regex pattern")
		}
	}

	return &Searcher{cmd: cmd, regexp: re}, nil
}

// IsMode returns true if the specified mode is enabled
func (s *Searcher) IsMode(mode KeyValueMode) bool {
	return s.cmd.GetSearchParams().Mode&mode == mode
}

func (s *Searcher) DoMatch(path string, k string, v string) (m []*Match) {
	if s.regexp != nil {
		return s.regexpMatch(path, k, v)
	}
	return s.substrMatch(path, k, v)
}

// find all indices for matches in key and value
func (s *Searcher) substrMatch(path string, k string, v string) (m []*Match) {
	substrLength := len(s.cmd.GetSearchParams().Search)
	keyMatchPairs := make([][]int, 0)
	if s.IsMode(ModeKeys) {
		keyIndex := suffixarray.New([]byte(k))
		keyMatches := keyIndex.Lookup([]byte(s.cmd.GetSearchParams().Search), -1)
		sort.Ints(keyMatches)
		for _, offset := range keyMatches {
			keyMatchPairs = append(keyMatchPairs, []int{offset, substrLength})
		}
	}

	valueMatchPairs := make([][]int, 0)
	if s.IsMode(ModeValues) {
		valueIndex := suffixarray.New([]byte(v))
		valueMatches := valueIndex.Lookup([]byte(s.cmd.GetSearchParams().Search), -1)
		sort.Ints(valueMatches)
		for _, offset := range valueMatches {
			valueMatchPairs = append(valueMatchPairs, []int{offset, substrLength})
		}
	}

	if len(keyMatchPairs) > 0 || len(valueMatchPairs) > 0 {
		m = []*Match{
			{
				path:       path,
				term:       s.cmd.GetSearchParams().Search,
				key:        k,
				value:      v,
				keyIndex:   keyMatchPairs,
				valueIndex: valueMatchPairs,
			},
		}
	}
	return m
}

func (s *Searcher) regexpMatch(path string, k string, v string) (m []*Match) {
	keyMatches := make([][]int, 0)
	if s.IsMode(ModeKeys) {
		keyMatches = s.regexp.FindAllIndex([]byte(k), -1)
	}
	valueMatches := make([][]int, 0)
	if s.IsMode(ModeValues) {
		valueMatches = s.regexp.FindAllIndex([]byte(v), -1)
	}

	if len(keyMatches) > 0 || len(valueMatches) > 0 {
		m = []*Match{
			{
				path:       path,
				term:       s.regexp.String(),
				key:        k,
				value:      v,
				keyIndex:   keyMatches,
				valueIndex: valueMatches,
			},
		}
	}
	return m
}
