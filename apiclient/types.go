package apiclient

import (
	"encoding/json"
	"fmt"

	"github.com/Enclave-Markets/enclave-go/models"
)

// Requests

type WebSocketAPIRequest struct {
	Op            WSRequestType   `json:"op,omitempty"`
	Channel       ChannelType     `json:"channel,omitempty"`
	Markets       []models.Market `json:"markets,omitempty"`
	NumPastTrades *int            `json:"numPastTrades,omitempty"` // This is an internally used parameter for Spot Trades
	Timeout       *int            `json:"timeout_seconds,omitempty"`
	Args          *RequestArgs    `json:"args,omitempty"`
}

type RequestArgs struct {
	Token          string `json:"token"`
	KeyId          string `json:"key"`
	SubaccountId   string `json:"subaccount"`
	TimeUnixMillis string `json:"time"`
	Sign           string `json:"sign"`
}

func (a *RequestArgs) IsJWTLogin() bool {
	return a.Token != ""
}

type WSRequestType string

const (
	Ping        WSRequestType = "ping"
	Subscribe   WSRequestType = "subscribe"
	Unsubscribe WSRequestType = "unsubscribe"
	Login       WSRequestType = "login"
)

type ChannelType string

func TopOfBooksPerps() ChannelType {
	return ChannelType("topOfBooksPerps")
}
func FillsPerps() ChannelType {
	return ChannelType("fillsPerps")
}
func PerpsPositions() ChannelType {
	return "positionsPerps"
}
func PerpsMarkPrices() ChannelType {
	return "markPricesPerps"
}
func TopOfBooksSpot() ChannelType {
	return ChannelType("topOfBooksSpot")
}
func FillsSpot() ChannelType {
	return ChannelType("fillsSpot")
}

// Responses

type WebSocketAPIResponse struct {
	Type    WSResponseType `json:"type"`
	Channel ChannelType    `json:"channel,omitempty"`
	Code    int            `json:"code,omitempty"`
	Msg     string         `json:"msg,omitempty"`
	Data    any            `json:"data,omitempty"`
}

func UnmarshalWebSocketAPIResponse(data []byte) (*WebSocketAPIResponse, error) {
	type Temp struct {
		Type    WSResponseType  `json:"type"`
		Channel ChannelType     `json:"channel,omitempty"`
		Code    int             `json:"code,omitempty"`
		Msg     string          `json:"msg,omitempty"`
		Data    json.RawMessage `json:"data,omitempty"`
	}

	var result *WebSocketAPIResponse

	var parsed Temp
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return result, err
	}
	result = &WebSocketAPIResponse{
		Type:    parsed.Type,
		Channel: parsed.Channel,
		Code:    parsed.Code,
		Msg:     parsed.Msg,
	}
	switch parsed.Type {
	case Pong:
		result.Data = nil
	case Error:
		result.Data = nil
	case Subscribed:
		var req WebSocketAPIRequest
		err = json.Unmarshal(parsed.Data, &req)
		result.Data = req
	case Unsubscribed:
		var req WebSocketAPIRequest
		err = json.Unmarshal(parsed.Data, &req)
		result.Data = req
	case Update:
		switch parsed.Channel {
		case TopOfBooksPerps(), TopOfBooksSpot():
			var temp []*models.ApiBookSnapshot
			err = json.Unmarshal(parsed.Data, &temp)
			result.Data = temp
		case FillsPerps(), FillsSpot():
			var temp []*models.ApiFill
			err = json.Unmarshal(parsed.Data, &temp)
			result.Data = temp
		case PerpsPositions():
			var temp []*models.ApiPosition
			err = json.Unmarshal(parsed.Data, &temp)
			result.Data = temp
		case PerpsMarkPrices():
			var temp []*models.GetMarkPriceRes
			err = json.Unmarshal(parsed.Data, &temp)
			result.Data = temp
		default:
			return nil, fmt.Errorf("couldn't unmarshal response: unknown channel %s", parsed.Channel)
		}
	case LoggedIn:
		result.Data = nil
	}

	return result, err
}

type WSResponseType string

const (
	Pong         WSResponseType = "pong"
	Error        WSResponseType = "error"
	Subscribed   WSResponseType = "subscribed"
	Unsubscribed WSResponseType = "unsubscribed"
	Update       WSResponseType = "update"
	LoggedIn     WSResponseType = "loggedIn"
)
