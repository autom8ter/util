package netutil

import (
	"github.com/autom8ter/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	req *http.Request
	*http.Client
	root *cobra.Command
}

func NewClient(name, short, long string, u *url.URL, method string) *Client {
	var r = &http.Request{
		URL:    u,
		Method: method,
	}
	return &Client{
		req:    r,
		Client: http.DefaultClient,
		root: &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
		},
	}
}

func (c *Client) RequestHeaders(headers map[string]string, r *http.Request) *http.Request {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}

func (c *Client) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
}

func (c *Client) JSONify(obj interface{}) []byte {
	return util.ToPrettyJson(obj)
}

func (c *Client) ReqPOST(r *http.Request) *http.Request {
	r.Method = "POST"
	return r
}

func (c *Client) ReqGET(r *http.Request) *http.Request {
	r.Method = "GET"
	return r
}
func (c *Client) ReqURL(r *http.Request, url *url.URL) *http.Request {
	r.URL = url
	return r
}

func (c *Client) AsCsv(s string) ([]string, error) {
	return util.ReadAsCSV(s)
}

func (c *Client) WriteTo(r *http.Request, w io.Writer) *http.Request {
	r.Write(w)
	return r
}

func (c *Client) RequestBasicAuth(userName, password string, r *http.Request) *http.Request {
	r.SetBasicAuth(userName, password)
	return r
}

func (c *Client) Prompt(q string) string {
	return util.Prompt(q)
}

func (c *Client) AddCommands(cmds ...*cobra.Command) {
	c.root.AddCommand(cmds...)
}

func (c *Client) AddClients(clients ...*Client) {
	for _, cmd := range clients {
		c.root.AddCommand(cmd.root)
	}
}

func (c *Client) Execute() error {
	return c.root.Execute()
}

func (c *Client) RunFunc(clientFunc ClientFunc) {
	c.root.Run = func(cmd *cobra.Command, args []string) {
		cmd = c.root
		clientFunc(c, args)
	}
}

func (c *Client) Debug() {
	c.root.DebugFlags()
}

func (c *Client) PersistentFlags() *pflag.FlagSet {
	return c.root.PersistentFlags()
}

func (c *Client) Flags() *pflag.FlagSet {
	return c.root.Flags()
}

func (c *Client) Version(version string) {
	c.root.Version = version
}

func (c *Client) SetOutput(writer io.Writer) {
	c.root.SetOutput(writer)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func (c *Client) GenerateJWT(signingKey string, claims map[string]interface{}) (string, error) {
	return util.GenerateJWT(signingKey, claims)
}

func (c *Client) Init(cfgFile string, envPrefix string, contentType []string, authorization []string, origin []string, cookie []string) {

	c.root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "relative path to config file")
	c.root.PersistentFlags().StringSliceVar(&contentType, "content-type", []string{}, "request content-type headers")
	c.root.PersistentFlags().StringSliceVar(&authorization, "authorization", []string{}, "request authorization headers")
	c.root.PersistentFlags().StringSliceVar(&origin, "origin", []string{}, "request origin headers")
	c.root.PersistentFlags().StringSliceVar(&cookie, "cookie", []string{}, "request cookie headers")

	c.req.Header = http.Header{
		"Content-Type":  contentType,
		"Authorization": authorization,
		"Origin":        origin,
		"Cookie":        cookie,
	}
	viper.BindPFlags(c.root.PersistentFlags())
	viper.BindPFlags(c.root.Flags())

}

type ReqFunc func(r *http.Request) *http.Request
type ClientFunc func(c *Client, args []string)
