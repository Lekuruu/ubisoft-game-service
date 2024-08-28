package router

// A map to store the handlers for each message type
var RouterHandlers = map[uint8]func(*GSMessage, *Client) (*GSMessage, error){}
