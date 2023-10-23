package opds

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/go-chi/chi/v5"
)

type OPDS struct {
	*chi.Mux
	libraryDirectory      string
	searchResultDirectory string
}

func NewOPDS(libraryDirectory, searchResultDirectory string) *OPDS {
	r := chi.NewRouter()
	opds := &OPDS{
		Mux:                   r,
		libraryDirectory:      libraryDirectory,
		searchResultDirectory: searchResultDirectory,
	}

	opds.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// spew.Dump(r.Header)
			next.ServeHTTP(w, r)
		})
	})

	opds.Get("/", rootHandler())
	opds.Get("/library", opds.libraryHandler())
	opds.Get("/library/*", opds.downloadBookHandler())
	opds.Get("/opensearch", opds.opensearchMetaHandler())
	opds.Get("/search", opds.searchHandler())
	opds.Get("/search-results", opds.searchResultsHandler())
	opds.Get("/search-results/*", opds.displaySearchResultsHandler())

	return opds
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		selfHref := "/opds"
		f := NewFeed("OpenBooks OPDS", "", selfHref)
		f.Link = append(f.Link, Link{
			Rel: FeedSearchLinkRel, Href: "/opds/opensearch", Type: FeedSearchLinkType, Title: "Search on catalog",
		})
		f.Entry = []*Entry{
			{
				Title:   "Library",
				ID:      "library",
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/library", Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: "Choose an author of a book",
				},
			},
			{
				Title:   "Search Results",
				ID:      "search-results",
				Updated: f.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/search-results", Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: "Choose a genre of a book",
				},
			},
		}

		writeFeed(w, http.StatusOK, *f)
	}
}

func (o *OPDS) opensearchMetaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data :=
			`
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
<ShortName>OpenBooks OPDS</ShortName>
<Description>New IRC Search</Description>
<InputEncoding>UTF-8</InputEncoding>
<OutputEncoding>UTF-8</OutputEncoding>
<Url type="application/atom+xml;profile=opds-catalog;kind=navigation" template="/opds/search?q={searchTerms}"/>
</OpenSearchDescription>	
`
		w.Header().Add("Content-Type", "application/atom+xml")
		w.WriteHeader(http.StatusOK)

		s := fmt.Sprintf("%s%s", xml.Header, data)
		w.Write([]byte(s))
	}
}

func (o *OPDS) libraryHandler() http.HandlerFunc {
	selfHref := "/opds/library"
	return func(w http.ResponseWriter, r *http.Request) {

		bookID := r.URL.Query().Get("download")

		if bookID != "" {
			slog.Info("Sending download request.", slog.String("book-id", bookID))
		}

		files, err := os.ReadDir(o.libraryDirectory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		feed := NewFeed("OpenBooks OPDS", "", selfHref)

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.IsDir() {
				continue
			}

			fmt.Println("File: ", file.Name())
			author := strings.Split(file.Name(), " - ")[0]
			title := strings.Split(file.Name(), " - ")[1]

			feed.Entry = append(feed.Entry, &Entry{
				Title:   title,
				Authors: []Author{{Name: author}},
				ID:      file.Name(),
				Updated: feed.Time(info.ModTime()),
				Link: []Link{
					{Rel: "http://opds-spec.org/acquisition/open-access", Length: fmt.Sprint(file.Name()), Href: "/opds/library/" + file.Name(), Type: "application/epub+zip"},
				},
			})
		}

		writeFeed(w, http.StatusOK, *feed)
	}
}

func (o *OPDS) downloadBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, fileName := path.Split(r.URL.Path)
		bookPath := filepath.Join(o.libraryDirectory, fileName)

		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		w.Header().Add("Content-Type", fmt.Sprintf("application/epub+zip; name=%s", fileName))
		w.Header().Add("Content-Transfer-Encoding", "binary")

		http.ServeFile(w, r, bookPath)
	}
}

