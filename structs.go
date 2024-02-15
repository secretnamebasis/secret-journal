/*
secret-journal: a journal stored on a decentralized, homomorphically-encrypted, blockchain called DERO (https://dero.io)
Copyright (C) 2024  secretnamebasis

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type transportWithBasicAuth struct {
	Username string
	Password string
	Base     http.RoundTripper
}

type appThemes struct {
	main theme1
	alt  theme2
}

type theme1 struct{}

type theme2 struct{}

type appSession struct {
	window fyne.Window
	domain string
}

type appUI struct {
	padding   float32
	maxwidth  float32
	width     float32
	maxheight float32
	height    float32
}

type entryInfo struct {
	Index          int
	TimeStr        string
	Text           string
	Author         string
	AuthorLabel    *widget.Label
	TimeLabel      *widget.Label
	TextLabel      *widget.RichText
	IndexLabel     *widget.Label
	EntryContainer *fyne.Container
}
