package goscp

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	hostKey        ssh.PublicKey
	defaultKeyPath = path.Join(userHome(), ".ssh/id_rsa")
)

// ScpOptions ID
type ScpOptions struct {
	Hostname    string
	Port        string `json:"port,omitempty"`
	Username    string
	KeyPath     string `json:"keyPath,omitempty"`
	Password    string `json:"password,omitempty"`
	Src         string
	Dst         string
	Connections int `json:"connections,omitempty"`
}

func userHome() string {
	usr, err := user.Current()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get user.")
	}

	return usr.HomeDir
}

// Scp ID
func Scp(op ScpOptions) error {
	// should establish a connection & copy file to remote

	// check src exists
	if _, err := os.Stat(op.Src); err != nil {
		logrus.WithError(err).Fatal("source does not exist")
	}

	var config *ssh.ClientConfig
	if op.Password != "" {
		config = &ssh.ClientConfig{
			User: op.Username,
			Auth: []ssh.AuthMethod{
				// Use the password for remote authentication.
				ssh.Password(op.Password),
			},
			// HostKeyCallback: ssh.FixedHostKey(hostKey),
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		if op.KeyPath == "" {
			op.KeyPath = defaultKeyPath
		}

		// Parse keys
		key, err := ioutil.ReadFile(op.KeyPath)
		if err != nil {
			logrus.WithError(err).Fatal("unable to read private key")
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logrus.WithError(err).Fatal("unable to parse private key")
		}

		config = &ssh.ClientConfig{
			User: op.Username,
			Auth: []ssh.AuthMethod{
				// Use the PublicKeys method for remote authentication.
				ssh.PublicKeys(signer),
			},
			// HostKeyCallback: ssh.FixedHostKey(hostKey),
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	client, err := ssh.Dial("tcp", op.Hostname+":"+op.Port, config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to dial")
	}

	handlerClient := Client{client: client, MaxConnections: op.Connections}
	return handlerClient.copy(op.Src, op.Dst)
}
