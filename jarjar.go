package gungan

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/client9/misspell"
)

const matchingPattern = `(?i)\b%s\b`

type lexicon struct {
	firstTier  map[*regexp.Regexp]string
	secondTier map[*regexp.Regexp]string
	suffixes   map[string]string
}

type JarJar struct {
	typoToleranceEnabled bool
	typoChecker          *misspell.Replacer
	dictionary           lexicon
}

func NewJarJar(enableTypoTolerance bool) (*JarJar, error) {
	misspellReplacer := misspell.New()

	rawLex := struct {
		FirstTier  map[string]string `json:"first_tier"`
		SecondTier map[string]string `json:"second_tier"`
		Suffixes   map[string]string `json:"suffixes"`
	}{}
	if err := json.Unmarshal([]byte(englishGunganeseLexicon), &rawLex); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gunganese dictionary > %w", err)
	}

	dict := lexicon{
		firstTier:  map[*regexp.Regexp]string{},
		secondTier: map[*regexp.Regexp]string{},
		suffixes:   rawLex.Suffixes,
	}

	for k, v := range rawLex.FirstTier {
		dict.firstTier[regexp.MustCompile(fmt.Sprintf(matchingPattern, k))] = v
	}

	for k, v := range rawLex.SecondTier {
		dict.secondTier[regexp.MustCompile(fmt.Sprintf(matchingPattern, k))] = v
	}

	return &JarJar{
		typoToleranceEnabled: enableTypoTolerance,
		typoChecker:          misspellReplacer,
		dictionary:           dict,
	}, nil
}
