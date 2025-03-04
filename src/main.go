/*
check_ssl_expiration

Copyright 2024 Rodolfo González González.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/xorpaul/go-nagios"
)

//-----------------------------------------------------------------------------
// Flags

var (
	host    string
	port    int
	warn    int
	crit    int
	version bool
)

//-----------------------------------------------------------------------------

func init() {
	pflag.StringVarP(&host, "host", "H", "localhost", "The hostname from where to get the SSL certificate")
	pflag.IntVarP(&port, "port", "p", 443, "The URL from where to get the SSL certificate")
	pflag.IntVarP(&warn, "warn", "w", 15, "How many days til expiration constitutes a WARNING?")
	pflag.IntVarP(&crit, "crit", "c", 7, "How many days til expiration constitutes a CRITICAL alert?")
	pflag.BoolVarP(&version, "version", "v", false, "Show version number")
}

//-----------------------------------------------------------------------------

// getHostWithPort takes a URL string, parses it, and returns the host with ":443" appended
func getHostWithPort(rawURL string, port int) (string, error) {
	// Check if the protocol is missing and prepend "https://" if necessary
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Get the host part and append the port
	var hostWithPort string
	if parsedURL.Port() == "" {
		hostWithPort = parsedURL.Hostname() + ":" + strconv.Itoa(port)
	} else {
		hostWithPort = parsedURL.Hostname() + ":" + parsedURL.Port()
	}

	return hostWithPort, nil
}

//-----------------------------------------------------------------------------

func getSSLCertificate(url string) (*x509.Certificate, error) {
	// Dial a TLS connection
	conn, err := tls.Dial("tcp", url, &tls.Config{
		InsecureSkipVerify: true, // Don't verify the certificate for simplicity
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Get the connection's state to extract certificates
	state := conn.ConnectionState()

	// Get the first certificate in the chain (this is the server certificate)
	if len(state.PeerCertificates) > 0 {
		return state.PeerCertificates[0], nil
	}

	return nil, fmt.Errorf("no certificates found")
}

//-----------------------------------------------------------------------------

func getDaysUntilExpiry(cert *x509.Certificate) int {
	// Convert the duration to hours and then to days
	due := time.Until(cert.NotAfter).Hours() / 24

	return int(due)
}

//-----------------------------------------------------------------------------

func main() {
	var nr nagios.NagiosResult

	pflag.Parse()

	if version {
		fmt.Println("check_ssl_expiration v0.2")
		os.Exit(0)
	}

	netUrl, err := getHostWithPort(host, port)
	if err != nil {
		nr = nagios.NagiosResult{ExitCode: 3, Text: err.Error(), Perfdata: ""}
		fmt.Println("Error:", err)
		nagios.NagiosExit(nr)
	}

	cert, err := getSSLCertificate(netUrl)
	if err != nil {
		nr = nagios.NagiosResult{ExitCode: 3, Text: err.Error(), Perfdata: ""}
		fmt.Println("Error:", err)
		nagios.NagiosExit(nr)
	}

	due := getDaysUntilExpiry(cert)
	exp := cert.NotAfter.Local().Format(time.RFC1123)

	if due < 0 {
		nr = nagios.NagiosResult{ExitCode: 2, Text: fmt.Sprintf("the SSL certificate has already expired %d days ago (%s)", due*-1, exp), Perfdata: ""}
	} else if due <= warn && due > crit {
		nr = nagios.NagiosResult{ExitCode: 1, Text: fmt.Sprintf("the SSL certificate will expire in %d days (%s)", due, exp), Perfdata: ""}
	} else if due <= crit {
		nr = nagios.NagiosResult{ExitCode: 2, Text: fmt.Sprintf("the SSL certificate will expire in %d days (%s)", due, exp), Perfdata: ""}
	}

	nagios.NagiosExit(nr)
}
