package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// List of file extensions that I've encountered.
// Some of them aren't eBooks, but they were returned
// in previous search results.
var fileTypes = [...]string{
	"epub",
	"mobi",
	"azw3",
	"html",
	"rtf",
	"pdf",
	"cdr",
	"lit",
	"cbr",
	"doc",
	"htm",
	"jpg",
	"txt",
	"rar", // Compressed extensions should always be last 2 items
	"zip",
}

// BookDetail contains the details of a single Book found on the IRC server
type BookDetail struct {
	Server string `json:"server"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Format string `json:"format"`
	Size   string `json:"size"`
	Full   string `json:"full"`
}

type ParseError struct {
	Line  string `json:"line"`
	Error error  `json:"error"`
}

func (p *ParseError) MarshalJSON() ([]byte, error) {
	item := struct {
		Line  string `json:"line"`
		Error string `json:"error"`
	}{
		Line:  p.Line,
		Error: p.Error.Error(),
	}
	return json.Marshal(item)
}

func (p ParseError) String() string {
	return fmt.Sprintf("Error: %s. Line: %s.", p.Error, p.Line)
}

// ParseSearchFile converts a single search file into an array of BookDetail
func ParseSearchFile(filePath string) ([]BookDetail, []ParseError, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	books, errs := ParseSearchV2(file)
	return books, errs, nil
}

func ParseSearch(reader io.Reader) ([]BookDetail, []ParseError) {
	var books []BookDetail
	var parseErrors []ParseError

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "!") {
			dat, err := parseLine(line)
			if err != nil {
				parseErrors = append(parseErrors, ParseError{Line: line, Error: err})
			} else {
				books = append(books, dat)
			}
		}
	}

	books = normalizeBookDetailFields(books)
	sort.Slice(books, func(i, j int) bool { return books[i].Server < books[j].Server })

	return books, parseErrors
}

// Parse line extracts data from a single line
func parseLine(line string) (BookDetail, error) {

	//First check if it follows the correct format. Some servers don't include file info...
	if !strings.Contains(line, "::INFO::") {
		return BookDetail{}, errors.New("invalid line format. ::INFO:: not found")
	}

	var book BookDetail
	book.Full = line[:strings.Index(line, " ::INFO:: ")]
	var tmp int

	// Get Server
	if tmp = strings.Index(line, " "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse server")
	}
	book.Server = line[1:tmp] // Skip the "!"
	line = line[tmp+1:]

	// Get the Author
	if tmp = strings.Index(line, " - "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse author")
	}
	book.Author = line[:tmp]
	line = line[tmp+len(" - "):]

	// Get the Title
	for _, ext := range fileTypes { //Loop through each possible file extension we've got on record
		tmp = strings.Index(line, "."+ext) // check if it contains our extension
		if tmp == -1 {
			continue
		}
		book.Format = ext
		if ext == "rar" || ext == "zip" { // If the extension is .rar or .zip the actual format is contained in ()
			for _, ext2 := range fileTypes[:len(fileTypes)-2] { // Range over the eBook formats (exclude archives)
				if strings.Contains(line[:tmp], ext2) {
					book.Format = ext2
				}
			}
		}
		book.Title = line[:tmp]
		line = line[tmp+len(ext)+1:]
	}

	if book.Title == "" { // Got through the entire loop without finding a single match
		return BookDetail{}, errors.New("could not parse title")
	}

	// Get the Size
	if tmp = strings.Index(line, "::INFO:: "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse size")
	}

	line = strings.TrimSpace(line)
	splits := strings.Split(line, " ")

	if len(splits) >= 2 {
		book.Size = splits[1]
	}

	return book, nil
}

func ParseSearchV2(reader io.Reader) ([]BookDetail, []ParseError) {
	books := make([]BookDetail, 0)
	parseErrors := make([]ParseError, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "!") {
			dat, err := parseLineV2(line)
			if err != nil {
				parseErrors = append(parseErrors, ParseError{Line: line, Error: err})
			} else {
				books = append(books, dat)
			}
		}
	}

	books = normalizeBookDetailFields(books)
	sort.Slice(books, func(i, j int) bool { return books[i].Server < books[j].Server })

	return books, parseErrors
}

var (
	nonAuthorTokenPattern = regexp.MustCompile(`\d|[_:/\\|@]`)
	initialPattern        = regexp.MustCompile(`\b[A-Z]\.`)
	tokenSplitPattern     = regexp.MustCompile(`[^a-z0-9']+`)
)

var tokenStopWords = map[string]struct{}{
	"": {}, "a": {}, "an": {}, "and": {}, "the": {}, "of": {}, "in": {}, "on": {}, "to": {}, "for": {}, "from": {}, "with": {}, "by": {},
	"book": {}, "books": {}, "series": {}, "vol": {}, "volume": {}, "edition": {}, "ed": {},
}

func normalizeBookDetailFields(books []BookDetail) []BookDetail {
	if len(books) == 0 {
		return books
	}

	normalized := make([]BookDetail, len(books))
	copy(normalized, books)
	orientations := make([]bool, len(normalized))

	for i := range normalized {
		book := &normalized[i]

		baseScore := orientationScore(book.Author, book.Title)
		swappedScore := orientationScore(book.Title, book.Author)

		if swappedScore >= baseScore+4 {
			orientations[i] = true
		}
	}

	for iteration := 0; iteration < 4; iteration++ {
		authorFreq, titleFreq := buildColumnTokenFrequencies(normalized, orientations)
		changed := false

		for i := range normalized {
			book := normalized[i]

			authorKeep, titleKeep := book.Author, book.Title
			authorSwap, titleSwap := book.Title, book.Author

			keepScore := orientationScore(authorKeep, titleKeep)
			keepScore += tokenSimilarityScore(authorKeep, authorFreq) + tokenSimilarityScore(titleKeep, titleFreq)

			swapScore := orientationScore(authorSwap, titleSwap)
			swapScore += tokenSimilarityScore(authorSwap, authorFreq) + tokenSimilarityScore(titleSwap, titleFreq)

			shouldSwap := swapScore > keepScore
			if shouldSwap != orientations[i] {
				orientations[i] = shouldSwap
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	for i := range normalized {
		if orientations[i] {
			normalized[i].Author, normalized[i].Title = normalized[i].Title, normalized[i].Author
		}
	}

	return normalized
}

func buildColumnTokenFrequencies(books []BookDetail, orientations []bool) (map[string]int, map[string]int) {
	authorFreq := make(map[string]int)
	titleFreq := make(map[string]int)

	for i, book := range books {
		author, title := book.Author, book.Title
		if orientations[i] {
			author, title = title, author
		}

		for _, token := range scoreTokens(author) {
			authorFreq[token]++
		}
		for _, token := range scoreTokens(title) {
			titleFreq[token]++
		}
	}

	return authorFreq, titleFreq
}

func tokenSimilarityScore(value string, columnFreq map[string]int) int {
	score := 0
	for _, token := range scoreTokens(value) {
		freq := columnFreq[token]
		if freq <= 1 {
			continue
		}
		score += freq
	}
	return score
}

func scoreTokens(value string) []string {
	lower := strings.ToLower(strings.TrimSpace(value))
	if lower == "" {
		return nil
	}

	parts := tokenSplitPattern.Split(lower, -1)
	tokens := make([]string, 0, len(parts)*2)
	for _, part := range parts {
		if len(part) < 2 {
			continue
		}
		if _, isStopWord := tokenStopWords[part]; isStopWord {
			continue
		}
		tokens = append(tokens, part)
	}

	originalLen := len(tokens)
	for i := 0; i < originalLen-1; i++ {
		tokens = append(tokens, tokens[i]+"_"+tokens[i+1])
	}

	return tokens
}

func orientationScore(author string, title string) int {
	return authorLikelihood(author) + titleLikelihood(title)
}

func authorLikelihood(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return -10
	}

	words := strings.Fields(value)
	wordCount := len(words)
	score := 0

	switch {
	case wordCount <= 4:
		score += 3
	case wordCount <= 7:
		score += 1
	case wordCount >= 10:
		score -= 4
	}

	if strings.Contains(value, ",") {
		score += 2
	}

	if initialPattern.MatchString(value) {
		score += 1
	}

	if nonAuthorTokenPattern.MatchString(value) {
		score -= 3
	}

	lower := strings.ToLower(value)
	for _, token := range []string{" the ", " and ", " of ", " in ", " with ", " for ", " to ", " from ", " volume ", " edition "} {
		if strings.Contains(" "+lower+" ", token) {
			score -= 1
		}
	}

	if strings.ContainsAny(value, "[]()!?") {
		score -= 1
	}

	return score
}

func titleLikelihood(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return -10
	}

	words := strings.Fields(value)
	wordCount := len(words)
	score := 0

	if wordCount >= 3 {
		score += 2
	}
	if wordCount >= 7 {
		score += 1
	}

	lower := strings.ToLower(value)
	for _, token := range []string{" the ", " and ", " of ", " in ", " with ", " for ", " to ", " from ", " volume ", " edition "} {
		if strings.Contains(" "+lower+" ", token) {
			score += 1
		}
	}

	if strings.ContainsAny(value, "[]()!?:") {
		score += 1
	}

	if nonAuthorTokenPattern.MatchString(value) {
		score += 1
	}

	if strings.Contains(value, ",") && !strings.ContainsAny(lower, "0123456789") {
		score -= 1
	}

	return score
}

func parseLineV2(line string) (BookDetail, error) {
	getServer := func(line string) (string, error) {
		if line[0] != '!' {
			return "", errors.New("result lines must start with '!'")
		}

		firstSpace := strings.Index(line, " ")
		if firstSpace == -1 {
			return "", errors.New("unable parse server name")
		}

		return line[1:firstSpace], nil
	}

	getAuthor := func(line string) (string, error) {
		firstSpace := strings.Index(line, " ")
		dashChar := strings.Index(line, " - ")
		if dashChar == -1 {
			return "", errors.New("unable to parse author")
		}
		author := line[firstSpace+len(" ") : dashChar]

		// Handles case with weird author characters %\w% ("%F77FE9FF1CCD% Michael Haag")
		if strings.Contains(author, "%") {
			split := strings.SplitAfterN(author, " ", 2)
			return split[1], nil
		}

		return author, nil
	}

	getTitle := func(line string) (string, string, int) {
		title := ""
		fileFormat := ""
		endIndex := -1
		// Get the Title
		for _, ext := range fileTypes { //Loop through each possible file extension we've got on record
			endTitle := strings.Index(line, "."+ext) // check if it contains our extension
			if endTitle == -1 {
				continue
			}
			fileFormat = ext
			if ext == "rar" || ext == "zip" { // If the extension is .rar or .zip the actual format is contained in ()
				for _, ext2 := range fileTypes[:len(fileTypes)-2] { // Range over the eBook formats (exclude archives)
					if strings.Contains(strings.ToLower(line[:endTitle]), ext2) {
						fileFormat = ext2
					}
				}
			}
			startIndex := strings.Index(line, " - ")
			title = line[startIndex+len(" - ") : endTitle]
			endIndex = endTitle
		}

		return title, fileFormat, endIndex
	}

	getSize := func(line string) (string, int) {
		const delimiter = " ::INFO:: "
		infoIndex := strings.LastIndex(line, delimiter)

		if infoIndex != -1 {
			// Handle cases when there is additional info after the file size (ex ::HASH:: )
			parts := strings.Split(line[infoIndex+len(delimiter):], " ")
			return parts[0], infoIndex
		}

		return "N/A", len(line)
	}

	server, err := getServer(line)
	if err != nil {
		return BookDetail{}, err
	}

	author, err := getAuthor(line)
	if err != nil {
		return BookDetail{}, err
	}

	title, format, titleIndex := getTitle(line)
	if titleIndex == -1 {
		return BookDetail{}, errors.New("unable to parse title")
	}

	size, endIndex := getSize(line)

	return BookDetail{
		Server: server,
		Author: author,
		Title:  title,
		Format: format,
		Size:   size,
		Full:   strings.TrimSpace(line[:endIndex]),
	}, nil
}
