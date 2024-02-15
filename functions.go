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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/color"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	x "fyne.io/x/fyne/widget"
	"github.com/deroproject/derohe/rpc"
	"github.com/google/martian/log"

	"github.com/ybbus/jsonrpc"
)

func pause() { time.Sleep(DEFAULT_WAIT_TIME * time.Millisecond) }

// RoundTrip implements the RoundTripper interface
func (t *transportWithBasicAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Base.RoundTrip(req)
}
func hash(text string) []byte {
	hash := sha256.New()
	hash.Write([]byte(text))
	hashInBytes := hash.Sum(nil)
	return hashInBytes
}

func walletConnection(deroUsername, deroPassword string, reset bool) error {
	if reset {
		deroRpcClient = nil
		deroHttpClient = nil
		return nil
	}

	deroHttpClient = &http.Client{
		Transport: &transportWithBasicAuth{
			Username: deroUsername,
			Password: deroPassword,
			Base:     http.DefaultTransport,
		},
	}

	deroRpcClient = jsonrpc.NewClientWithOpts(
		deroEndpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: deroHttpClient,
		},
	)

	address, err := getAddress()

	if err != nil {
		return fmt.Errorf("error connecting to wallet: %s", err)
	}

	// Optional: Log the wallet version
	fmt.Printf("Client ID: %s\n", hex.EncodeToString(hash(address.String())))

	return nil
}

func getHeight() int {
	err = deroRpcClient.CallFor(
		&walletHeight,
		"GetHeight",
	)
	if err != nil || walletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(walletHeight.Height)
}

func getAddress() (*rpc.Address, error) {
	err = deroRpcClient.CallFor(&addr_result, "GetAddress")
	if err != nil || addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return nil, err
	}

	address, err = rpc.NewAddress(addr_result.Address)
	if err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", addr_result.Address, err)
		return nil, err
	}
	return address, nil
}

func getTransferByTXID(txid string) rpc.Get_Transfer_By_TXID_Result {
	var transfer rpc.Get_Transfer_By_TXID_Result
	var params rpc.Get_Transfer_By_TXID_Params
	params.TXID = txid
	_ = deroRpcClient.CallFor(
		&transfer,
		"GetTransferbyTXID",
		params,
	)

	if transfer.Entry.Time.String() == "" {
		log.Infof("Time is \"\" string")
	}

	return transfer
}

func getBalance() (rpc.GetBalance_Result, error) {

	err = deroRpcClient.CallFor(
		&balance,
		"GetBalance",
	)
	return balance, nil
}

func sendTransfer(params rpc.Transfer_Params) (rpc.Transfer_Result, error) {
	var transfers rpc.Transfer_Result
	err = deroRpcClient.CallFor(
		&transfers,
		"Transfer",
		params,
	)

	if err != nil {
		return transfers, err
	}

	return transfers, nil
}


func getAllTransfers() (rpc.Get_Transfers_Result, error) {

	err = deroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:              true,
			Out:             true,
			Coinbase:        false,
			DestinationPort: uint64(0),
		},
	)
	if err != nil {
		log.Errorf("Could not obtain gettransfers from wallet: %v", err)
		return transfers, err
	}

	return transfers, nil
}

func getIncomingTransfers() (rpc.Get_Transfers_Result, error) {

	err = deroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:              true,
			Out:             false,
			Coinbase:        false,
			DestinationPort: uint64(0),
		},
	)
	if err != nil {
		log.Errorf("Could not obtain gettransfers from wallet: %v", err)
		return transfers, err
	}

	return transfers, nil
}


func getOutgoingTransfers() (rpc.Get_Transfers_Result, error) {

	err = deroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:              false,
			Out:             true,
			Coinbase:        false,
			DestinationPort: uint64(0),

			// Receiver:        destinationAddress,

		},
	)
	if err != nil {
		log.Errorf("Could not obtain gettransfers from wallet: %v", err)
		return transfers, err
	}

	return transfers, nil
}

func truncateAddress(addr string, prefixLen, suffixLen int) string {
	if len(addr) <= prefixLen+suffixLen {
		return addr
	}
	return fmt.Sprintf("%s....%s", addr[:prefixLen], addr[len(addr)-suffixLen:])
}

