package gungan

import (
	"fmt"
	"strings"
)

func (jj *JarJar) Spake(englishText string) string {
	return jj.translateText(jj.normalizeText(englishText))
}

func (jj *JarJar) normalizeText(raw string) string {
	var normalizedText string
	if jj.typoToleranceEnabled {
		updated, _ := jj.typoChecker.Replace(raw)
		normalizedText = updated
	}
	return normalizedText
}

func (jj *JarJar) translateText(englishText string) string {
	translatedText := englishText

	var firstTierMatches []string
	for pattern, gunganeseTerm := range jj.dictionary.firstTier {
		if pattern.MatchString(englishText) {
			firstTierMatches = append(firstTierMatches, gunganeseTerm)
			translatedText = pattern.ReplaceAllString(translatedText, gunganeseTerm)
		}
	}

	for pattern, gunganeseTerm := range jj.dictionary.secondTier {
		for _, match := range firstTierMatches {
			if pattern.MatchString(match) {
				continue
			}
		}
		translatedText = pattern.ReplaceAllString(translatedText, gunganeseTerm)
	}

	splitWords := strings.Split(translatedText, " ")

	for i, word := range splitWords {
		processedWord := word
		for suffix, replacementSuffix := range jj.dictionary.suffixes {
			if strings.HasSuffix(word, suffix) {
				processedWord = fmt.Sprintf("%s%s", strings.TrimSuffix(word, suffix), replacementSuffix)
				break
			}
		}
		splitWords[i] = processedWord
	}

	return strings.Join(splitWords, " ")
}
