package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

// preprocessText - cleans text from \t and \n symbols for easier splitting.
func preprocessText(text string) string {
	return strings.Replace(strings.Replace(text, "\t", " ", -1), "\n", " ", -1)
}

// preprocessWord - normalizes word to lower case and deletes redundant characters.
func preprocessWord(word string) string {
	word = strings.ToLower(word)
	word = strings.Trim(word, " ,.:;\"")
	return word
}

// Top10 - function for finding top-10 most frequent words in a text.
// Initially text goes through preprocessing, word-normalizing.
func Top10(text string) []string {
	if text == "" {
		return make([]string, 0)
	}
	// Accepted word expression looks like "abc-dfg", where "-" is optional
	wordRegexp := regexp.MustCompile("[a-zA-Zа-яА-Я]+(-[a-zA-Zа-яА-Я]+)?")

	// Preprocessing whole text
	text = preprocessText(text)

	// Splitting text to words and initializing helper variables
	splittedText := strings.Split(text, " ")
	top10Words := make([]string, 0, len(splittedText))
	wordCountMap := make(map[string]int)

	for _, word := range splittedText {
		if word == "" {
			continue
		}
		word = preprocessWord(word)
		if !wordRegexp.MatchString(word) {
			continue
		}
		_, wordExists := wordCountMap[word]
		// If words appears for the first time, append it to resulting slice
		if !wordExists {
			top10Words = append(top10Words, word)
		}
		// In both cases we should count this word
		wordCountMap[word]++
	}

	sort.Slice(top10Words, func(i, j int) bool {
		return wordCountMap[top10Words[i]] > wordCountMap[top10Words[j]]
	})

	if len(top10Words) > 10 {
		return top10Words[:10]
	}

	return top10Words
}
