package netutil

import (
	"github.com/autom8ter/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"net/http"
)

type Client struct {
	*http.Client
	root *cobra.Command
}

func NewClient(name, short, long string) *Client {
	return &Client{
		Client: http.DefaultClient,
		root: &cobra.Command{
			Use: name,
			Short: short,
			Long: long,
		},
	}
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

func (c *Client) Version(version string)  {
	c.root.Version = version
}

func (c *Client) SetOutput(writer io.Writer)  {
	c.root.SetOutput(writer)
}

func (c *Client) Do(reqfns ...ReqFunc) (*http.Response, error){
	req := &http.Request{}
	for _, r := range reqfns {
		r(req)
	}
	return c.Client.Do(req)
}

func (c *Client) GenerateJWT(signingKey string, claims map[string]interface{}) string {
	return ""
}
type ReqFunc func(r *http.Request)
type ClientFunc func(c *Client, args[]string)