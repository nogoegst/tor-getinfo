// tor-getinfo.go - just issue GETINFO into stdout.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of tor-getinfo, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"github.com/yawning/bulb"
)

func main() {
	var debug_flag = flag.Bool("debug", false, "Show what's happening")
	var control = flag.String("control-port", "9051", "Set Tor ControlPort to be used")
	flag.Parse()
	debug := *debug_flag
	var tail = flag.Args()
	if len(tail) != 1 {
		log.Fatalf("You should specify exactly one keyword")
	}
	var keyword = tail[0]

	// Connect to a running tor instance.
	c, err := bulb.Dial("tcp4", "127.0.0.1:"+*control)
	//c, err := bulb.Dial("unix", "/var/run/tor/control")
	if err != nil {
		log.Fatalf("Failed to connect to control socket: %v", err)
	}
	defer c.Close()

	// See what's really going on under the hood.
	// Do not enable in production.
	c.Debug(debug)

	// Authenticate with the control port.  The password argument
	// here can be "" if no password is set (CookieAuth, no auth).
	if err := c.Authenticate("ExamplePassword"); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// At this point, c.Request() can be used to issue requests.

	// Request GETINFO with provided keyword
	resp, err := c.Request(fmt.Sprintf("GETINFO %v", keyword)) //
	if err != nil {
		log.Fatalf("GETINFO failed: %v", err)
	}
	// Extract real data without keyword prefix
	resp_data := strings.Replace(resp.Data[0], keyword+"=", "", 1)
	// Print to stdout
	fmt.Printf("%v", resp_data)
	c.Close()
}
