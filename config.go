package eastmoney

type config struct {
	disableKeepAlive bool
	proxy            string
	ua               string
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