func truncateTXID(txid string, prefixLen, suffixLen int) string {
	if len(txid) <= prefixLen+suffixLen {
		return txid
	}
	return fmt.Sprintf("%s....%s", txid[:prefixLen], txid[len(txid)-suffixLen:])
}

func updateContacts(deroDestination *x.CompletionEntry, resultLabel *widget.Label) {
	var truncatedOptions []string

	data, _ := getOutgoingTransfers()

	for _, e := range data.Entries {
		option := e.Destination

		// Check if the option is not already in the map
		if !uniqueOptions[option] {
			uniqueOptions[option] = true

			// Truncate the option
			truncatedOption := truncateAddress(option, 4, 4)

			// Check if the truncated option is not already in the map
			if !uniqueOptions[truncatedOption] {
				uniqueOptions[truncatedOption] = true
				truncatedOptions = append(truncatedOptions, truncatedOption)

				// Map the original option to its truncated version
				originalToTruncated[truncatedOption] = option
			}
		}
	}
	deroDestination.SetOptions(truncatedOptions)
}

func connectWallet() {
	go func() {
		updateContacts(deroDestination, resultLabel)

		updateWallet(lblAddress)
		updateHeight(lblHeight)
		updateBalance(lblBalance)
		updateTransfers(contentContainer)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				updateHeight(lblHeight)
				updateBalance(lblBalance)
			case <-time.After(time.Second):
				// Add any other periodic tasks here

			case <-resetCh:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateWallet(lblAddress *widget.Label) {
	addr, err := getAddress()
	if err != nil {
		log.Errorf("not connected to wallet: %s", err)
		return
	}
	truncatedAddr := truncateAddress(addr.String(), 4, 2)
	lblAddress.Text = truncatedAddr

	lblAddress.SetText("A: " + truncatedAddr)

}

func updateHeight(lblHeight *widget.Label) {
	h := getHeight()
	lbl := fmt.Sprintf("H: %d", h)
	lblHeight.SetText(lbl)

}

func updateBalance(lblBalance *widget.Label) {
	b, _ := getBalance()

	balanceFloat := float64(b.Balance) / 1e5

	formattedBalance := strconv.FormatFloat(balanceFloat, 'f', 5, 64)

	lbl := fmt.Sprintf("B: %s", formattedBalance)
	lblBalance.SetText(lbl)
}

// Function to split a string into chunks of a specified size
func chunkString(s string, chunkSize int) []string {

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func processEntry(text string) (rpc.Transfer_Result, error) {

	// Check if the entry exceeds the character limit
	if len(text) <= CHUNKSIZE {
		// If the entry is within the character limit, proceed as before
		return processSingleEntry(text)
	}

	chunks := chunkString(text, CHUNKSIZE)
	// Create a list to store transfers
	var txs []rpc.Transfer

	// Iterate over the chunks and create a transfer for each
	for _, text := range chunks {
		tx := prepareTransfer(text)
		txs = append(txs, tx...)
	}

	tip := prepareTip()
	endTx := append(txs, tip)

	params := prepareParams(
		endTx,
	)

	result, err := sendTransfer(params)

	if err != nil {
		return result, err
	}
	// Send the transfers
	return result, nil
}

func prepareTransfer(text string) []rpc.Transfer {
	payload := rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    text,
		},
		{
			Name:     rpc.RPC_REPLYBACK_ADDRESS,
			DataType: rpc.DataAddress,
			Value:    address,
		},
	}

	transfer := rpc.Transfer{
		Destination: destinationAddress,
		Amount:      uint64(0),
		Payload_RPC: payload,
	}

	return []rpc.Transfer{transfer}
}

func prepareTip() rpc.Transfer {

	receipt := rpc.Arguments{
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    tipMsg,
		},
	}

	transfer := rpc.Transfer{
		Destination: DEVELOPER_ADDRESS,
		Amount:      tipAmt,
		Payload_RPC: receipt,
	}

	return transfer
}

func prepareParams(transfers []rpc.Transfer) rpc.Transfer_Params {

	transferParams := rpc.Transfer_Params{
		Transfers: transfers,
	}

	return transferParams

}

func processSingleEntry(text string) (rpc.Transfer_Result, error) {
	// Send the transfers

	return sendTransfer(prepareParams(prepareTransfer(text)))
}

