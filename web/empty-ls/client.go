package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func request(client *http.Client) {
	//r, err := client.Get("https://admin.zone443.dev/")
	r, err := client.Get("https://lrindels2.zone443.dev/")
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Println(string(body))
}

func readPem(path string) *x509.Certificate {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// convert to certificate
	block, _ := pem.Decode(raw)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return cert
}

func main() {
	// load our client certificate and key for user 'test'
	cert, err := tls.LoadX509KeyPair("certs/test.pem", "certs/test.key")
	if err != nil {
		log.Fatal(err)
	}

	admin_cert := readPem("certs/admin.pem")

	// set leaf certificate
	// cert.Leaf = admin_cert // fails, we don't have key for admin cert

	// add admin cert to chain
	// surprisingly this throws no error
	// server can see this in PeerCertificates[1]
	cert.Certificate = append(cert.Certificate, admin_cert.Raw)

	// swap leaf with next certificate in chain. Fails
	/*
		var tmp []byte
		tmp = cert.Certificate[0]
		cert.Certificate[0] = cert.Certificate[1]
		cert.Certificate[1] = tmp
	*/

	// Create a HTTPS client and supply the certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion:   tls.VersionTLS13, // test different TLS versions
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	request(client)
}
