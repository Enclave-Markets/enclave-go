package apiclient

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketConn struct {
	wsConn        *websocket.Conn
	readDeadline  time.Duration
	writeDeadline time.Duration
}

const DefaultTimeout = 5 * time.Second

func GetTlsConfig(apiEndpoint string) *tls.Config {
	insecureSkipVerify := false
	u, err := url.Parse(apiEndpoint)
	if err == nil {

		host, _, _ := net.SplitHostPort(u.Host)
		switch host {
		case "localhost", "127.0.0.1", "0.0.0.0", "engine", "monolith":
			insecureSkipVerify = true
		default:
			// single name is probably from Docker
			parts := strings.Split(host, ".")
			if len(parts) == 1 {
				insecureSkipVerify = true
			}
		}
	}

	return &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
}

func (client *ApiClient) NewWebsocketConnection() (*WebsocketConn, error) {
	u, err := url.Parse(client.ApiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the api endpoint %s: %s", client.ApiEndpoint, err.Error())
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		if u.Host != "" {
			host = u.Host
		} else {
			return nil, err
		}
	}
	spotWsEndpoint := fmt.Sprintf("wss://%s/ws", host)

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = GetTlsConfig(spotWsEndpoint)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	conn, _, err := dialer.DialContext(ctx, spotWsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	wsConn := &WebsocketConn{
		wsConn:        conn,
		readDeadline:  DefaultTimeout,
		writeDeadline: DefaultTimeout,
	}

	err = wsConn.websocketLogin(client)
	if err != nil {
		return nil, err
	}

	return wsConn, nil
}

func (c *ApiClient) GetWebsocketLoginArgs() *RequestArgs {
	if c.apiKeyArgs != nil {
		c.computeApiKeyArgs("enclave_ws_login", "", nil)
		return &RequestArgs{
			KeyId:          c.apiKeyArgs.KeyId,
			TimeUnixMillis: c.apiKeyArgs.Timestamp,
			Sign:           c.apiKeyArgs.Sign,
		}
	} else {
		return nil
	}
}

func (c *WebsocketConn) SetReadDeadline(t time.Duration) {
	c.readDeadline = t
}

func (c *WebsocketConn) SetWriteDeadline(t time.Duration) {
	c.writeDeadline = t
}

func (c *WebsocketConn) SendMessage(req WebSocketAPIRequest) error {
	var err error
	if c.writeDeadline == 0 {
		err = c.wsConn.SetWriteDeadline(time.Time{})
	} else {
		err = c.wsConn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
	}
	if err != nil {
		return err
	}
	err = c.wsConn.WriteJSON(req)
	_ = c.wsConn.SetWriteDeadline(time.Time{})
	return err
}

func (c *WebsocketConn) ReadMessage() (*WebSocketAPIResponse, error) {
	var err error
	if c.readDeadline == 0 {
		err = c.wsConn.SetReadDeadline(time.Time{})
	} else {
		err = c.wsConn.SetReadDeadline(time.Now().Add(c.readDeadline))
	}
	if err != nil {
		return nil, err
	}

	_, p, err := c.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	err = c.wsConn.SetReadDeadline(time.Time{})
	if err != nil {
		return nil, err
	}

	res, err := UnmarshalWebSocketAPIResponse(p)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *WebsocketConn) Close() error {
	return c.wsConn.Close()
}

func IsCloseError(err error) bool {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return true
	}
	netErr, ok := err.(*net.OpError)
	if ok && netErr.Unwrap() == net.ErrClosed {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	return false
}

func (conn *WebsocketConn) websocketLogin(client *ApiClient) error {
	// Skip login if no credentials are provided
	args := client.GetWebsocketLoginArgs()
	if args == nil {
		return nil
	}
	req := WebSocketAPIRequest{Op: Login, Args: args}
	err := conn.SendMessage(req)
	if err != nil {
		conn.wsConn.Close()
		return fmt.Errorf("failed to write login message to websocket: %s", err.Error())
	}

	res, err := conn.ReadMessage()
	if err != nil {
		conn.wsConn.Close()
		return fmt.Errorf("failed to read message from websocket: %s", err.Error())
	}

	if res.Type != LoggedIn {
		conn.wsConn.Close()
		return fmt.Errorf("unexpected response to login message: %s %s", res.Type, res.Msg)
	}

	return nil
}
