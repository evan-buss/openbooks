package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSearchParser(t *testing.T) {
	reader := strings.NewReader(sampleData)
	results, errors := ParseSearch(reader)

	if len(errors) > 1 {
		t.Errorf("Expected 1 errors but got %d\n", len(errors))
		for _, parseError := range errors {
			t.Log(parseError)
		}
	}

	if len(results) != 26 {
		t.Errorf("Expected 26 results but got %d\n", len(results))
	}
}

func TestSearchParserV2(t *testing.T) {
	reader := strings.NewReader(sampleData)
	results, errors := ParseSearchV2(reader)

	if len(errors) != 1 {
		t.Errorf("Expected 1 errors but got %d\n", len(errors))
		for _, parseError := range errors {
			t.Log(parseError)
		}
	}

	if len(results) != 26 {
		t.Errorf("Expected 26 results but got %d\n", len(results))
	}
}

func TestSpecialCases(t *testing.T) {
	cases := []struct {
		reason   string
		original string
		download BookDetail
	}{
		{
			"info block, file size, title case file format, rar file",
			"!DV8 F. Scott Fitzgerald - The Great Gatsby (Epub).rar  ::INFO:: 394.7KB",
			BookDetail{
				Server: "DV8",
				Author: "F. Scott Fitzgerald",
				Title:  "The Great Gatsby (Epub)",
				Format: "epub",
				Size:   "394.7KB",
				Full:   "!DV8 F. Scott Fitzgerald - The Great Gatsby (Epub).rar",
			},
		},
		{
			"no info block, no file size",
			"!Horla F Scott Fitzgerald - The Great Gatsby (retail) (epub).epub",
			BookDetail{
				Server: "Horla",
				Author: "F Scott Fitzgerald",
				Title:  "The Great Gatsby (retail) (epub)",
				Format: "epub",
				Size:   "N/A",
				Full:   "!Horla F Scott Fitzgerald - The Great Gatsby (retail) (epub).epub",
			},
		},
		{
			"hash code, author/title swapped position",
			"!Ook So we Read on -How the Great Gatsby came to be and why it Endures (2014) - Maureen Corrigan.epub  ::INFO:: 5MB ::HASH:: dde55317998f25aa",
			BookDetail{
				Server: "Ook",
				Author: "So we Read on -How the Great Gatsby came to be and why it Endures (2014)",
				Title:  "Maureen Corrigan",
				Format: "epub",
				Size:   "5MB",
				Full:   "!Ook So we Read on -How the Great Gatsby came to be and why it Endures (2014) - Maureen Corrigan.epub",
			},
		},
	}

	for _, input := range cases {
		result, err := parseLineV2(input.original)
		require.NoError(t, err)
		assert.Equal(t, result, input.download)
	}
}

var sampleData = `Search results from SearchBot v3.00.07 by Ook, searching dll written by Iczelion, Based on Searchbot v2.22 by Dukelupus
Searched 20 lists for "the great gatsby" , found 27 matches. Enjoy!
This list includes results from ALL the lists SearchBot v3.00.07 currently has, some of these servers may be offline.
Always check to be sure the server you want to make a request from is actually in the channel, otherwise your request will have no effect.
For easier searching, use sbClient script (also very fast local searches). You can get that script by typing @sbClient in the channel.




!dragnbreaker Fitzgerald, F Scott - Novel 03 - The Great Gatsby (retail).epub  ::INFO:: 1.7MB
!DV8 F. Scott Fitzgerald - The Great Gatsby (Epub).rar  ::INFO:: 394.7KB
!Horla F Scott Fitzgerald - The Great Gatsby (retail) (epub).epub
!Horla F. Scott Fitzgerald - The Great Gatsby (V1.5 RTF).rtf
!Horla Sarah Churchwell - Careless People- Murder, Mayhem, and the Invention of the Great Gatsby (epub).epub
!JimBob420 F. Scott Fitzgerald - The Great Gatsby (V1.5 RTF).rar ::INFO:: 272.23KB
!JimBob420 F Scott Fitzgerald - The Great Gatsby (epub).rar ::INFO:: 204.54KB
!JimBob420 F Scott Fitzgerald - The Great Gatsby (retail) (epub).rar ::INFO:: 1.65MB
!JimBob420 Sarah Churchwell - Careless People- Murder, Mayhem, and the Invention of the Great Gatsby (epub).rar ::INFO:: 8.44MB
!MusicWench F Scott Fitzgerald - The Great Gatsby.epub  ::INFO:: 332.7KB
!MusicWench F Scott Fitzgerald - The Great Gatsby.mobi  ::INFO:: 376.6KB
!Oatmeal F. Scott Fitzgerald - The Great Gatsby (V1.5 RTF).rar ::INFO:: 272.23KB
!Oatmeal F Scott Fitzgerald - The Great Gatsby (epub).rar ::INFO:: 204.55KB
!Oatmeal F Scott Fitzgerald - The Great Gatsby (retail) (epub).rar ::INFO:: 1.65MB
!Oatmeal Sarah Churchwell - Careless People- Murder, Mayhem, and the Invention of the Great Gatsby (epub).rar ::INFO:: 8.44MB
!Ook So we Read on -How the Great Gatsby came to be and why it Endures (2014) - Maureen Corrigan.epub  ::INFO:: 5MB ::HASH:: dde55317998f25aa
!Ook Sarah Churchwell - Careless People- Murder, Mayhem, and the Invention of the Great Gatsby (epub).rar  ::INFO:: 8MB ::HASH:: 348c62174a5c5c29
!Ook F Scott Fitzgerald - The Great Gatsby (retail) (epub).rar  ::INFO:: 1MB ::HASH:: 8d860602f0f43789
!peapod F Scott Fitzgerald - Great Gatsby, The.azw3  ::INFO:: 260.46KB
!peapod F Scott Fitzgerald - Great Gatsby, The.epub  ::INFO:: 373.54KB
!peapod F Scott Fitzgerald - Great Gatsby, The.mobi  ::INFO:: 368.87KB
!peapod Sarah Churchwell - Careless People- Murder, Mayhem, and the Invention of the Great Gatsby (epub).rar  ::INFO:: 8.44MB
!peapod The Great Gatsby.pdf  ::INFO:: 254.73KB
!peapod The Great Gatsby - F Scott Fitzgerald.mobi  ::INFO:: 246.10KB
!phoomphy Fitzgerald, F. Scott - The Great Gatsby (1925).epub     ::INFO:: 205.10 KiB
!phoomphy Fitzgerald, F. Scott - The Great Gatsby.pdf     ::INFO:: 775.69 KiB
!phoomphy Call of Cthulhu - Gatsby and the Great Race (monograph #0324).pdf     ::INFO:: 20.23 MiB
`
