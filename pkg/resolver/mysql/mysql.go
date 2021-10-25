package mysql

import (
	"bytes"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/google/martian/log"
	"github.com/pkg/errors"
	"io"
	"net"
	"strings"
	"vitess.io/vitess/go/tb"
)

const initClientConnStatus = ServerStatusAutocommit

type MysqlResolver struct {
	// This is the main listener socket.
	listener net.Listener
	// Incrementing ID for connection id.
	connectionID uint32
	// connReadBufferSize is size of buffer for reads from underlying connection.
	// Reads are unbuffered if it's <=0.
	connReadBufferSize int

	salt []byte

	// ServerVersion is the version we will advertise.
	ServerVersion string

	// Capabilities is the current set of features this connection
	// is using.  It is the features that are both supported by
	// the client and the server, and currently in use.
	// It is set during the initial handshake.
	//
	// It is only used for CapabilityClientDeprecateEOF
	// and CapabilityClientFoundRows.
	Capabilities uint32

	// CharacterSet is the character set used by the other side of the
	// connection.
	// It is set during the initial handshake.
	// See the values in constants.go.
	CharacterSet uint8

	// schemaName is the default database name to use. It is set
	// during handshake, and by ComInitDb packets. Both client and
	// servers maintain it. This member is private because it's
	// non-authoritative: the client can change the schema name
	// through the 'USE' statement, which will bypass this variable.
	schemaName string

	users map[string]string
}

func NewMysqlResolver(listener net.Listener, conf *model.MysqlConfig) *MysqlResolver {
	return &MysqlResolver{
		listener:           listener,
		salt:               []byte(conf.Salt),
		ServerVersion:      conf.ServerVersion,
		users:              conf.Users,
	}
}

func (r *MysqlResolver) Accept() {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			return
		}

		connectionID := r.connectionID
		r.connectionID++

		go r.handle(conn, connectionID)
	}
}

func (r *MysqlResolver) handle(conn net.Conn, connectionID uint32) {
	c := newConn(conn)
	c.ConnectionID = connectionID

	// Catch panics, and close the connection in any case.
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("mysql_server caught panic:\n%v\n%s", x, tb.Stack(4))
		}

		conn.Close()
	}()

	err := r.handshake(c)
	if err != nil {
		c.writeErrorPacketFromError(err)
		return
	}

	// Negotiation worked, send OK packet.
	if err := c.writeOKPacket(0, 0, c.StatusFlags, 0); err != nil {
		log.Errorf("Cannot write OK packet to %s: %v", c, err)
		return
	}

	for {
		// todo parse sql, execute sql
	}
}

func (r *MysqlResolver) handshake(conn *Conn) error {
	// First build and send the server handshake packet.
	err := r.writeHandshakeV10(conn, false, r.salt)
	if err != nil {
		if err != io.EOF {
			log.Errorf("Cannot send HandshakeV10 packet to %s: %v", conn, err)
		}
		return err
	}

	// Wait for the client response. This has to be a direct read,
	// so we don't buffer the TLS negotiation packets.
	response, err := conn.readEphemeralPacketDirect()
	if err != nil {
		// Don't log EOF errors. They cause too much spam, same as main read loop.
		if err != io.EOF {
			log.Infof("Cannot read client handshake response from %s: %v, it may not be a valid MySQL client", conn, err)
		}
		return err
	}

	user, _, authResponse, err := r.parseClientHandshakePacket(conn, true, response)
	if err != nil {
		log.Errorf("Cannot parse client handshake response from %s: %v", conn, err)
		return err
	}
	conn.recycleReadPacket()

	err = r.ValidateHash(user, authResponse)
	if err != nil {
		log.Errorf("Error authenticating user using MySQL native password: %v", err)
		return err
	}
	return nil
}

