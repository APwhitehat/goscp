package internal

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var hostKey ssh.PublicKey

type ScpOptions struct {
	Hostname string
	Port     int
	Username string
	KeyPath  string
	Password string
}

func userHome() string {
	usr, err := user.Current()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get user.")
	}

	return usr.HomeDir
}

// Scp insertInfo
func Scp(op ScpOptions, src, dsc string) {
	// should establish a connec & copy file to remote

	// check src exists
	if _, err := os.Stat(src); err != nil {
		logrus.WithError(err).Fatal("Sourse does not exist")
	}

	// Parse keys
	key, err := ioutil.ReadFile(path.Join(userHome(), ".ssh/id_rsa"))
	if err != nil {
		logrus.WithError(err).Fatal("unable to read private key")
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logrus.WithError(err).Fatal("unable to parse private key")
	}

	config := &ssh.ClientConfig{
		User: op.Username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", op.Hostname, config)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to dial")
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
