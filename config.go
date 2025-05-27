package eastmoney

import "time"

type config struct {
	timeout          time.Duration
	disableKeepAlive bool
	proxy            string
	ua               string
	preCall          func(*Client) error
}

type configFn func(*config)

// DisableKeepAlive disables HTTP keep-alive connections.
func DisableKeepAlive() configFn {
	return func(c *config) {
		c.disableKeepAlive = true
	}
}

// Proxy sets the HTTP proxy for the client.
func Proxy(proxy string) configFn {
	return func(c *config) {
		c.proxy = proxy
	}
}

// UserAgent sets the User-Agent header for HTTP requests.
func UserAgent(ua string) configFn {
	return func(c *config) {
		c.ua = ua
	}
}

// PreCall sets a function to be called before each API call.
func PreCall(fn func(*Client) error) configFn {
	return func(c *config) {
		c.preCall = fn
	}
}

// Timeout sets the timeout for HTTP requests.
func Timeout(d time.Duration) configFn {
	return func(c *config) {
		c.timeout = d
	}
}
