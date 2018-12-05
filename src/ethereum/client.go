package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"

	hierr "github.com/reconquest/hierr"
)

const (
	// BlockLatest represents latest available block.
	BlockLatest = "latest"
)

const (
	ErrorCodeAuthenticationNeeded = -32000
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
func (client *Client) GetBalance(address string) (*Wei, error) {
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

	return &Wei{result}, nil
}

// GetGasPrice returns current gas price in wei.
func (client *Client) GetGasPrice() (*Wei, error) {
	reply, err := client.Call("eth_gasPrice")
	if err != nil {
		return nil, fmt.Errorf(
			`unable to get gas price: %s`, err,
		)
	}

	var (
		result   big.Int
		envelope string
	)

	err = json.Unmarshal(reply, &envelope)
	if err != nil {
		return nil, err
	}

	err = DecodeHex(envelope, &result)
	if err != nil {
		return nil, err
	}

	return &Wei{result}, nil
}

// EstimateGas returns gas estimation for transaction.
func (client *Client) EstimateGas(
	from string,
	to string,
) (*Wei, error) {
	transaction := &Transaction{
		From: from,
		To:   to,
	}

	reply, err := client.Call("eth_estimateGas", transaction)
	if err != nil {
		return nil, err
	}

	var (
		result   big.Int
		envelope string
	)

	err = json.Unmarshal(reply, &envelope)
	if err != nil {
		return nil, err
	}

	err = DecodeHex(envelope, &result)
	if err != nil {
		return nil, err
	}

	return &Wei{result}, nil
}

func (client *Client) EstimateContractGas(
	from string,
	to string,
	data string,
) (*Wei, error) {
	transaction := &TokenTransaction{
		From: from,
		To:   to,
		Data: data,
	}

	reply, err := client.Call("eth_estimateGas", transaction)
	if err != nil {
		return nil, err
	}

	var (
		result   big.Int
		envelope string
	)

	err = json.Unmarshal(reply, &envelope)
	if err != nil {
		return nil, err
	}

	err = DecodeHex(envelope, &result)
	if err != nil {
		return nil, err
	}

	return &Wei{result}, nil
}

// SendTransaction sends specified amount of wei to specified address.
// func (client *Client) SendTransaction(
// 	from string,
// 	to string,
// 	value *Wei,
// 	options ...string,
// ) (*Transaction, error) {
// 	amount, err := EncodeHex(value)
// 	if err != nil {
// 		return nil, err
// 	}

// 	transaction := &Transaction{
// 		From:  from,
// 		To:    to,
// 		Value: amount,
// 	}

// 	response, err := client.Call("eth_sendTransaction", transaction)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(response, &transaction.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"unable to unmarshal transaction ID %q: %s",
// 			string(response),
// 			err,
// 		)
// 	}

// 	return transaction, nil
// }

func (client *Client) PersonalSendTransaction(
	from string,
	to string,
	pwd string,
	value *Wei,
	options ...string,
) (*Transaction, error) {
	amount, err := EncodeHex(value)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		From:  from,
		To:    to,
		Value: amount,
	}

	response, err := client.Call("personal_sendTransaction", transaction, pwd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &transaction.ID)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal transaction ID %q: %s",
			string(response),
			err,
		)
	}

	return transaction, nil
}

func (client *Client) PersonalSendTokenTransaction(
	from string,
	to string, //contract address
	pwd string,
	data string,
) (*TokenTransaction, error) {
	transaction := &TokenTransaction{
		From: from,
		To:   to,
		Data: data,
	}

	response, err := client.Call("personal_sendTransaction", transaction, pwd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &transaction.ID)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal transaction ID %q: %s",
			string(response),
			err,
		)
	}

	return transaction, nil
}

// UnlockAccount unlocks specified account (only if feature is enabled in geth).
func (client *Client) UnlockAccount(
	address string,
	password string,
	duration int,
) error {
	_, err := client.Call("personal_unlockAccount", address, password, duration)
	if err != nil {
		return err
	}

	return nil
}

// GetVersion returns version of connected network.
func (client *Client) GetVersion() (string, error) {
	response, err := client.Call("net_version")
	if err != nil {
		return "", err
	}

	var version string

	err = json.Unmarshal(response, &version)
	if err != nil {
		return "", hierr.Errorf(
			err,
			"unable to unmarshal version: %q",
			response,
		)
	}

	return version, nil
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
		if reply.Error.Code == ErrorCodeAuthenticationNeeded {
			return nil, AuthenticationNeededError{reply.Error.Message}
		}

		return nil, fmt.Errorf(
			"error while processing request (code %d): %s",
			reply.Error.Code,
			reply.Error.Message,
		)
	}

	return reply.Result, nil
}
