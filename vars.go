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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/blang/semver/v4"
	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc"
)

var (
	// app
	resetCh = make(chan struct{})
	err     error
	session appSession
	ui      appUI

	version        = semver.MustParse("0.0.1")
	versionMsg     = "secret-journal | version: %s \n"
	copyrightMsg   = "Copyright 2024 secretnamebasis. All rights reserved."
	osArchGoMaxMsg = "OS: %s ARCH: %s GOMAXPROCS: %d\n\n"
	quoteMsg       = "Imitation is the sincerest form of flattery that mediocrity can pay to greatness. ― Oscar Wilde"
	transfersMap   = make(map[string]string)
	sortedTimes    []string
	// domains
	appLanding = "app.main.landing"
	appMain    = "app.main"

	// developer support
	tipMsg = "secret-journal support"
	tipAmt = uint64(200)

	// dero
	deroUsername   string
	deroPassword   string
	deroAddress    = DEVELOPER_ADDRESS
	deroIp         = "127.0.0.1"
	deroPort       = "10103"
	deroEndpoint   = "http://" + deroIp + ":" + deroPort + "/json_rpc"
	deroHttpClient *http.Client
	deroRpcClient  jsonrpc.RPCClient
	addr_result    rpc.GetAddress_Result
	address        *rpc.Address
	balance        rpc.GetBalance_Result
	transfers      rpc.Get_Transfers_Result
	walletHeight   *rpc.GetHeight_Result
	payload        = rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    "",
		},
		{
			Name:     rpc.RPC_REPLYBACK_ADDRESS,
			DataType: rpc.DataAddress,
			Value:    "",
		},
	}

	transfer = rpc.Transfer{
		Destination: DEVELOPER_ADDRESS,
		Amount:      uint64(0),
		Payload_RPC: payload,
	}

	receipt = rpc.Arguments{
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    tipMsg,
		},
	}

	tip = rpc.Transfer{
		Destination: DEVELOPER_ADDRESS,
		Amount:      tipAmt,
		Payload_RPC: receipt,
	}

	transferParams = rpc.Transfer_Params{
		Transfers: []rpc.Transfer{transfer, tip},
	}

	// fyne
	themes                 appThemes
	modal                  *widget.PopUp
	visbilityButton        *widget.Button
	entryButton            *widget.Button
	logoutButton           *widget.Button
	refreshButton          *widget.Button
	isVisibilityOn         bool
	bottomButtonsContainer *fyne.Container
	entryContainer         *fyne.Container
	entryForm              *widget.Entry
	scrollContainer        *container.Scroll
	searchEntry            *widget.Entry
	lblHeight              = widget.NewLabel("Height: N/A")
	lblAddress             = widget.NewLabel("Wallet: N/A")
	lblBalance             = widget.NewLabel("Balance: N/A")
	contentContainer       = container.NewVBox()
	padding                = layout.NewSpacer()
)
