package gungan

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/client9/misspell"
)

const matchingPattern = `(?i)(\s|^)(%s)(\s|$)`

type lexicon struct {
	terms    map[*regexp.Regexp]dictionaryEntry
	suffixes map[string]string
}

type dictionaryEntry struct {
	englishTerm   string
	gunganeseTerm string
	rgx           *regexp.Regexp
	dependencies  []dictionaryEntry
}

type JarJar struct {
	typoToleranceEnabled bool
	typoChecker          *misspell.Replacer
	dictionary           lexicon
}

func NewJarJar(enableTypoTolerance bool) (*JarJar, error) {
	misspellReplacer := misspell.New()

	rawLex := struct {
		Terms    map[string]string `json:"terms"`
		Suffixes map[string]string `json:"suffixes"`
	}{}
	if err := json.Unmarshal([]byte(englishGunganeseLexicon), &rawLex); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gunganese dictionary > %w", err)
	}

	dict := lexicon{
		terms:    map[*regexp.Regexp]dictionaryEntry{},
		suffixes: rawLex.Suffixes,
	}

	for englishTerm, gunganeseTerm := range rawLex.Terms {
		rgx := regexp.MustCompile(fmt.Sprintf(matchingPattern, englishTerm))
		entry := dictionaryEntry{
			englishTerm:   englishTerm,
			gunganeseTerm: gunganeseTerm,
			rgx:           rgx,
		}
		dict.terms[rgx] = entry
	}

	for rgx, associatedEntry := range dict.terms {
		englishTerm := associatedEntry.englishTerm
		for _, v := range dict.terms {
			if englishTerm != v.englishTerm && rgx.MatchString(v.englishTerm) {
				associatedEntry.dependencies = append(associatedEntry.dependencies, v)
			}
		}
		dict.terms[rgx] = associatedEntry
	}

	return &JarJar{
		typoToleranceEnabled: enableTypoTolerance,
		typoChecker:          misspellReplacer,
		dictionary:           dict,
	}, nil
}
