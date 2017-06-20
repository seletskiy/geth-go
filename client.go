package geth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
)

const (
	// BlockLatest represents latest available block.
	BlockLatest = "latest"
)

// Client connects to geth daemon and provides API for JSON RPC.
type Client struct {
	Host string
	Port int
}

// NewClient returns new client for connection to specified host and port.
func NewClient(host string, port int) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}

// NewLocalClient returns new client for localhost:8545, which is default
// geth listen address.
func NewLocalClient() *Client {
	return NewClient("localhost", 8545)
}

// GetBalance returns balance for specified address. Note, that address should
// begin with 0x prefix.
func (client *Client) GetBalance(address string) (*big.Int, error) {
	reply, err := client.Call("eth_getBalance", address, BlockLatest)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to get address balance "%s": %s`,
			address,
			err,
		)
	}

	var envelope string

	err = json.Unmarshal(reply, &envelope)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal balance envelope: %s", err,
		)
	}

	var result big.Int

	err = DecodeHex(envelope, &result)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to decode hex value of balance: %s", err,
		)
	}

	return &result, nil
}

// Call is a generic call to JSON RPC of geth daemon.
func (client *Client) Call(
	method string,
	args ...interface{},
) (json.RawMessage, error) {
	var request = struct {
		Header string `json:"jsonrpc"`
		Method string
		Params []interface{}
		ID     int
	}{
		Header: "2.0",
		Method: method,
		Params: args,
		ID:     rand.Int(),
	}

	packet, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to marshal request (method "%s"): %s`,
			method,
			err,
		)
	}

	response, err := http.Post(
		fmt.Sprintf("http://%s:%d", client.Host, client.Port),
		"application/json",
		bytes.NewBuffer(packet),
	)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to make post request to "%s:%d": %s`,
			client.Host,
			client.Port,
			err,
		)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf(
			`unexpected non-200 status code from "%s:%d": %s`,
			client.Host,
			client.Port,
			response.StatusCode,
		)
	}

	var reply struct {
		Error struct {
			Code    int
			Message string
		}

		Result json.RawMessage
	}

	err = json.NewDecoder(response.Body).Decode(&reply)
	if err != nil {
		return nil, fmt.Errorf(
			`unable to unmarshal response from "%s:%d": %s`,
			client.Host,
			client.Port,
			err,
		)
	}

	if reply.Error.Code != 0 {
		return nil, fmt.Errorf(
			"error while processing request (code %d): %s",
			reply.Error.Code,
			reply.Error.Message,
		)
	}

	return reply.Result, nil
}
