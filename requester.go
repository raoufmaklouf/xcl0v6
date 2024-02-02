package main

func Request1(path string, host string, payload string, port int) string {

	var responses1 string = ""

	// Create a TCP connection
	tcpConn, err := createTCPConnection(host, port)
	if err == nil {
		defer tcpConn.Close()
		tlsConn, err = createTLSConnection(tcpConn)
		if err == nil {
			defer tlsConn.Close()
			err = sendRequest("POST %s HTTP/1.1\r\nHost: %s\r\nConnection: keep-alive\r\nContent-Type: application/x-www-form-urlencoded\r\n%sGET /robots.txt HTTP/1.1\r\nFoo: x", path, host, payload)
			if err == nil {
				responsePrefix := "HTTP/1.1"
				responseCount := 2
				combinedResponse, err := readFullResponse(responsePrefix, responseCount)
				if err == nil {
					responses1 = combinedResponse

				}

			}

		}

	}

	return responses1
}

func Request2(host string, port int) string {

	var responses1 string = ""

	// Create a TCP connection
	tcpConn, err := createTCPConnection(host, port)
	if err == nil {
		defer tcpConn.Close()
		tlsConn, err = createTLSConnection(tcpConn)
		if err == nil {
			defer tlsConn.Close()
			err = sendRequest("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\n\r\n", host)
			if err == nil {
				responsePrefix := "HTTP/1.1"
				responseCount := 2
				combinedResponse, err := readFullResponse(responsePrefix, responseCount)
				if err == nil {
					responses1 = combinedResponse

				}

			}

		}

	}

	return responses1
}

func Request3(host string, port int) string {

	var responses1 string = ""

	// Create a TCP connection
	tcpConn, err := createTCPConnection(host, port)
	if err == nil {
		defer tcpConn.Close()
		tlsConn, err = createTLSConnection(tcpConn)
		if err == nil {
			defer tlsConn.Close()
			err = sendRequest("GET /robots.txt HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\n\r\n", host)
			if err == nil {
				responsePrefix := "HTTP/1.1"
				responseCount := 2
				combinedResponse, err := readFullResponse(responsePrefix, responseCount)
				if err == nil {
					responses1 = combinedResponse

				}

			}

		}

	}
	return responses1
}
