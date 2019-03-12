package netutil

import (
	"github.com/autom8ter/util"
	"github.com/gorilla/sessions"
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

func (c *Client) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.req.Header.Set(k, v)
	}
}

func (c *Client) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
}

func (c *Client) JSONify(obj interface{}) []byte {
	return util.ToPrettyJson(obj)
}

func (c *Client) AsCsv(s string) ([]string, error) {
	return util.ReadAsCSV(s)
}

func (c *Client) WriteTo(w io.Writer) {
	c.req.Write(w)
}

func (c *Client) RequestBasicAuth(userName, password string) {
	c.req.SetBasicAuth(userName, password)
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

func (c *Client) Init(cfgFile string, envPrefix string, headers map[string]string) {
	if headers == nil {
		headers = make(map[string]string)
	}
	c.root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "relative path to config file")
	c.root.PersistentFlags().StringToStringVar(&headers, "headers", nil, "request headers")
	if viper.ConfigFileUsed() == "" {
		util.InitConfig(cfgFile, envPrefix)
	}
	viper.BindPFlags(c.root.PersistentFlags())
	viper.BindPFlags(c.root.Flags())
	c.SetHeaders(headers)
}

func (r *Client) Render(s string, data interface{}) string {
	return util.Render(s, data)
}

func (r *Client) NewSessionStore(key string) *sessions.CookieStore {
	return NewSessionStore(key)
}

type ClientFunc func(c *Client, args []string)
