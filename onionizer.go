package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	verbose := flag.Bool("verbose", false, "should every proxy request be logged to stdout")

	http_addr := flag.String("http_addr", ":8080", "proxy listen address")
	https_addr := flag.String("https_addr", ":8081", "proxy https listen address")

	cert_file := flag.String("cert", "cert.pem", "https certificate")
	key_file := flag.String("key", "key.pem", "https private key")

	origin := flag.String("origin", "example.com", "origin domain")
	onion := flag.String("onion", "example.onion", "onion domain")

	server := flag.String("server", "", "proxy requests to host (origin domain by default")

	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == "" {
			fmt.Fprintln(w, "Cannot handle requests without Host header, e.g., HTTP 1.0")
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = req.Host
		proxy.ServeHTTP(w, req)
	})

	proxy.OnResponse().DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			for key, value := range resp.Header {
				for index, _ := range value {
					resp.Header[key][index] = strings.Replace(value[index], *origin, *onion, -1)
				}
			}

			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			new_body := strings.Replace(string(body), *origin, *onion, -1)

			buf := bytes.NewBufferString(new_body)
			resp.Body = ioutil.NopCloser(buf)

			return resp
		})

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if *server == "" {
				r.URL.Host = strings.Replace(r.URL.Host, *onion, *origin, -1)
			} else {
				r.URL.Host = *server
			}

			for key, value := range r.Header {
				for index, _ := range value {
					r.Header[key][index] = strings.Replace(value[index], *onion, *origin, -1)
				}
			}

			for key, value := range r.Form {
				for index, _ := range value {
					r.Form[key][index] = strings.Replace(value[index], *onion, *origin, -1)
				}
			}

			for key, value := range r.PostForm {
				for index, _ := range value {
					r.PostForm[key][index] = strings.Replace(value[index], *onion, *origin, -1)
				}
			}

			r.Host = strings.Replace(r.Host, *onion, *origin, -1)

			return r, nil
		})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.Header.Add("X-Forwarded-Proto", "https")
		proxy.ServeHTTP(w, req)
	})

	go func() {
		log.Fatal(http.ListenAndServe(*http_addr, proxy))
	}()

	err := http.ListenAndServeTLS(*https_addr, *cert_file, *key_file, nil)
	if err != nil {
		log.Fatal(err)
	}
}
