package frcntgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/techplexengineer/frc-networktables-go/entryupdate"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/techplexengineer/frc-networktables-go/entry"
	"github.com/techplexengineer/frc-networktables-go/message"
	"github.com/techplexengineer/frc-networktables-go/util"
)

// ClientStatus is the enum type to represent the different
// states/statuses the client could have
type ClientStatus int

const (
	// ClientDisconnected indicates that the client cannot reach
	// the server
	ClientDisconnected ClientStatus = iota
	// ClientConnected indicates that the client has connected to
	// the server but has not began actual communication
	ClientConnected
	// ClientSentHello indicates that the client has sent the hello
	// packets and is waiting for a response from the server
	ClientSentHello
	// ClientStartingSync indicates that the client has received the
	// server hello and is beginning to synchronize values with the
	// server.
	ClientStartingSync
	// ClientInSync indicates that the client is completely in sync
	// with the server and has all the correct values.
	ClientInSync
	// keepAliveTime is the amount of time (seconds) between packets
	// that the client waits before it sends a KeepAlive message.
	// It is advised to never have this lower than one second so as to
	// prevent overloading the server.
	keepAliveTime int64 = 1
)

var (
	lastSent            int64
	messageOutgoingChan = make(chan message.IMessage)
)

// Client is the NetworkTables Client
type Client struct {
	//handler ClientMessageHandler
	conn    net.Conn
	entries map[string]entry.IEntry
	status  ClientStatus
}

func (c *Client) GetStatus() ClientStatus {
	return c.status
}

// Create a new network tables client connecting to a server on localhost
func NewClientLocalhost() (*Client, error) {
	address := "0.0.0.0"
	port := "1735"
	return NewClient(address, port)
}

// Create a new network tables client using mdns given a team number
func NewClientTeam(teamNumber int) (*Client, error) {
	address := fmt.Sprintf("roboRIO-%d-FRC.local", teamNumber) //no leading zeros
	port := "1735"
	return NewClient(address, port)
}

// Create a new Network Tables client
// connAddr can be an IP or hostname
// connPort is the tcp port to connect to. Usually Network Tables uses port 1735
func NewClient(connAddr, connPort string) (*Client, error) {
	tcpConn, err := net.Dial("tcp", util.ConcatAddress(connAddr, connPort))
	if err != nil {
		return &Client{
			conn:    nil,
			entries: map[string]entry.IEntry{},
			status:  ClientDisconnected,
		}, err
	}
	client := Client{
		conn:    tcpConn,
		entries: map[string]entry.IEntry{},
		status:  ClientConnected,
	}
	defer client.connect()
	return &client, nil
}

func (c *Client) connect() {
	go c.processOutgoingQueue()
	go c.receiveIncoming()
	c.startHandshake()
}

func (c *Client) startHandshake() {
	clientName := []byte("frc-nt-golang")
	clientLength := util.EncodeULeb128(uint32(len(clientName)))
	clientName = append(clientLength, clientName...)
	VERSION := [2]byte{0x03, 0x00}

	// Step 1: Client sends Client Hello
	helloMessage := message.ClientHelloFromItems(VERSION, clientName)
	c.QueueMessage(helloMessage)
	c.status = ClientSentHello
}

// Close disconnects and closes the client from the server.
func (c *Client) Close() error {
	if c.status == ClientDisconnected {
		return errors.New("client: Already disconnected")
	}
	c.status = ClientDisconnected
	c.conn.Close()
	return nil
}

// QueueMessage prepares the message that has been provided for
// sending.
func (c *Client) QueueMessage(message message.IMessage) error {
	if c.status != ClientDisconnected {
		fmt.Printf("<=== Sending msg type %#x - %s\n", message.GetType().Byte(), message.GetType().String())
		messageOutgoingChan <- message
		return nil
	}
	return errors.New("client: server could not be reached")
}

// process the queue of outgoing messages, should be called as a gofun
func (c *Client) processOutgoingQueue() {
	for c.status != ClientDisconnected {
		sending := <-messageOutgoingChan
		c.conn.Write(sending.CompressToBytes())
		defer updateLastSent() //@todo why is this deferred?
	}
	log.Printf("client: server could not be reached")
}

