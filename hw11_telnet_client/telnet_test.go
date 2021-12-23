package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type caseParams struct {
	timeout          time.Duration
	address2Listen   string
	address2Connect  string
	messages2Send    []string
	messages2Receive []string
	waitTimeoutError bool
}

type caseStatus struct {
	portOpened      bool
	clientConnected bool
	toServerR       *io.PipeReader
	fromClientW     *io.PipeWriter
	toClientR       *io.PipeReader
	fromServerW     *io.PipeWriter
	err             error
}

type telnetTestSuite struct {
	suite.Suite
	params   caseParams
	status   caseStatus
	listener net.Listener
	client   TelnetClient
}

func (s *telnetTestSuite) TearDownTest() {
	if s.status.portOpened {
		s.NoError(s.listener.Close())
	}
	if s.status.clientConnected {
		s.NoError(s.client.Close())
	}
}

func (s *telnetTestSuite) TestBasicCase() {
	s.params = caseParams{
		timeout:          10 * time.Second,
		messages2Send:    []string{"hello"},
		messages2Receive: []string{"world"},
	}

	s.initConnections()
	s.RunTest()
}

func (s *telnetTestSuite) TestSeveralMessagesCase() {
	s.params = caseParams{
		timeout:          10 * time.Second,
		messages2Send:    []string{"hello", "cruel", "world"},
		messages2Receive: []string{"bye-bye", "calm", "heaven"},
	}

	s.initConnections()
	s.RunTest()
}

func (s *telnetTestSuite) TestTimeoutCase() {
	s.params = caseParams{
		timeout:          5 * time.Second,
		address2Connect:  "127.0.0.2:5768",
		waitTimeoutError: true,
	}

	s.initConnections()
}

func (s *telnetTestSuite) RunTest() {
	conn, err := s.listener.Accept()
	s.NoError(err)
	s.NotNil(conn)
	defer func() { s.NoError(conn.Close()) }()

	var wg sync.WaitGroup
	wg.Add(4)

	go func() { // send messages from client
		defer wg.Done()
		for _, m := range s.params.messages2Send {
			s.status.fromClientW.Write([]byte(m + "\n"))
			time.Sleep(time.Second)
		}
	}()

	go func() { // read messages from client
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for i := 0; i < len(s.params.messages2Send); i++ {
			scanner.Scan()
			s.Equal(s.params.messages2Send[i], scanner.Text())
		}
		s.NoError(scanner.Err())
	}()

	go func() { // send messages from server
		defer wg.Done()
		for _, m := range s.params.messages2Receive {
			conn.Write([]byte(m + "\n"))
			time.Sleep(time.Second)
		}
	}()

	go func() { // read messages from server
		defer wg.Done()
		scanner := bufio.NewScanner(s.status.toClientR)
		for i := 0; i < len(s.params.messages2Receive); i++ {
			scanner.Scan()
			s.Equal(s.params.messages2Receive[i], scanner.Text())
		}
		s.NoError(scanner.Err())
	}()

	go func() {
		s.client.Send()
	}()

	go func() {
		s.client.Receive()
	}()

	wg.Wait()
}

func (s *telnetTestSuite) connectCheckStatus() {
	var wg sync.WaitGroup
	var completed int32
	wg.Add(1)
	go func(c *int32) {
		defer wg.Done()
		s.status.err = s.client.Connect()
		s.status.clientConnected = s.status.err == nil
		atomic.AddInt32(c, 1)
	}(&completed)

	if s.params.waitTimeoutError {
		s.Eventuallyf(func() bool {
			return !atomic.CompareAndSwapInt32(&completed, 1, 0)
		}, s.params.timeout-time.Second, s.params.timeout-2*time.Second, "")

		s.Eventuallyf(func() bool {
			return atomic.CompareAndSwapInt32(&completed, 1, 1)
		}, 4*time.Second, time.Second, "")

		wg.Wait()

		var result net.Error
		s.True(errors.As(s.status.err, &result))
		s.True(result.Timeout())

		return
	}

	wg.Wait()
	s.True(s.status.clientConnected)
}

func (s *telnetTestSuite) initConnections() {
	s.status = caseStatus{}

	s.status.toServerR, s.status.fromClientW = io.Pipe()
	s.status.toClientR, s.status.fromServerW = io.Pipe()

	address2Listen := s.params.address2Listen
	if address2Listen == "" {
		address2Listen = "127.0.0.1:"
	}

	var err error
	s.listener, err = net.Listen("tcp", address2Listen)
	s.NoError(err)
	s.status.portOpened = true

	address2Connect := s.params.address2Connect
	if address2Connect == "" {
		address2Connect = s.listener.Addr().String()
	}

	s.client = NewTelnetClient(address2Connect, s.params.timeout, s.status.toServerR, s.status.fromServerW)
	s.connectCheckStatus()
}

func TestTelnetClient(t *testing.T) {
	suite.Run(t, new(telnetTestSuite))
}