func showContactWindow(
	w fyne.Window,
) {

	deroDestination.Validator = func(s string) (err error) {

		switch {
		case deroDestination.Text == "":
			resultLabel.SetText("Enter receiving address")
		case validateAddress(deroDestination.Text):
			resultLabel.SetText(":)")
			destinationAddress = deroDestination.Text
		case deroDestination.Text != "":

			truncatedText := deroDestination.Text
			if original, ok := originalToTruncated[truncatedText]; ok {
				// Original value found
				resultLabel.SetText(":)")
				destinationAddress = original

			} else {
				// Truncated value not found
				resultLabel.SetText(":(")
			}

		}

		return nil
	}

	lbl := widget.NewLabel("Choose Contact")

	formWidget := container.NewGridWrap(
		fyne.NewSize(ui.width, 45),
		deroDestination,
	)

	closeButton := widget.NewButton(
		"Close",
		func() {
			if modal != nil {
				modal.Hide()
			}
		})

	closeButton.OnTapped = func() {
		if modal != nil {
			modal.Hide()
		}
	}

	modalContent := container.NewCenter(
		container.NewVBox(
			lbl,
			container.NewCenter(
				container.NewHBox(
					formWidget,
					closeButton,
				),
			),
		),
	)
	modal = widget.NewModalPopUp(modalContent, w.Canvas())

	modal.Show()
}

func showSettingsWindow(
	w fyne.Window,
) {
	current := container.NewGridWithColumns(
		3,
		lblAddress,
		lblHeight,
		lblBalance,
	)
	lblLogin := widget.NewLabel("Wallet RPC Login")

	username := widget.NewEntry()
	username.PlaceHolder = "Username"

	password := widget.NewPasswordEntry()
	password.PlaceHolder = "Password"

	loginButton := widget.NewButtonWithIcon(
		"Login",
		theme.LoginIcon(),
		func() {
			deroUsername = username.Text
			deroPassword = password.Text

			if err := walletConnection(
				deroUsername,
				deroPassword,
				false,
			); err != nil {
				log.Errorf("wallet connection err: %s", err)
				return
			}

			connectWallet()
			pause()
			scrollContainer.ScrollToBottom()
			pause()
			entryButton.Enable()
			searchEntry.Enable()
			visbilityButton.Enable()
			logoutButton.Enable()
			refreshButton.Enable()
			deroDestination.Enable()

		},
	)

	formWidget := container.NewVBox(
		lblLogin,
		username,
		password,
	)

	closeButton := widget.NewButton(
		"Close",
		func() {
			if modal != nil {
				modal.Hide()
			}
		})

	closeButton.OnTapped = func() {
		if modal != nil {
			modal.Hide()
		}
	}

	modalContent := container.NewCenter(
		container.NewVBox(
			current,
			formWidget,
			padding,
			container.NewCenter(
				container.NewHBox(
					loginButton,
					closeButton,
					logoutButton,
				),
			),
		),
	)
	modal = widget.NewModalPopUp(modalContent, w.Canvas())

	modal.Show()
}

func validateAddress(a string) bool {
	return len(a) == 66
}

func processEntrySubmission(
	entry *widget.Entry,
	entryButton *widget.Button,
	resultLabel *widget.Label,
	contentContainer *fyne.Container,
) {
	entryButton.SetText("")
	entryButton.SetIcon(theme.MediaPauseIcon())
	entryButton.Disable()

	resultLabel.SetText("Obtaining TXID")

	entry.Disable()

	result, err := processEntry(entry.Text)
	if err != nil {
		resultLabel.SetText(
			fmt.Sprintf(
				"Status: Error %s",
				truncateTXID(result.TXID, 4, 4),
			),
		)
		resetEntryAfterSubmission(entry, entryButton, resultLabel, contentContainer)
	}
	resultLabel.SetText(
		fmt.Sprintf(
			"Status: Processing %s",
			truncateTXID(result.TXID, 4, 4),
		),
	)

	// initialHeight := getHeight()

	for {
		// currentHeight := getHeight()
		currentTXID := getTransferByTXID(result.TXID)

		if !currentTXID.Entry.Time.IsZero() {
			resetEntryAfterSubmission(entry, entryButton, resultLabel, contentContainer)
			return
		}

		pause()
	}
}

