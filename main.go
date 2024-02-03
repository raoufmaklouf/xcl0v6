package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			Port, host, path, err := parseURL(u)
			if err == nil {
				port, _ := strconv.Atoi(Port)
				for x, payload := range payload_list {
					r1 := Request1(path, host, payload, port)

					if len(r1) > 1 {
						_, b1, err := splitHTTPResponse(r1)
						if err == nil {

							if err == nil {

								r2 := Request2(host, port)
								if len(r2) > 1 {
									h2, b2, err := splitHTTPResponse(r2)
									if err == nil {

										scode2, err := extractStatusCode(h2)
										if err == nil {
											r3 := Request3(host, port)
											if len(r3) > 1 {
												h3, b3, err := splitHTTPResponse(r3)
												scode3, err := extractStatusCode(h3)
												if err == nil {
													if b2 == b3 && scode2 == scode3 && b1 != b2 {
														//fmt.Println("\n======================================R1=========================================\n", r1)
														//fmt.Println("\n======================================R2=========================================\n", r2)
														//fmt.Println("\n======================================R3=========================================\n", r3)
														fmt.Println("\033[31m", u, "is vuln ", x, payload, "\033[0m")
														break
													}
												}

											}

										}
									}

								}
							}

						}

					}

				}

			}
		}(url)
		wg.Wait()
	}

}
