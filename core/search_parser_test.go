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

	if len(results) != 57 {
		t.Errorf("Expected 57 results but got %d\n", len(results))
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

	if len(results) != 57 {
		t.Errorf("Expected 57 results but got %d\n", len(results))
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
		{
			"has a weird %some_text% prefix on the title",
			"!FWServer %F77FE9FF1CCD% Michael Haag - Inferno Decoded - The Essential Companion To The Myths, Mysteries And Locations Of Dan Brown's Inferno.epub  ::INFO:: 8.00MB",
			BookDetail{
				Server: "FWServer",
				Author: "Michael Haag",
				Title:  "Inferno Decoded - The Essential Companion To The Myths, Mysteries And Locations Of Dan Brown's Inferno",
				Format: "epub",
				Size:   "8.00MB",
				Full:   "!FWServer %F77FE9FF1CCD% Michael Haag - Inferno Decoded - The Essential Companion To The Myths, Mysteries And Locations Of Dan Brown's Inferno.epub",
			},
		},
		{
			"has a weird %some_text% prefix on the title, audiobook with valid eBook format",
			"!FWServer %DE7B9E7F6F34% Brown, Dan - Robert Langdon 04 - Inferno - Audiobook.zip  ::INFO:: 445.09MB",
			BookDetail{
				Server: "FWServer",
				Author: "Brown, Dan",
				Title:  "Robert Langdon 04 - Inferno - Audiobook",
				Format: "zip",
				Size:   "445.09MB",
				Full:   "!FWServer %DE7B9E7F6F34% Brown, Dan - Robert Langdon 04 - Inferno - Audiobook.zip",
			},
		},
	}

	for _, input := range cases {
		result, err := parseLineV2(input.original)
		require.NoError(t, err)
		assert.Equal(t, input.download, result)
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
!FWServer %F77FE9FF1CCD% Michael Haag - Inferno Decoded - The Essential Companion To The Myths, Mysteries And Locations Of Dan Brown's Inferno.epub  ::INFO:: 8.00MB
!FWServer %DE7B9E7F6F34% Brown, Dan - Robert Langdon 04 - Inferno - Audiobook.zip  ::INFO:: 445.09MB
!DeathCookie Emma_L_Adams_Heritage_of_Fire_03_Inferno.epub.rar  ::INFO:: 530.0KB
!DeathCookie J._L._Weil_Divisa_Huntress_02_Inferno_of_Darkness.epub.rar  ::INFO:: 330.9KB
!DeathCookie Jan_Stryvant_Dan's_Inferno_01_Cursed!.epub.rar  ::INFO:: 212.5KB
!DeathCookie Jan_Stryvant_Dan's_Inferno_02_BeDeviled.epub.rar  ::INFO:: 222.4KB
!DeathCookie Jan_Stryvant_Dan's_Inferno_03_Heritage.epub.rar  ::INFO:: 311.6KB
!DeathCookie Jan_Stryvant_Dan's_Inferno_04_Vengeance.epub.rar  ::INFO:: 341.4KB
!DeathCookie Simon_Archer_Super_Hero_Academy_03_Inferno_Island.epub.rar  ::INFO:: 293.3KB
!DeathCookie Travis_Bagwell_Tarot_03_Inferno.epub.rar  ::INFO:: 579.5KB
!dragnbreaker Cooper, Louise - Indigo 02 - Inferno.doc  ::INFO:: 791.9KB
!dragnbreaker Allen, Roger MacBride - Isaac Asimov's Caliban 02 - Inferno.htm  ::INFO:: 109.0B
!dragnbreaker Allen, Roger MacBride - Isaac Asimov's Caliban 02 - Inferno.jpg  ::INFO:: 52.5KB
!dragnbreaker Inferno! 030 [Black Library] (2002) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 17.1MB
!dragnbreaker Inferno! 001 [Black Library] (1997) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 16.6MB
!dragnbreaker Inferno! 003 [Black Library] (1997) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 15.6MB
!dragnbreaker Inferno! 004 (Black Library) (1997) (UK) (CBRed by Discovery-DCP).cbr  ::INFO:: 33.2MB
!dragnbreaker Inferno! 006 (Black Library) (1998) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 14.9MB
!dragnbreaker Inferno! 013 [Black Library] (1999) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 17.5MB
!dragnbreaker Inferno! 024 [Black Library] (2001) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 16.5MB
!dragnbreaker Inferno! 025 [Black Library] (2001) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 15.8MB
!dragnbreaker Inferno! 029 [Black Library] (2002) (U.K.) (CBRed by Discovery-DCP).cbr  ::INFO:: 17.5MB
!Horla Annmarie Ortega - Dante's Inferno (lit).lit
!Horla Bianca D'Arc - Inferno (lit).lit
!Horla Niven, Larry - [Inferno 01] - Inferno (v5.0).lit
!Horla Linda Howard -[Raintree 01]- Inferno.doc
!Horla Linda Howard - [Raintree 01] - Inferno (lit).lit
!Horla Martin Amis - The Moronic Inferno & Other Visits to America.lit
!Horla Monnery, David - The Bosnian Inferno.txt.RAR
!Horla Roger MacBride Allen - Caliban 02 - Inferno.lit
!Horla Annmarie Ortega - Dante's Inferno (lit).lit
!Horla Bianca D'Arc - Inferno (lit).lit
`
