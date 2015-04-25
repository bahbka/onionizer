# onionizer
Onionize any web site without modifying its code

You can make [Tor Hidden Service](https://www.torproject.org/docs/tor-hidden-service.html) with any web site without changing original site code. This is very simple proxy which receives request, replaces onion hostname with original hostname, sends modified request to original site and replace original hostname with onion at response to client. This happens also with headers (cookies, redirects, referrers...). Writen in [Go](https://golang.org/) with [goproxy](https://github.com/elazarl/goproxy).

```
Usage:
  ./onionizer [OPTIONS]

Options:
  -http_addr <ADDR:PORT> proxy listen address (default :8080)
  -https_addr <ADDR:PORT> proxy https listen address (default :8081)

  -cert <FILE> certificate for https (default cert.pem)
  -key <FILE> private key for https (default key.pem)

  -origin <DOMAIN> original domain name (default example.com)
  -onion <DOMAIN> onion domain name (default example.onion)

  -server <ADDR> proxy all requests to (default empty, proxy to origin)

  -verbose every proxy request will be logged to STDOUT
```