func resetEntryAfterSubmission(entry *widget.Entry, entryButton *widget.Button, resultLabel *widget.Label, contentContainer *fyne.Container) {
	resultLabel.SetText("Status: New")
	entry.SetText("")
	entry.Enable()
	resetButtons()
	updateTransfers(contentContainer)
	entryButton.Enable()

}

func resetButtons() {
	entryButton.SetIcon(theme.MailSendIcon())

}

func displayTransfers(
	entriesData []rpc.Entry,
	contentContainer *fyne.Container,
) {
	resetButtons()

	contentContainer.Objects = nil

	transfersMap := organizeTransfersByTime(entriesData)

	sortedTimes := sortTimestamps(transfersMap)

	// Display the first page of entries
	displayPage(contentContainer, transfersMap, sortedTimes, 0)

	contentContainer.Refresh()

}

func displayPage(
	contentContainer *fyne.Container,
	transfersMap map[string]string,
	sortedTimes []string,
	pageIndex int,
) {
	loadMoreButton := widget.NewButton(
		"Load More",
		nil,
	)
	loadMoreButton.Hide()
	loadMoreContainer := container.NewBorder(
		loadMoreButton,
		nil,
		nil,
		nil,
	)

	startIndex := pageIndex * PAGINATION
	endIndex := startIndex + PAGINATION - 1

	var entryInfos []*entryInfo

	if startIndex >= len(sortedTimes) {
		return
	}

	if endIndex >= len(sortedTimes) {
		endIndex = len(sortedTimes) - 1
	}

	// Create widgets for the specified range of entries
	entryInfos = createWidgetsAndAddToContainer(
		len(sortedTimes)-endIndex-1, // Start from the end
		len(sortedTimes)-startIndex-1,
		sortedTimes,
		transfersMap,
		contentContainer,
	)
	contentContainer.Objects = append([]fyne.CanvasObject{loadMoreContainer}, contentContainer.Objects...)

	// Add "Load More" button if there are more entries
	if endIndex < len(sortedTimes)-1 {
		loadMoreButton.Show()
		loadMoreButton.OnTapped = func() {
			// Load the next page of entries when the button is clicked
			contentContainer.Remove(loadMoreContainer)
			displayPage(
				contentContainer,
				transfersMap,
				sortedTimes,
				pageIndex+1,
			)
		}

	}

	// Use the last entry to trigger scrolling to the bottom
	if len(entryInfos) > 0 {
		lastEntry := entryInfos[len(entryInfos)-1]
		if lastEntry != nil {
			contentContainer.Refresh()
		}
	}
}

func organizeTransfersByTime(entriesData []rpc.Entry) map[string]string {
	transfersMap := make(map[string]string)

	for _, e := range entriesData {

		if e.Amount == 0 &&
			e.Payload_RPC.Has(
				rpc.RPC_COMMENT,
				rpc.DataString,
			) &&
			e.Payload_RPC.Has(
				rpc.RPC_REPLYBACK_ADDRESS,
				rpc.DataAddress,
			) {

			timeStr := e.Time.Format("2006-01-02 15:04:05")

			if _, ok := transfersMap[timeStr]; !ok {
				transfersMap[timeStr] = ""
			}

			transfersMap[timeStr] += e.Payload_RPC.Value(
				rpc.RPC_COMMENT,
				rpc.DataString,
			).(string)
		}
	}

	return transfersMap
}

func sortTimestamps(transfersMap map[string]string) []string {
	var s []string
	for timeStr := range transfersMap {
		s = append(s, timeStr)
	}
	slices.Sort(s)

	return s
}

