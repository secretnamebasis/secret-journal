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
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2/app"
)

func main() {
	// create a new application
	a := app.NewWithID(APP_ID)

	// create themes
	a.Settings().SetTheme(themes.main)

	session.window = // give this session a window
		a.NewWindow(APP_NAME) // and name it

	// this is the master window
	session.window.SetMaster()

	// set closing procedure
	session.window.SetCloseIntercept(
		func() { // use a anonymous function
			reset()                // reset the application
			session.window.Close() // close the window
			// and tell os to exit app
			os.Exit(0)
		},
	)

	// give window some padding
	session.window.SetPadded(true)

	// establish where you are
	session.domain = appLanding

	// get centered
	session.window.CenterOnScreen()

	// set icon
	a.SetIcon(resourceIconPng)
	session.window.SetIcon(resourceIconPng)

	// create a welcome screen
	fmt.Printf(versionMsg, version)
	fmt.Println(copyrightMsg)
	fmt.Printf(
		osArchGoMaxMsg,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.GOMAXPROCS(0),
	)
	fmt.Println(quoteMsg)

	// you are in the main part of the app
	session.domain = appMain

	ui.maxwidth = 360
	ui.maxheight = 680

	ui.width = ui.maxwidth * 0.9
	ui.height = ui.maxheight
	ui.padding = ui.maxwidth * 0.5

	resizeWindow(ui.maxwidth, ui.maxheight)

	session.window.SetContent(layoutMain())
	session.window.ShowAndRun()

}
