/*
MIT License

Copyright (c) 2016 Hackzzila

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

// Import packages
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ini/ini"
)

// Define global variables.
var track string = "Nothing"
var ver string = "2"
var dg *discordgo.Session
var err error

func handler(w http.ResponseWriter, r *http.Request) {
	// Print to page
	fmt.Fprintf(w, "  ___  _                   _ __  __         _    \n")
	fmt.Fprintf(w, " |   \\(_)___ __ ___ _ _ __| |  \\/  |_  _ __(_)__ \n")
	fmt.Fprintf(w, " | |) | (_-</ _/ _ \\ '_/ _` | |\\/| | || (_-< / _|\n")
	fmt.Fprintf(w, " |___/|_/__/\\__\\___/_| \\__,_|_|  |_|\\_,_/__/_\\__|\n")
	fmt.Fprintf(w, "    Version: "+ver)
	fmt.Fprintf(w, "\nNow Playing: "+track)
	fmt.Println(track)
}

func newHttp() {
	//Bind and start server
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}

func statusLoop() {
	// Open Snip.txt
	file, err := ioutil.ReadFile("snip/Snip.txt")
	if err != nil {
		panic(err)
	}

	// Set status
	if string(file) != "" {
		dg.UpdateStatus(0, "🎧 "+string(file))
		fmt.Println("Now Playing: " + string(file))
		track = string(file)
	} else {
		dg.UpdateStatus(0, "🎧 Nothing")
		fmt.Println("Now Playing: Nothing")
		track = "Nothing"
	}

	for {
		// Open Snip.txt
		new, err := ioutil.ReadFile("snip/Snip.txt")
		if err != nil {
			panic(err)
		}

		if string(file) != string(new) {
			track = string(new)
			if track == "" {
				track = "Nothing"
			}
			dg.UpdateStatus(0, ":headphones: "+track)
			fmt.Println("Now Playing: " + track)
		}
		// Sleep for 5 seconds
		time.Sleep(5000000000)
	}
}

func main() {
	// Start webserver
	go newHttp()

	// Print layout
	fmt.Println("  ___  _                   _ __  __         _    ")
	fmt.Println(" |   \\(_)___ __ ___ _ _ __| |  \\/  |_  _ __(_)__ ")
	fmt.Println(" | |) | (_-</ _/ _ \\ '_/ _` | |\\/| | || (_-< / _|")
	fmt.Println(" |___/|_/__/\\__\\___/_| \\__,_|_|  |_|\\_,_/__/_\\__|")
	fmt.Println("Version " + ver)
	fmt.Println("Serving at http://localhost:5000")
	fmt.Println("")
	fmt.Println("___MUSIC___")

	// Start Snip
	cmd := exec.Command("snip/Snip.exe")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// Load the config file
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the email and password from the config file
	email := cfg.Section("Credentials").Key("email").String()
	password := cfg.Section("Credentials").Key("password").String()

	// Call the helper function New() passing username and password command
	// line arguments. This returns a new Discord session, authenticates,
	// connects to the Discord data websocket, and listens for events.
	dg, err = discordgo.New(email, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open the websocket and begin listening.
	dg.Open()

	// Start the status loop.
	go statusLoop()

	// Keep script open
	<-make(chan struct{})
}
