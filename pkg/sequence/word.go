package sequence

import "strings"

func CountWordsLongerThan(words []string, length int) int {
	count := 0
	for _, word := range words {
		if len(word) > length {
			count++
		}
	}
	return count
}

func CountWordsWithRepeatingChars(words []string, minRepeat int) int {
	count := 0
	for _, word := range words {
		charCount := make(map[rune]int)
		for _, ch := range word {
			charCount[ch]++
		}
		for _, v := range charCount {
			if v >= minRepeat {
				count++
				break
			}
		}
	}
	return count
}

func CountWordsSameStartEnd(words []string) int {
	count := 0
	for _, word := range words {
		if len(word) >= 1 && word[0] == word[len(word)-1] {
			count++
		}
	}
	return count
}

func CapitalizeFirstLetter(words []string) []string {
	newWords := make([]string, len(words))
	for i, word := range words {
		if len(word) == 0 {
			newWords[i] = word
			continue
		}
		newWords[i] = strings.ToUpper(string(word[0])) + word[1:]
	}
	return newWords
}
