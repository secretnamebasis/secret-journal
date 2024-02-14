package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func layoutMain() fyne.CanvasObject {
	session.window.SetFixedSize(false)

	contentContainer.Hide()
	scrollContainer = container.NewVScroll(contentContainer)
	scrollContainer.SetMinSize(fyne.NewSize(ui.maxwidth, ui.maxheight))
	entryForm := widget.NewEntry()
	entryForm.MultiLine = true
	entryForm.Wrapping = fyne.TextWrapWord
	entryForm.SetMinRowsVisible(1)
	entryForm.PlaceHolder = "Enter Text Here..."

	resultLabel := widget.NewLabel("Status: New")
	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")
	searchEntry.OnSubmitted = func(query string) {
		searchTransfers(
			query,
			contentContainer,
		)
		pause()
		contentContainer.Refresh()
	}
	searchEntry.Disable()

	visbilityButton = widget.NewButtonWithIcon(
		"",
		theme.VisibilityOffIcon(),
		func() {
			// Toggle visibility state
			isVisibilityOn = !isVisibilityOn
			pause()
			// Update the button's icon based on the visibility state
			if isVisibilityOn {
				visbilityButton.SetIcon(theme.VisibilityIcon())
				contentContainer.Show()

			} else {

				visbilityButton.SetIcon(theme.VisibilityOffIcon())
				contentContainer.Hide()

			}

		},
	)

	visbilityButton.Disable()
	refreshButton = widget.NewButtonWithIcon(
		"",
		theme.ViewRefreshIcon(),
		func() {
			pause()
			updateTransfers(contentContainer)
			scrollContainer.Refresh()

		},
	)
	refreshButton.Disable()

	entryButton = widget.NewButtonWithIcon(
		"",
		theme.MailSendIcon(),
		func() {
			pause()

			processEntrySubmission(
				entryForm,
				entryButton,
				resultLabel,
				contentContainer,
			)
			scrollContainer.ScrollToBottom()

		},
	)
	entryButton.Disable()

	logoutButton = widget.NewButtonWithIcon(
		"Logout",
		theme.LogoutIcon(),
		func() {
			reset()
		},
	)
	logoutButton.Disable()

	connectButton := widget.NewButtonWithIcon(
		"",
		theme.SettingsIcon(),
		func() {
			showSettingsWindow(
				session.window,
			)

		},
	)

	toolbarContainer := container.NewBorder(
		container.NewVBox(
			padding,
			container.NewGridWithColumns(
				4,
				connectButton,
				refreshButton,
				visbilityButton,
				searchEntry,
			),
		),
		scrollContainer,
		nil,
		nil,
	)
	buttonContainer := container.NewGridWrap(
		fyne.NewSize(ui.width*0.18, ui.maxheight*.065),
		entryButton,
	)

	entryContainer := container.NewGridWrap(
		fyne.NewSize(ui.width*0.9, ui.maxheight*.065),
		entryForm,
	)

	buttonsContainer := container.NewHBox(
		entryContainer,
		buttonContainer,
	)

	chatBarContainer := container.NewVBox(
		resultLabel,
		buttonsContainer,
	)

	layout := container.NewBorder(
		toolbarContainer,
		chatBarContainer,
		nil,
		nil,
	)

	return layout
}
