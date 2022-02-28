package main

import "github.com/evan-buss/openbooks/core"

type Book struct {
	core.BookDetail
	Line string `json:"line"`
}

func NewBook(detail core.BookDetail, line string) Book {
	return Book{detail, line}
}

func (b Book) Type() string {
	return "book"
}
