package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"sync"
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
	in              *bytes.Buffer
	out             *bytes.Buffer
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
		timeout:          3 * time.Second,
		address2Connect:  "127.1.0.1:5768",
		waitTimeoutError: true,
	}

	s.initConnections()
}

func (s *telnetTestSuite) RunTest() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		s.sendMessagesFromClient()

		s.NoError(s.client.Send())
		s.NoError(s.client.Receive())

		s.checkMessagesToClient()
	}()

	go func() {
		defer wg.Done()
		conn, err := s.listener.Accept()
		s.NoError(err)
		s.NotNil(conn)
		defer func() { s.NoError(conn.Close()) }()

		s.checkMessagesFromClient(conn)
		s.sendMessagesToClient(conn)
	}()

	wg.Wait()
}

func (s *telnetTestSuite) connectCheckStatus() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.status.err = s.client.Connect()
		s.status.clientConnected = s.status.err == nil
	}()

	if s.params.waitTimeoutError {
		s.Eventuallyf(func() bool {
			return s.status.err == nil
		}, s.params.timeout, s.params.timeout-2*time.Second, "")

		s.Eventuallyf(func() bool {
			return s.status.err != nil
		}, 4*time.Second, time.Second, "")

		var result net.Error
		s.True(errors.As(s.status.err, &result))
		s.True(result.Timeout())
		return
	}

	wg.Wait()
	s.True(s.status.clientConnected)
}

func (s *telnetTestSuite) initConnections() {
	s.status = caseStatus{
		in:  &bytes.Buffer{},
		out: &bytes.Buffer{},
	}

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

	s.client = NewTelnetClient(address2Connect, s.params.timeout, ioutil.NopCloser(s.status.in), s.status.out)
	s.connectCheckStatus()
}

func (s *telnetTestSuite) sendMessagesFromClient() {
	s.sendMessages(s.params.messages2Send, s.status.in)
}

func (s *telnetTestSuite) sendMessagesToClient(conn net.Conn) {
	s.sendMessages(s.params.messages2Receive, conn)
}

func (s *telnetTestSuite) sendMessages(messages []string, w io.Writer) {
	for _, m := range messages {
		w.Write([]byte(m + "\n"))
	}
}

func (s *telnetTestSuite) checkMessagesFromClient(conn net.Conn) {
	s.checkMessages(s.params.messages2Send, conn)
}

func (s *telnetTestSuite) checkMessagesToClient() {
	s.checkMessages(s.params.messages2Receive, s.status.out)
}

func (s *telnetTestSuite) checkMessages(expected []string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for i := 0; i < len(expected); i++ {
		scanner.Scan()
		s.Equal(expected[i], scanner.Text())
	}
	s.NoError(scanner.Err())
}

func TestTelnetClient(t *testing.T) {
	suite.Run(t, new(telnetTestSuite))
}
