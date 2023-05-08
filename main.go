// netcat go version
// Usage:
//
//	go-nc hostname port
//
// - connect to hostname:port
// - read stdin and send to socket
// - read socket and write to stdout
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeout time.Duration
)

type deadlineReader struct {
	conn     net.Conn
	deadline time.Duration
}

type deadlineWriter struct {
	conn     net.Conn
	deadline time.Duration
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (r *deadlineReader) Read(p []byte) (n int, err error) {
	must(r.conn.SetReadDeadline(time.Now().Add(r.deadline)))

	return r.conn.Read(p)
}

func (r *deadlineWriter) Write(p []byte) (n int, err error) {
	must(r.conn.SetWriteDeadline(time.Now().Add(r.deadline)))

	return r.conn.Write(p)
}

func proxy(conn net.Conn) {
	var wg sync.WaitGroup

	for _, rw := range [...]struct {
		reader io.Reader
		writer io.Writer
	}{
		{
			os.Stdin,
			&deadlineWriter{conn, timeout},
		},
		{
			&deadlineReader{conn, timeout},
			os.Stdout,
		},
	} {
		wg.Add(1)
		go func(src io.Reader, dst io.Writer) {
			defer wg.Done()

			io.Copy(dst, src)
		}(rw.reader, rw.writer)
	}

	wg.Wait()
}

func main() {
	root := &cobra.Command{
		Use:     "gn hostname port",
		Short:   "gn",
		Long:    "gn",
		Example: "gn target.host 22",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			hostname := args[0]
			port := args[1]

			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port))
			if err != nil {
				return err
			}
			defer conn.Close()

			proxy(conn)
			return nil
		},
	}
	fs := root.Flags()
	fs.DurationVarP(&timeout, "timeout", "w", time.Minute,
		"If a connection or stdin are idle for more than timeout seconds, then the connection is silently closed.")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
