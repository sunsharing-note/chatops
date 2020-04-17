package sshd
import (
	gossh "golang.org/x/crypto/ssh"
	"net"
)

type Cli struct {
	user    string
	pwd     string
	addr    string
	client  *gossh.Client
	session *gossh.Session
	LastResult string
}

func NewSSH(user,pwd,addr string) *Cli{
	return &Cli{
		user:       user,
		pwd:        pwd,
		addr:       addr,
	}
}

func (c *Cli) Connect() (*Cli, error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []gossh.AuthMethod{gossh.Password(c.pwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	client, err := gossh.Dial("tcp", c.addr, config)

	if nil != err {
		return c, err
	}
	c.client = client
	return c, nil
}

func (c Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if _,err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

