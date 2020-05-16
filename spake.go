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

func (jj *JarJar) translateText(translationStr string) string {
	for rgx, entry := range jj.dictionary.terms {
		if len(entry.dependencies) > 0 && hasDependencyClash(translationStr, entry.dependencies) {
			continue
		}
		translationStr = rgx.ReplaceAllString(translationStr, fmt.Sprintf("${1}%s${3}", entry.gunganeseTerm))
	}

	splitWords := strings.Split(translationStr, " ")

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

func hasDependencyClash(translatedText string, dependencies []dictionaryEntry) bool {
	for _, dep := range dependencies {
		if dep.rgx.MatchString(translatedText) {
			return true
		}
	}
	return false
}
