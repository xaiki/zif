// a few network helpers

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"golang.org/x/crypto/ed25519"

	log "github.com/sirupsen/logrus"
)

type ConnHeader struct {
	conn   net.Conn
	header ProtocolHeader
}

func net_recvall(buf []byte, conn net.Conn) error {
	read := 0

	for read < len(buf) {

		if conn == nil {
			return errors.New("Cannot read, connection nil")
		}

		r, err := conn.Read(buf[read:])

		if err != nil {
			return err
		}

		read += r
	}

	return nil
}

func check_ok(conn net.Conn) bool {
	buf := make([]byte, 2)

	net_recvall(buf, conn)

	return bytes.Equal(buf, proto_ok)
}

func handshake_recieve(conn net.Conn) (ProtocolHeader, error) {
	check := func(e error) bool {
		if e != nil {
			log.Error(e.Error())
			conn.Close()
			return true
		}

		return false
	}

	header := make([]byte, ProtocolHeaderSize)
	err := net_recvall(header, conn)
	if check(err) {
		conn.Write(proto_no)
		return ProtocolHeader{}, err
	}

	pHeader, err := ProtocolHeaderFromBytes(header)
	if check(err) {
		conn.Write(proto_no)
		return pHeader, err
	}

	conn.Write(proto_ok)

	log.Info("Incoming connection from ", pHeader.zifAddress.Encode())

	// Send the client a cookie for them to sign, this proves they have the
	// private key, and it is highly unlikely an attacker has a signed cookie
	// cached.
	cookie, err := RandBytes(20)
	if check(err) {
		return pHeader, err
	}

	conn.Write(cookie)

	sig := make([]byte, ed25519.SignatureSize)
	net_recvall(sig, conn)

	verified := ed25519.Verify(pHeader.PublicKey[:], cookie, sig)

	if !verified {
		log.Error("Failed to verify peer ", pHeader.zifAddress.Encode())
		conn.Write(proto_no)
		conn.Close()
		return pHeader, errors.New("Signature not verified")
	}

	conn.Write(proto_ok)

	log.Info(fmt.Sprintf("Verified %s", pHeader.zifAddress.Encode()))

	return pHeader, nil
}

func handshake_send(conn net.Conn, lp *LocalPeer) error {
	log.Debug("Handshaking with ", conn.RemoteAddr().String())
	//ph := c.localPeer.ProtocolHeader()

	header := lp.ProtocolHeader()
	conn.Write(header.Bytes())

	if !check_ok(conn) {
		return errors.New("Peer refused header")
	}

	// The server will want us to sign this. Proof of identity and all that.
	cookie := make([]byte, 20)
	net_recvall(cookie, conn)

	sig := lp.Sign(cookie)
	conn.Write(sig)

	if !check_ok(conn) {
		return errors.New("Peer refused signature")
	}

	return nil
}

func recieve_entry(conn net.Conn) (Entry, []byte, error) {
	length_b := make([]byte, 8)
	net_recvall(length_b, conn)
	length, _ := binary.Varint(length_b)

	if length > EntryLengthMax {
		return Entry{}, nil, errors.New("Peer entry larger than max")
	}

	entry_json := make([]byte, length)
	net_recvall(entry_json, conn)

	entry, err := JsonToEntry(entry_json)

	sig := make([]byte, ed25519.SignatureSize)
	net_recvall(sig, conn)

	if !ValidateEntry(&entry, sig) {
		return entry, sig, errors.New("Failed to validate entry")
	}

	return entry, sig, err
}

// TODO: Make this check using UpNp/NAT_PMP first, then query services.
func external_ip() string {
	resp, err := http.Get("https://api.ipify.org/")
	defer resp.Body.Close()

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	ret, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("Failed to get external ip: try setting manually")
		return ""
	}

	return string(ret)
}
