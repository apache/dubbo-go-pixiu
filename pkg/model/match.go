package model

// StringMatcher matcher string
type StringMatcher struct {
	Matcher MatcherType
}

// Match
func (sm *StringMatcher) Match() (bool, error) {
	return true, nil
}

// MatcherType matcher type
type MatcherType int32

const (
	Exact MatcherType = 0 + iota
	Prefix
	Suffix
	Regex
)

var MatcherTypeName = map[int32]string{
	0: "Exact",
	1: "Prefix",
	2: "Suffix",
	3: "Regex",
}

var MatcherTypeValue = map[string]int32{
	"Exact":  0,
	"Prefix": 1,
	"Suffix": 2,
	"Regex":  3,
}

type HeaderMatcher struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
	Regex bool   `yaml:"regex" json:"regex"`
}
