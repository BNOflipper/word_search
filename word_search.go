package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

var dictionary map[string]int

// Functions using ByCount implements the sort.Interface for the words occurence
type ByCount []string

func (w ByCount) Len() int           { return len(w) }
func (w ByCount) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w ByCount) Less(i, j int) bool { return dictionary[w[i]] < dictionary[w[j]] }

func removeCss(htmlCode string) string {
	return removeHtmlTag(htmlCode, "<style", "</style>")
}

func removeJavaScript(htmlCode string) string {
	return removeHtmlTag(htmlCode, "<script", "</script>")
}

func removeHtmlTag(htmlCode, startTag, endTag string) string {
	var newHtmlCode strings.Builder
	stringArray := strings.Split(htmlCode, endTag)

	for _, substring := range stringArray {
		substring, _, _ = strings.Cut(substring, startTag)
		newHtmlCode.WriteString(substring)
	}

	return newHtmlCode.String()
}

func getPlainText(htmlCode string) string {
	var plainText strings.Builder
	htmlCode = removeJavaScript(htmlCode)
	htmlCode = removeCss(htmlCode)
	stringArray := strings.Split(htmlCode, ">")

	for _, substring := range stringArray {
		substring, _, _ = strings.Cut(substring, "<")
		plainText.WriteString(substring)
		plainText.WriteString(" ")
	}

	return plainText.String()
}

func main() {
	if len(os.Args) == 1 {
		panic("No web page entered")
	}

	webPage := os.Args[1]

	resp, err := http.Get(webPage)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var htmlCode strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan(); i++ {
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		htmlCode.WriteString(scanner.Text())
	}

	words := strings.Split(getPlainText(htmlCode.String()), " ")
	dictionary = make(map[string]int)

	for _, word := range words {
		word = strings.Trim(word, ".,:;'?!()[]{}\"-")

		if len(word) > 0 {
			dictionary[word]++
		}
	}

	fmt.Println("Type a word to search for. Type \"-all\" to display all words.")

	for {
		var searchWord string
		fmt.Scanf("%s", &searchWord)

		if searchWord == "-all" {
			i := 0
			foundWords := make([]string, len(dictionary))

			for word := range dictionary {
				foundWords[i] = word
				i++
			}

			sort.Sort(ByCount(foundWords))

			for i := 0; i < len(foundWords); i++ {
				fmt.Printf("\"%s\" was found %d time(s)\n", foundWords[i], dictionary[foundWords[i]])
			}
		} else if len(searchWord) > 0 {
			fmt.Printf("\"%s\" was found %d time(s)\n", searchWord, dictionary[searchWord])
		}
	}
}
