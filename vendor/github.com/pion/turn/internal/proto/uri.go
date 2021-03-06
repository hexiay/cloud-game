package proto

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Scheme definitions from RFC 7065 Section 3.2.
const (
	Scheme       = "turn"
	SchemeSecure = "turns"
)

// Transport definitions as in RFC 7065.
const (
	TransportTCP = "tcp"
	TransportUDP = "udp"
)

// URI as defined in RFC 7065.
type URI struct {
	Scheme    string
	Host      string
	Port      int
	Transport string
}

func (u URI) String() string {
	transportSuffix := ""
	if len(u.Transport) > 0 {
		transportSuffix = "?transport=" + u.Transport
	}
	if u.Port != 0 {
		return fmt.Sprintf("%s:%s:%d%s",
			u.Scheme, u.Host, u.Port, transportSuffix,
		)
	}
	return u.Scheme + ":" + u.Host + transportSuffix
}

// ParseURI parses URI from string.
func ParseURI(rawURI string) (URI, error) {
	// Carefully reusing URI parser from net/url.
	u, urlParseErr := url.Parse(rawURI)
	if urlParseErr != nil {
		return URI{}, urlParseErr
	}
	if u.Scheme != Scheme && u.Scheme != SchemeSecure {
		return URI{}, fmt.Errorf("unknown uri scheme %q", u.Scheme)
	}
	if u.Opaque == "" {
		return URI{}, errors.New("invalid uri format: expected opaque")
	}
	// Using URL methods to split host.
	u.Host = u.Opaque
	host, rawPort := u.Hostname(), u.Port()
	uri := URI{
		Scheme:    u.Scheme,
		Host:      host,
		Transport: u.Query().Get("transport"),
	}
	if len(rawPort) > 0 {
		port, portErr := strconv.Atoi(rawPort)
		if portErr != nil {
			return uri, fmt.Errorf("failed to parse %q as port: %v", rawPort, portErr)
		}
		uri.Port = port
	}
	return uri, nil
}
