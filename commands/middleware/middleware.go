package middleware

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	apiClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	"github.com/netlify/open-api/go/porcelain"
	apiContext "github.com/netlify/open-api/go/porcelain/context"

	"github.com/aj0strow/netlifyctl/auth"
	"github.com/aj0strow/netlifyctl/context"
)

const defaultAPIPath = "/api/v1"

type CommandFunc func(context.Context, *cobra.Command, []string) error
type Middleware func(CommandFunc) CommandFunc

func NewRunFunc(f CommandFunc, mm []Middleware) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		runf := f
		for _, m := range mm {
			runf = m(runf)
		}

		return runf(ctx, cmd, args)
	}
}

func LoggingMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		ctx = apiContext.WithLogger(ctx, logrus.NewEntry(logrus.StandardLogger()))
		logrus.Debugf("setup logger middleware: %v", logrus.StandardLogger().Level)
		return cmd(ctx, c, args)
	}
}

func AuthMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		creds := auth.ClientCredentials()
		logrus.WithField("credentials", creds).Debug("setup credentials")

		ctx = apiContext.WithAuthInfo(ctx, creds)

		return cmd(ctx, c, args)
	}
}

func ClientMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		var transport *apiClient.Runtime

		if endpoint := c.Flag("endpoint"); endpoint != nil {
			if v := endpoint.Value.String(); v != "" {
				logrus.WithField("endpoint", v).Debug("setup API endpoint")

				u, err := url.Parse(v)
				if err != nil {
					return err
				}

				if u.Scheme == "" {
					u.Scheme = "http"
				}

				if u.Path == "" {
					u.Path = defaultAPIPath
				}

				transport = apiClient.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpClient())
			}
		}

		if transport == nil {
			logrus.WithField("endpoint", "https://api.netlify.com").Debug("setup default API endpoint")

			transport = apiClient.NewWithClient("api.netlify.com", "", []string{"https"}, httpClient())
		}

		client := porcelain.New(transport, strfmt.Default)
		ctx = context.WithClient(ctx, client)

		return cmd(ctx, c, args)
	}
}

func httpClient() *http.Client {
	protoUpgrade := map[string]func(string, *tls.Conn) http.RoundTripper{
		"ignore-h2": func(string, *tls.Conn) http.RoundTripper { return nil },
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSNextProto:          protoUpgrade,
	}

	return &http.Client{Transport: tr}
}
