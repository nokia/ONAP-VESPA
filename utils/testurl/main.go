/*
	Copyright 2019 Nokia

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
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	certFile    = flag.String("cert", "", "Path to PEM certificate file")
	endpointURL = flag.String("url", "https://localhost/", "URL to test")
	insecure    = flag.Bool("insecure", false, "Ignore untrusted certificates")
)

func main() {
	flag.Parse()

	var tlsConfig *tls.Config
	if *certFile != "" || *insecure {
		/* #nosec */
		tlsConfig = &tls.Config{
			InsecureSkipVerify: *insecure,
		}
		if *certFile != "" {
			rootCa := x509.NewCertPool()
			caBytes, err := ioutil.ReadFile(*certFile)
			if err != nil {
				log.Fatalln("Cannot load CA certificate file:", err.Error())
			}
			if !rootCa.AppendCertsFromPEM(caBytes) {
				log.Fatalln("Cannot load root CA. PEM not valid")
			}
			tlsConfig.RootCAs = rootCa
		}
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
	}

	resp, err := client.Get(*endpointURL)
	if err != nil {
		log.Fatalln("Cannot connect:", err.Error())
	}
	if err = resp.Body.Close(); err != nil {
		log.Fatalln("Cannot close response body stream:", err.Error())
	}
}
