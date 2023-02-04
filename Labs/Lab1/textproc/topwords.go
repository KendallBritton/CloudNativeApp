// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

package textproc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// Function to find the top K words in a document
func TopWords(path string, K int) []WordCount {

	wordList := make([]WordCount, 0)      		// Splice of words with its associated counts
	wordMap := make(map[string]WordCount) 		// Map to detect if words are found within the text

	file, err := os.Open(path) 					// Opens the file specified

	checkError(err) 				// Checks for error

	defer file.Close() 				// Cleans up the file once finished, i.e closes file

	scanDoc := bufio.NewScanner(file) 		// Takes in contents of whole file
	scanDoc.Split(bufio.ScanWords)    		// Seperates the text file word by word

	for scanDoc.Scan() { 			// Loops through whole text file until last word

		scannedWord := scanDoc.Text() 		// Takes in scanned word from document

		_, inListTF := wordMap[scannedWord] 		// Checks if scanned word is already in map

		if inListTF == false { 		// If not found in map, initialize it

			var temp WordCount
			temp.Count = 1
			temp.Word = scannedWord
			wordMap[scannedWord] = temp

		} else { 					// If found in map, update map key

			update := wordMap[scannedWord]
			update.Count += 1
			update.Word = scannedWord
			wordMap[scannedWord] = update

		}

	}

	for _, v := range wordMap { 			// Assign map contents to splice to be sorted

		wordList = append(wordList, v)

	}

	sortWordCounts(wordList) 			// Sorts to find the top K words with associated count values

	for i := 0; i < K; i++ { 			// Prints out the top K words with associated count values

		fmt.Print(wordList[i].String() + " ")

	}

	fmt.Print("\n")

	return wordList[:K]
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