// readMessage
func (c *Client) receiveIncoming() {
	var potentialMessage [1]byte
	for c.status != ClientDisconnected {
		_, ioError := io.ReadFull(c.conn, potentialMessage[:])
		if ioError != nil {
			if ioError == io.EOF {
				continue
			}
			c.Close()
			log.Printf("io error: %s", ioError)
			return //don't attempt to process any further
		}
		tempPacket, messageError := message.BuildFromReader(message.MessageType(potentialMessage[0]), c.conn)
		if messageError != nil {
			c.Close()
			log.Printf("Message Error: %s", messageError)
			return //don't attempt to process any further
		}
		fmt.Printf("===> got %s\n", tempPacket.GetType())
		switch tempPacket.GetType() {

		case message.TypeServerHello:
			// Step 2: Server replies to ClientHello with ServerHello
			msg := tempPacket.(*message.ServerHello)
			fmt.Printf("Connected to %s\n", msg.GetServerIdentity())

		case message.TypeEntryAssign:
			// Step 3: Server sends EntryAssign messages for each entry
			// if we just sent a client hello then the sync is beginning
			if c.status == ClientSentHello {
				c.status = ClientStartingSync
			}
			msg := tempPacket.(*message.EntryAssign)
			c.entries[msg.GetEntry().GetName()] = msg.GetEntry()
		case message.TypeServerHelloComplete:
			// Step 4: The Server sends a Server Hello Complete message.
			// Server is done sending entryAssigns

			// Step 5: For all Entries the Client recognizes that the Server
			// did not identify with a Entry Assignment.
			// we can now send any entries the server should have @todo not implemented

			// Step 6: The Client sends a Client Hello Complete message.
			msg := message.ClientHelloCompleteFromItems()
			c.QueueMessage(msg)
			c.status = ClientInSync
		case message.TypeEntryUpdate:
			msg := tempPacket.(*message.EntryUpdate)
			up := msg.GetUpdate()
			for _, e := range c.entries {
				if up.GetID() == e.GetID() {
					if up.GetType() != e.GetType() {
						log.Printf("Types differ. Ignoring update")
						break
					}
					switch up.GetType() {
					case entry.TypeBoolean:
						eu := up.(*entryupdate.Boolean)
						e.SetValue(eu.GetValue())

					case entry.TypeDouble:
					case entry.TypeString:
					case entry.TypeRaw:
					case entry.TypeBooleanArr:
					case entry.TypeDoubleArr:
					case entry.TypeStringArr:
					case entry.TypeRPCDef:
					default:
						log.Printf("Unknown type! Unable to update entry")
					}
					fmt.Printf("%v", e.GetID())
				}
			}
		case message.TypeClientHelloComplete:
			// only expect to get this message on the server
		case message.TypeKeepAlive:
			// can be safely ignored
		case message.TypeClientHello:
			// only expected on the server
		case message.TypeProtoUnsupported:
			// @todo
		case message.TypeEntryFlagUpdate:
			// @todo
		case message.TypeEntryDelete:
			// @todo
		case message.TypeClearAllEntries:
			// @todo
		case message.TypeRPCExec:
			// @todo
		case message.TypeRPCResponse:
			// @todo
		default:
			fmt.Printf("===> got UNKNOWN\n")
		}
	}
}

func updateLastSent() {
	currentTime := time.Now()
	lastSent = currentTime.Unix()
}

// keepAlive should be run in a Go routine. It sends a
// the provided packet after the provided time (seconds) have
// passed between the last packet.
func (c *Client) keepAlive(packet message.IMessage) {
	for c.status == ClientInSync {
		currentTime := time.Now()
		currentSeconds := currentTime.Unix()
		if (currentSeconds - lastSent) >= keepAliveTime {
			go c.QueueMessage(packet)
		}
	}
}

// GetBoolean fetches a boolean at the specified key
func (c *Client) GetBoolean(key string) (bool, error) {
	key = util.SanitizeKey(key)
	entry, ok := c.entries[key]
	if !ok {
		return false, fmt.Errorf("key is missing")
	}

	return entry.GetValue().(bool), nil
}

//Set function to be called when robot connects/disconnects
//func (c Client) AddRobotConnectionListener(callback func()) {}
//func (c Client) AddKeyListener(key string, callback func()) {}

func (c Client) GetKeys(prefix string) []string {
	keys := []string{}
	for k, _ := range c.entries {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys
}

// Determines whether the given key is in this table.
func (c Client) ContainsKey(key string) bool {
	_, ok := c.entries[key]
	return ok
}

func (c Client) GetEntry(key string) interface{} {
	e := c.entries[key]
	return e.GetValue()
}

type SnapShotEntry struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Datatype string `json:"type"`
}

func (c Client) GetSnapshot(prefix string) []SnapShotEntry {
	keys := []SnapShotEntry{}
	for k, v := range c.entries {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			valueStr := fmt.Sprintf("%#v", v.GetValue())
			valueByt, err := json.Marshal(v.GetValue())
			if err == nil {
				valueStr = string(valueByt)
			}

			keys = append(keys, SnapShotEntry{
				Key:      k,
				Value:    valueStr,
				Datatype: v.GetType().String(),
			})
		}
	}
	return keys
}
