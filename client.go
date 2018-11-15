package goscp

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	defaultConnections        = 4
	SCP_PUSH_BEGIN_FILE       = "C"
	SCP_PUSH_BEGIN_FOLDER     = "D"
	SCP_PUSH_BEGIN_END_FOLDER = " 0"
	SCP_PUSH_END_FOLDER       = "E"
	SCP_PUSH_END              = "\x00"
)

type Client struct {
	client         *ssh.Client
	MaxConnections int
}

func (c *Client) copy(src, dst string) error {
	if c.MaxConnections == 0 {
		c.MaxConnections = defaultConnections
	}
	logrus.Info("performing copy with ", c.MaxConnections, " connections")

	// limit number of connections
	sessions := make(chan struct{}, c.MaxConnections)
	var wg sync.WaitGroup

	wg.Add(1)
	c.recursiveCopy(wg, sessions, src, dst)
	wg.Wait()

	return nil
}

func (c *Client) recursiveCopy(wg sync.WaitGroup, sessions chan struct{}, src, dst string) {
	defer wg.Done()

	srcinfo, err := os.Stat(src)
	if err != nil {
		logrus.WithError(err).Fatal("file/folder not found")
	}
	// get permission for a session otherwise wait.
	sessions <- struct{}{}
	if srcinfo.IsDir() {
		dst = path.Join(dst, srcinfo.Name())
		if err = c.makeDir(dst); err != nil {
			logrus.WithError(err).Fatal("failed to create directory on host")
		}

		// release a session
		<-sessions

		fds, err := ioutil.ReadDir(src)
		if err != nil {
			logrus.WithError(err).Fatal("error reading directory")
		}

		for _, fd := range fds {
			srcfp := path.Join(src, fd.Name())
			dstfp := path.Join(dst, fd.Name())
			wg.Add(1)
			go c.recursiveCopy(wg, sessions, srcfp, dstfp)
		}
	} else {
		logrus.Info("copying ", src, " to ", dst)
		c.pushFile(src, dst)
		logrus.Info("done copying ", src, " to ", dst)
		// release a session
		<-sessions
	}
}

// makeDir creates a dir on the remote
func (c *Client) makeDir(dst string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()
	if err := session.Run("/usr/bin/mkdir " + dst); err != nil {
		return err
	}

	return nil
}

// following code is based on https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
// Code is modified version of https://github.com/gnicod/goscplib

// pushFile pushes one file to server
func (c *Client) pushFile(src string, dst string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		prepareFile(w, src)
	}()
	if err := session.Run("/usr/bin/scp -tr " + filepath.Dir(dst)); err != nil {
		return err
	}

	return nil
}

func prepareFile(w io.WriteCloser, src string) {
	file, err := os.Open(src)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open file")
	}

	defer file.Close()
	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		logrus.WithError(err).Fatal("failed to stat file")
	}

	// Print the file content
	mode := SCP_PUSH_BEGIN_FILE + getPerm(fileInfo)
	fmt.Fprintln(w, mode, fileInfo.Size(), filepath.Base(src))
	io.Copy(w, file)
	fmt.Fprint(w, SCP_PUSH_END)
}

func getPerm(f os.FileInfo) (perm string) {
	mod := f.Mode().Perm()
	return fmt.Sprintf("%04o", uint32(mod))
}