// writeHandshakeV10 writes the Initial Handshake Packet, server side.
// It returns the salt data.
func (r *MysqlResolver) writeHandshakeV10(c *Conn, enableTLS bool, salt []byte) error {
	capabilities := CapabilityClientLongPassword |
		CapabilityClientFoundRows |
		CapabilityClientLongFlag |
		CapabilityClientConnectWithDB |
		CapabilityClientProtocol41 |
		CapabilityClientTransactions |
		CapabilityClientSecureConnection |
		CapabilityClientMultiStatements |
		CapabilityClientMultiResults |
		CapabilityClientPluginAuth |
		CapabilityClientPluginAuthLenencClientData |
		CapabilityClientDeprecateEOF |
		CapabilityClientConnAttr
	if enableTLS {
		capabilities |= CapabilityClientSSL
	}

	length :=
		1 + // protocol version
			lenNullString(r.ServerVersion) +
			4 + // connection ID
			8 + // first part of salt data
			1 + // filler byte
			2 + // capability flags (lower 2 bytes)
			1 + // character set
			2 + // status flag
			2 + // capability flags (upper 2 bytes)
			1 + // length of auth plugin data
			10 + // reserved (0)
			13 + // auth-plugin-data
			lenNullString(MysqlNativePassword) // auth-plugin-name

	data := c.startEphemeralPacket(length)
	pos := 0

	// Protocol version.
	pos = writeByte(data, pos, protocolVersion)

	// Copy server version.
	pos = writeNullString(data, pos, r.ServerVersion)

	// Add connectionID in.
	pos = writeUint32(data, pos, c.ConnectionID)

	pos += copy(data[pos:], salt[:8])

	// One filler byte, always 0.
	pos = writeByte(data, pos, 0)

	// Lower part of the capability flags.
	pos = writeUint16(data, pos, uint16(capabilities))

	// Character set.
	pos = writeByte(data, pos, CharacterSetUtf8)

	// Status flag.
	pos = writeUint16(data, pos, initClientConnStatus)

	// Upper part of the capability flags.
	pos = writeUint16(data, pos, uint16(capabilities>>16))

	// Length of auth plugin data.
	// Always 21 (8 + 13).
	pos = writeByte(data, pos, 21)

	// Reserved 10 bytes: all 0
	pos = writeZeroes(data, pos, 10)

	// Second part of auth plugin data.
	pos += copy(data[pos:], salt[8:])
	data[pos] = 0
	pos++

	// Copy authPluginName. We always start with mysql_native_password.
	pos = writeNullString(data, pos, MysqlNativePassword)

	// Sanity check.
	if pos != len(data) {
		return errors.Errorf("error building Handshake packet: got %v bytes expected %v", pos, len(data))
	}

	if err := c.writeEphemeralPacket(); err != nil {
		if strings.HasSuffix(err.Error(), "write: connection reset by peer") {
			return io.EOF
		}
		if strings.HasSuffix(err.Error(), "write: broken pipe") {
			return io.EOF
		}
		return err
	}

	return nil
}

