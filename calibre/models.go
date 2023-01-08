package calibre

import "github.com/go-playground/validator/v10"

type Config struct {
	// Calibre Content Server URL. Example) http://calibre.homeserver, http://192.168.1.1:8080
	URL           string `validate:"required,url"`
	Username      string `validate:"required_with=Password"`
	Password      string `validate:"required_with=Username"`
	AddDuplicates bool
}

func (config *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(config)
}

type AddBookResponse struct {
	Title     string   `json:"title"`
	Authors   []string `json:"authors"`
	Languages []string `json:"languages"`
	Filename  string   `json:"filename"`
	Id        string   `json:"id"`
	// Only populated if a book was created
	BookId int `json:"book_id"`
	// Only populated if there are duplicate titles
	Duplicates []struct {
		Title   string   `json:"title"`
		Authors []string `json:"authors"`
	} `json:"duplicates"`
}
