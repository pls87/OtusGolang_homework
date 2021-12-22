package main

import (
	"bufio"
	"bytes"
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
	address          string
	messages2Send    []string
	messages2Receive []string
}

type telnetTestSuite struct {
	suite.Suite
	params   caseParams
	listener net.Listener
	in       *bytes.Buffer
	out      *bytes.Buffer

	client TelnetClient
}

func (s *telnetTestSuite) TearDownTest() {
	s.NoError(s.listener.Close())
	s.NoError(s.client.Close())
}

func (s *telnetTestSuite) TestBasicCase() {
	s.params = caseParams{
		timeout:          10 * time.Second,
		address:          "127.0.0.1:",
		messages2Send:    []string{"hello"},
		messages2Receive: []string{"world"},
	}

	s.RunTest()
}

func (s *telnetTestSuite) initTest() {
	s.in = &bytes.Buffer{}
	s.out = &bytes.Buffer{}

	var err error
	s.listener, err = net.Listen("tcp", s.params.address)
	s.NoError(err)
	s.client = NewTelnetClient(s.listener.Addr().String(), s.params.timeout, ioutil.NopCloser(s.in), s.out)
	s.NoError(s.client.Connect())
}

func (s *telnetTestSuite) sendMessages(messages []string, w io.Writer) {
	for _, m := range messages {
		w.Write([]byte(m + "\n"))
	}
}

func (s *telnetTestSuite) checkMessages(expected []string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for i := 0; i < len(expected); i++ {
		scanner.Scan()
		s.Equal(expected[i], scanner.Text())
	}
	s.NoError(scanner.Err())
}

func (s *telnetTestSuite) RunTest() {
	s.initTest()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.sendMessages(s.params.messages2Send, s.in)

		s.NoError(s.client.Send())
		s.NoError(s.client.Receive())

		s.checkMessages(s.params.messages2Receive, s.out)
	}()

	go func() {
		defer wg.Done()
		conn, err := s.listener.Accept()
		s.NoError(err)
		s.NotNil(conn)
		defer func() { s.NoError(conn.Close()) }()

		s.checkMessages(s.params.messages2Send, conn)
		s.sendMessages(s.params.messages2Receive, conn)
	}()

	wg.Wait()
}

func TestTelnetClient(t *testing.T) {
	suite.Run(t, new(telnetTestSuite))
}
