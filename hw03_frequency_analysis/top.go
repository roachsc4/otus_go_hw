package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

// preprocessText - cleans text from \t and \n symbols for easier splitting.
func preprocessText(text string) string {
	return strings.ReplaceAll(text, "\n\t", " ")
}

// preprocessWord - normalizes word to lower case and deletes redundant characters.
func preprocessWord(word string) string {
	word = strings.ToLower(word)
	word = strings.Trim(word, " ,.:;\"")
	return word
}

type WordCounterMap map[string]int

// Count occurrences for all words and place it to map.
func CalculateWordCounterMap(words []string) WordCounterMap {
	wordCountMap := make(WordCounterMap)

	for _, word := range words {
		if word == "" {
			continue
		}
		word = preprocessWord(word)
		if !wordRegexp.MatchString(word) {
			continue
		}
		wordCountMap[word]++
	}

	return wordCountMap
}

// Get slice of Top-10 words sorted by occurrences rate.
func GetTop10WordsFromWordCounterMap(wordCountMap WordCounterMap) []string {
	top10Words := make([]string, 0, len(wordCountMap))
	for word := range wordCountMap {
		top10Words = append(top10Words, word)
	}
	sort.Slice(top10Words, func(i, j int) bool {
		return wordCountMap[top10Words[i]] > wordCountMap[top10Words[j]]
	})

	if len(top10Words) > 10 {
		return top10Words[:10]
	}

	return top10Words
}

// Accepted word expression looks like "abc-dfg", where "-" is optional.
var wordRegexp = regexp.MustCompile("[a-zA-Zа-яА-Я]+(-[a-zA-Zа-яА-Я]+)?")

// Top10 - function for finding Top-10 most frequent words in a text.
// Initially text goes through preprocessing, word-normalizing, counting and filling Top-10 words slice.
func Top10(text string) []string {
	if text == "" {
		return nil
	}
	// Preprocessing whole text
	text = preprocessText(text)
	// Splitting text to words and initializing helper variables
	wordCountMap := CalculateWordCounterMap(strings.Split(text, " "))
	return GetTop10WordsFromWordCounterMap(wordCountMap)
}