// parseClientHandshakePacket parses the handshake sent by the client.
// Returns the username, auth method, auth data, error.
// The original data is not pointed at, and can be freed.
func (r *MysqlResolver) parseClientHandshakePacket(c *Conn, firstTime bool, data []byte) (string, string, []byte, error) {
	pos := 0

	// Client flags, 4 bytes.
	clientFlags, pos, ok := readUint32(data, pos)
	if !ok {
		return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read client flags")
	}
	if clientFlags&CapabilityClientProtocol41 == 0 {
		return "", "", nil, errors.Errorf("parseClientHandshakePacket: only support protocol 4.1")
	}

	// Remember a subset of the capabilities, so we can use them
	// later in the protocol. If we re-received the handshake packet
	// after SSL negotiation, do not overwrite capabilities.
	if firstTime {
		r.Capabilities = clientFlags & (CapabilityClientDeprecateEOF | CapabilityClientFoundRows)
	}

	// set connection capability for executing multi statements
	if clientFlags&CapabilityClientMultiStatements > 0 {
		r.Capabilities |= CapabilityClientMultiStatements
	}

	// Max packet size. Don't do anything with this now.
	// See doc.go for more information.
	_, pos, ok = readUint32(data, pos)
	if !ok {
		return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read maxPacketSize")
	}

	// Character set. Need to handle it.
	characterSet, pos, ok := readByte(data, pos)
	if !ok {
		return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read characterSet")
	}
	r.CharacterSet = characterSet

	// 23x reserved zero bytes.
	pos += 23

	//// Check for SSL.
	//if firstTime && r.TLSConfig != nil && clientFlags&CapabilityClientSSL > 0 {
	//	// Need to switch to TLS, and then re-read the packet.
	//	conn := tls.Server(c.conn, r.TLSConfig)
	//	c.conn = conn
	//	c.bufferedReader.Reset(conn)
	//	r.Capabilities |= CapabilityClientSSL
	//	return "", "", nil, nil
	//}

	// username
	username, pos, ok := readNullString(data, pos)
	if !ok {
		return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read username")
	}

	// auth-response can have three forms.
	var authResponse []byte
	if clientFlags&CapabilityClientPluginAuthLenencClientData != 0 {
		var l uint64
		l, pos, ok = readLenEncInt(data, pos)
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read auth-response variable length")
		}
		authResponse, pos, ok = readBytesCopy(data, pos, int(l))
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read auth-response")
		}

	} else if clientFlags&CapabilityClientSecureConnection != 0 {
		var l byte
		l, pos, ok = readByte(data, pos)
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read auth-response length")
		}

		authResponse, pos, ok = readBytesCopy(data, pos, int(l))
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read auth-response")
		}
	} else {
		a := ""
		a, pos, ok = readNullString(data, pos)
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read auth-response")
		}
		authResponse = []byte(a)
	}

	// db name.
	if clientFlags&CapabilityClientConnectWithDB != 0 {
		dbname := ""
		dbname, pos, ok = readNullString(data, pos)
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read dbname")
		}
		r.schemaName = dbname
	}

	// authMethod (with default)
	authMethod := MysqlNativePassword
	if clientFlags&CapabilityClientPluginAuth != 0 {
		authMethod, pos, ok = readNullString(data, pos)
		if !ok {
			return "", "", nil, errors.Errorf("parseClientHandshakePacket: can't read authMethod")
		}
	}

	// The JDBC driver sometimes sends an empty string as the auth method when it wants to use mysql_native_password
	if authMethod == "" {
		authMethod = MysqlNativePassword
	}

	// Decode connection attributes send by the client
	if clientFlags&CapabilityClientConnAttr != 0 {
		if _, _, err := parseConnAttrs(data, pos); err != nil {
			logger.Warnf("Decode connection attributes send by the client: %v", err)
		}
	}

	return username, authMethod, authResponse, nil
}

func (r *MysqlResolver) ValidateHash(user string, authResponse []byte) error {
	password, ok := r.users[user]
	if !ok {
		return NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "Access denied for user '%v'", user)
	}
	computedAuthResponse := ScramblePassword(r.salt, []byte(password))
	if bytes.Equal(authResponse, computedAuthResponse) {
		return nil
	}
	return NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "Access denied for user '%v'", user)
}

func parseConnAttrs(data []byte, pos int) (map[string]string, int, error) {
	var attrLen uint64

	attrLen, pos, ok := readLenEncInt(data, pos)
	if !ok {
		return nil, 0, errors.Errorf("parseClientHandshakePacket: can't read connection attributes variable length")
	}

	var attrLenRead uint64

	attrs := make(map[string]string)

	for attrLenRead < attrLen {
		var keyLen byte
		keyLen, pos, ok = readByte(data, pos)
		if !ok {
			return nil, 0, errors.Errorf("parseClientHandshakePacket: can't read connection attribute key length")
		}
		attrLenRead += uint64(keyLen) + 1

		var connAttrKey []byte
		connAttrKey, pos, ok = readBytesCopy(data, pos, int(keyLen))
		if !ok {
			return nil, 0, errors.Errorf("parseClientHandshakePacket: can't read connection attribute key")
		}

		var valLen byte
		valLen, pos, ok = readByte(data, pos)
		if !ok {
			return nil, 0, errors.Errorf("parseClientHandshakePacket: can't read connection attribute value length")
		}
		attrLenRead += uint64(valLen) + 1

		var connAttrVal []byte
		connAttrVal, pos, ok = readBytesCopy(data, pos, int(valLen))
		if !ok {
			return nil, 0, errors.Errorf("parseClientHandshakePacket: can't read connection attribute value")
		}

		attrs[string(connAttrKey[:])] = string(connAttrVal[:])
	}

	return attrs, pos, nil

}
