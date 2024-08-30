package cdkey

// A map to store the handlers for each message type
var CDKeyHandlers = map[uint8]func(*CDKeyMessage, *Client) (*CDKeyMessage, error){}
