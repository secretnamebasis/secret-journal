package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	x "fyne.io/x/fyne/widget"
)

func layoutMain() fyne.CanvasObject {
	session.window.SetFixedSize(true)
	options := []string{}
	deroDestination = x.NewCompletionEntry(options)
	contentContainer.Hide()
	scrollContainer = container.NewVScroll(contentContainer)
	scrollContainer.SetMinSize(fyne.NewSize(ui.maxwidth, ui.maxheight))
	entryForm := widget.NewEntry()
	entryForm.MultiLine = true
	entryForm.Wrapping = fyne.TextWrapWord
	entryForm.SetMinRowsVisible(1)
	entryForm.PlaceHolder = "Enter Text Here..."

	resultLabel = widget.NewLabel("Status: Not Logged In")
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
			switch {
			case deroDestination.Text == "":
				resultLabel.SetText("Enter receiving address")
			case validateAddress(deroDestination.Text):
				resultLabel.SetText(":)")
				destinationAddress = deroDestination.Text
				updateTransfers(contentContainer)
			case deroDestination.Text != "":

				truncatedText := deroDestination.Text
				if original, ok := originalToTruncated[truncatedText]; ok {
					// Original value found
					resultLabel.SetText(":)")
					destinationAddress = original

					updateTransfers(contentContainer)
				} else {
					// Truncated value not found
					resultLabel.SetText(":(")
				}

			}
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
	deroDestination.PlaceHolder = truncateAddress(DEVELOPER_ADDRESS, 6, 7)
	// When the use typed text, complete the list.
	deroDestination.OnChanged = func(s string) {
		// Completion starts for text length >= 3
		if len(s) < 3 {
			deroDestination.HideCompletion()
			return
		}

		// Filter options that contain the typed text
		filteredOptions := filterOptions(originalToTruncated, deroDestination.Text)

		// No matching results
		if len(filteredOptions) == 0 {
			deroDestination.HideCompletion()
			return
		}

		// Show filtered options
		deroDestination.SetOptions(filteredOptions)
		deroDestination.ShowCompletion()
		updateTransfers(contentContainer)
	}
	contactButton := widget.NewButtonWithIcon(
		"",
		theme.AccountIcon(),
		func() {
			showContactWindow(
				session.window,
			)

		},
	)
	deroDestination.Disable()

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
		fyne.NewSize(ui.width*0.18, ui.maxheight*.1),
		entryButton,
	)

	entryContainer := container.NewGridWrap(
		fyne.NewSize(ui.width*0.92, ui.maxheight*.1),
		entryForm,
	)

	buttonsContainer := container.NewHBox(
		contactButton,
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