// createWidgetsAndAddToContainer creates widgets for each entry, adds them to the container,
// and returns a slice of EntryInfo to associate data with each entry.
func createWidgetsAndAddToContainer(
	startIndex, endIndex int,
	sortedTimes []string,
	transfersMap map[string]string,
	contentContainer *fyne.Container,
) []*entryInfo {
	var entryInfos []*entryInfo

	for i := endIndex; i >= startIndex && i < len(sortedTimes); i-- {
		timeStr := sortedTimes[i]
		text := transfersMap[timeStr]

		timeLabel := widget.NewLabelWithStyle(
			timeStr,
			fyne.TextAlignLeading,
			fyne.TextStyle{
				Bold:      false,
				Italic:    true,
				Monospace: false,
			},
		)

		textLabel := widget.NewRichTextFromMarkdown(text)
		textLabel.Wrapping = fyne.TextWrapWord

		entryContainer = container.NewVBox(
			container.NewHBox(
				timeLabel,
			),
			textLabel,
		)

		contentContainer.Objects = append([]fyne.CanvasObject{entryContainer}, contentContainer.Objects...)

		// Create entryInfo and store it in the slice
		entryInfo := &entryInfo{
			Index:          i,
			TimeStr:        timeStr,
			Text:           text,
			TimeLabel:      timeLabel,
			TextLabel:      textLabel,
			EntryContainer: entryContainer,
		}
		entryInfos = append(entryInfos, entryInfo)
	}

	return entryInfos
}

func getOutgoingData() []rpc.Entry {
	data, err := getOutgoingTransfers()
	if err != nil {
		fmt.Println("Error fetching outgoing transfers:", err)
		return nil
	}

	return data.Entries
}

func updateTransfers(contentContainer *fyne.Container) {
	// Combine outgoing and incoming entries

	data := getOutgoingData()

	// Display the combined entries in the content container
	displayTransfers(data, contentContainer)

	// Scroll to the bottom of the content container
	scrollContainer.ScrollToBottom()
}

func filterEntriesByCondition(entries []rpc.Entry, condition func(e rpc.Entry) bool) []rpc.Entry {
	var filteredEntries []rpc.Entry
	for _, e := range entries {
		if condition(e) {
			filteredEntries = append(filteredEntries, e)
		}
	}
	return filteredEntries
}

func hasCommentPayload(e rpc.Entry) bool {
	return e.Amount == 0 && e.Payload_RPC != nil && e.Payload_RPC.Has(
		rpc.RPC_COMMENT,
		rpc.DataString,
	)
}

func containsQueryIgnoreCase(s, query string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(query))
}

func getCommentPayloadValue(e rpc.Entry) string {
	payloadValue, _ := e.Payload_RPC.Value(
		rpc.RPC_COMMENT,
		rpc.DataString,
	).(string)
	return payloadValue
}

func searchTransfers(query string, contentContainer *fyne.Container) {

	// Filter entries based on conditions
	filteredEntries := filterEntriesByCondition(
		getOutgoingData(),
		func(e rpc.Entry) bool {
			return hasCommentPayload(e) && containsQueryIgnoreCase(getCommentPayloadValue(e), query)
		})

	displayTransfers(filteredEntries, contentContainer)
}

func resizeWindow(width float32, height float32) {
	s := fyne.NewSize(width, height)
	session.window.Resize(s)
}

func reset() {

	close(resetCh)
	contentContainer.RemoveAll()
	lblAddress.SetText("Wallet: N/A")
	lblHeight.SetText("Height: N/A")
	lblBalance.SetText("Balance: N/A")
	visbilityButton.Disable()
	entryButton.Disable()
	searchEntry.Disable()
	logoutButton.Disable()
	refreshButton.Disable()
	if err := walletConnection("", "", true); err != nil {
		panic(err)
	}
	contentContainer.Refresh()
	resetCh = make(chan struct{})
}

func (theme1) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return color.Black
	default:
		return theme.DefaultTheme().Color(c, v)
	}
}

func (theme1) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return resourceRobotoRegularTtf
	}
	if s.Bold {
		if s.Italic {
			return resourceRobotoBoldItalicTtf
		}
		return resourceRobotoBoldTtf
	}
	if s.Italic {
		return resourceRobotoItalicTtf
	}
	return resourceRobotoRegularTtf
}

func (theme1) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (theme1) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	default:
		return theme.DefaultTheme().Size(s)
	}
}

func filterOptions(originalToTruncated map[string]string, text string) []string {
	var filteredOptions []string
	for truncatedOption, _ := range originalToTruncated {
		if strings.Contains(strings.ToLower(truncatedOption), strings.ToLower(text)) {
			filteredOptions = append(filteredOptions, truncatedOption)
		}
	}
	return filteredOptions
}
