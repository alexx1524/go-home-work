package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("multiple connections", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			log.Fatalln(err)
		}

		require.NoError(t, err)
		defer func() { require.NoError(t, listener.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client1 := NewTelnetClient(listener.Addr().String(), time.Second*10, ioutil.NopCloser(in), out)
		err1 := client1.Connect()

		client2 := NewTelnetClient(listener.Addr().String(), time.Second*10, ioutil.NopCloser(in), out)
		err2 := client2.Connect()

		require.NoError(t, err1)
		require.NoError(t, err2)

		err1 = client1.Close()
		err2 = client2.Close()

		require.NoError(t, err1)
		require.NoError(t, err2)
	})

	t.Run("timeout", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			log.Fatalln(err)
		}

		require.NoError(t, err)
		defer func() { require.NoError(t, listener.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(listener.Addr().String(), time.Nanosecond, ioutil.NopCloser(in), out)
		err = client.Connect()

		require.Error(t, err)
	})

	t.Run("connection refused", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:2222", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Connect()

		var connRefusedError *os.SyscallError

		require.Error(t, err)
		require.True(t, errors.As(err, &connRefusedError))
		require.True(t, errors.Is(connRefusedError.Err, syscall.ECONNREFUSED))
	})
}