func (o *OPDS) searchHandler() http.HandlerFunc {
	selfHref := "/opds/search"
	return func(w http.ResponseWriter, r *http.Request) {
		// searchQuery := chi.URLParam(r, "query")
		searchQuery := r.URL.Query().Get("q")
		fmt.Println("Query: ", searchQuery)

		feed := NewFeed("OpenBooks OPDS IRC Search", "", selfHref)

		feed.Entry = []*Entry{
			{
				Title:   fmt.Sprintf("Search Undernet #bookz for '%s'", searchQuery),
				ID:      "search-undernet",
				Updated: feed.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/search-results?q=%s&server=%s", searchQuery, "undernet"), Type: FeedAcquisitionLinkType},
				},

				Content: &Content{
					Type:    FeedTextContentType,
					Content: fmt.Sprintf("Search for %s on Undernet #bookz", searchQuery),
				},
			},
			{
				Title:   fmt.Sprintf("Search IRC Highway #ebooks for '%s'", searchQuery),
				ID:      "search-irchighway",
				Updated: feed.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/search-results?q=%s&server=%s", searchQuery, "irchighway"), Type: FeedNavigationLinkType},
				},
				Content: &Content{
					Type:    FeedTextContentType,
					Content: "Choose a genre of a book",
				},
			},
		}

		writeFeed(w, http.StatusOK, *feed)
	}
}

// searchResultsHandler will send search to IRC server if the necessary query params are set.
// It also shows a list of past search results.
func (o *OPDS) searchResultsHandler() http.HandlerFunc {
	selfHref := "/opds/search-results"

	return func(w http.ResponseWriter, r *http.Request) {
		searchQuery := r.URL.Query().Get("q")
		server := r.URL.Query().Get("server")

		if searchQuery != "" && server != "" {
			// TODO: Send query to correct IRC Server here...
			go o.searcher(searchQuery, server)
		}

		feed := NewFeed("OpenBooks OPDS IRC Search", "", selfHref)

		fmt.Println("Search results listing")

		files, err := os.ReadDir(o.searchResultDirectory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		feed.Entry = []*Entry{
			{
				Title:   "Refresh Search Results",
				ID:      "search-undernet",
				Updated: feed.Time(time.Now()),
				Link: []Link{
					{Rel: FeedSelfLinkRel, Href: selfHref, Type: FeedNavigationLinkType},
				},
			},
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.IsDir() {
				continue
			}

			feed.Entry = append(feed.Entry, &Entry{
				Title:   file.Name(),
				ID:      file.Name(),
				Updated: feed.Time(info.ModTime()),
				Link: []Link{
					{Rel: FeedSubsectionLinkRel, Href: "/opds/search-results/" + file.Name(), Type: FeedAcquisitionLinkType},
				},
			})
		}

		writeFeed(w, http.StatusOK, *feed)
	}
}

func (o *OPDS) displaySearchResultsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, fileName := path.Split(r.URL.Path)
		resultsPath := filepath.Join(o.searchResultDirectory, fileName)

		feed := NewFeed("OpenBooks OPDS IRC Search", "", fmt.Sprintf("/opds/search-results/%s", fileName))

		results, _, err := core.ParseSearchFile(resultsPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, result := range results {
			feed.Entry = append(feed.Entry, &Entry{
				Title:   result.Title,
				Authors: []Author{{Name: result.Author}},
				Source:  result.Server,
				ID:      result.Title,
				Updated: feed.Time(time.Now()),
				Link: []Link{
					{Rel: "http://opds-spec.org/acquisition/buy", Href: fmt.Sprintf("/opds/library?download=%s", result.Full), Type: FeedAcquisitionLinkType},
				},
			})
		}

		// for _, error := range errors {
		// 	feed.Entry = append(feed.Entry, &Entry{
		// 		Title:   error.Line,
		// 		ID:      error.Line,
		// 		Updated: feed.Time(time.Now()),
		// 		Link: []Link{
		// 			{Rel: FeedSubsectionLinkRel, Href: fmt.Sprintf("/opds/library?download=%s", error.Line), Type: FeedNavigationLinkType},
		// 		},
		// 	})
		// }

		writeFeed(w, http.StatusOK, *feed)

	}
}

func (o *OPDS) searcher(query, server string) {
	fmt.Println("Searching for: ", query, " on ", server)
	inFile, err := os.Open("C:/Users/evanb/Downloads/oatmeal.txt")
	if err != nil {
		return
	}

	defer inFile.Close()

	outFile, err := os.Create(filepath.Join(o.searchResultDirectory, fmt.Sprintf("Searchbot_results_for_%s", query)))
	if err != nil {
		return
	}
	defer outFile.Close()

	reader := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	for reader.Scan() {
		text := reader.Text()
		if !strings.HasPrefix(text, "!") {
			continue
		}

		if !strings.Contains(text, query) {
			continue
		}

		writer.WriteString(fmt.Sprintf("%s\n", text))
	}
}
