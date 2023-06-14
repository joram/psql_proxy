package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/joram/psql_proxy/server/message"
	"io"
	"net"
	"time"
)

type AuthPhase int

const (
	PhaseStartup AuthPhase = iota
	PhaseGSS
	PhaseSASLInit
	PhaseSASL
	PhaseOK
)

const (
	BIND         = 'B'
	CLOSE        = 'C'
	COPYFAILED   = 'f'
	DESCRIBE     = 'D'
	EXECUTE      = 'E'
	FLUSH        = 'H'
	FUNCTIONCALL = 'F'
	PARSE        = 'P'
	QUERY        = 'Q'
	SYNC         = 'S'
	TERMINATE    = 'X'
	COPYDATA     = 'd'
	COPYDONE     = 'c'
	AUTH         = 'p'

	REJECT = 'N'
)

const InitMessageSizeLength = 4

type ClientHandler struct {
	ClientConn net.Conn
	AuthPhase  AuthPhase
	Context    context.Context
	Cancel     context.CancelFunc
}

func (c ClientHandler) handle() error {
	fmt.Println("handling new client connection")

	// if it's a TCP connection
	cc, ok := c.ClientConn.(*net.TCPConn)
	if ok {
		err := cc.SetKeepAlivePeriod(30 * time.Second)
		if err != nil {
			return err
		}
		err = cc.SetKeepAlive(true)
		if err != nil {
			return err
		}
	}

	fmt.Println("reading startup message")
	startup, err := c.readStartupMessage()
	if err != nil {
		io.Copy(c.ClientConn, message.ErrorResp("ERROR", "08P01", err.Error()).Reader())
		return err
	}
	_, isSSLRequest := startup.(*message.SSLRequest)
	if isSSLRequest {
		// TODO: SSLRequest, currently reject
		c.ClientConn.Write([]byte{REJECT})
		return nil
	}

	return nil
}

func (c ClientHandler) readStartupMessage() (message.Reader, error) {
	head := make([]byte, InitMessageSizeLength)
	if _, err := io.ReadFull(c.ClientConn, head); err != nil {
		return nil, err
	}

	data := make([]byte, binary.BigEndian.Uint32(head)-InitMessageSizeLength)
	if _, err := io.ReadFull(c.ClientConn, data); err != nil {
		return nil, err
	}

	msg := message.NewBaseFromBytes(data)
	code := msg.ReadUint32()
	if code == 80877103 {
		return &message.SSLRequest{RequestCode: code}, nil
	}

	if code == 80877102 {
		req := &message.CancelRequest{RequestCode: code}
		req.ProcessID = msg.ReadUint32()
		req.SecretKey = msg.ReadUint32()
		return req, nil
	}

	if majorVersion := code >> 16; majorVersion < 3 {
		return nil, errors.New("pg protocol < 3.0 is not supported")
	}

	startupMsg := &message.StartupMessage{ProtocolVersion: code, Parameters: make(map[string]string)}
	for {
		k := msg.ReadString()
		if k == "" {
			break
		}
		startupMsg.Parameters[k] = msg.ReadString()
	}
	fmt.Println("startup parameters", startupMsg.Parameters)
	fmt.Println("protocol version", startupMsg.ProtocolVersion)

	fmt.Println()
	return startupMsg, nil
}
