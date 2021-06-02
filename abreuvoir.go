package abreuvoir

var (
	address, port string = "0.0.0.0", "1735"
	client        *Client

	// VERSION is the current NetworkTables protocol revision
	VERSION = [2]byte{0x03, 0x00}
	// IDENTITY is the identity of the client or server
	IDENTITY = "ABREUVOIR"
)

// SetAddress sets the address of the remote server
// It's advised to use a local host name instead of
// an IP for better
func SetAddress(newAddress string) {
	address = newAddress
}

// SetPort sets the port of the remote server
func SetPort(newPort string) {
	port = newPort
}

// InitClient initializes the client and the connection to the remote server.
func InitClient() (*Client, error) {
	var tempClient, err = NewClient(address, port)
	if err != nil {
		return nil, err
	}
	client = tempClient
	return tempClient, nil
}
