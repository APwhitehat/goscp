package goscp

import "golang.org/x/crypto/ssh"

const (
	defaultConnections = 4
)

type Client struct {
	client         *ssh.Client
	MaxConnections int
}

func (c *Client) copy(src, dsc string) {
	if c != 0 {
		c.MaxConnections = defaultConnections
	}

	// // Each ClientConn can support multiple interactive sessions,
	// // represented by a Session.
	// session, err := client.NewSession()
	// if err != nil {
	//     log.Fatal("Failed to create session: ", err)
	// }
	// defer session.Close()

	// // Once a Session is created, you can execute a single command on
	// // the remote side using the Run method.
	// var b bytes.Buffer
	// session.Stdout = &b
	// if err := session.Run("/usr/bin/whoami"); err != nil {
	//     log.Fatal("Failed to run: " + err.Error())
	// }
	// fmt.Println(b.String())
}
