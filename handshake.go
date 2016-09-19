package main

import (
	"errors"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
)

func handshake(conn net.Conn, lp *LocalPeer) (ProtocolHeader, error) {
	header, err := handshake_recieve(conn)

	if err != nil {
		return header, err
	}

	if lp == nil {
		return header, errors.New("Handshake passed nil LocalPeer")
	}

	err = handshake_send(conn, lp)

	if err != nil {
		return header, err
	}

	return header, nil
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
	cookie, err := CryptoRandBytes(20)
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
	header.zifAddress.Generate(header.PublicKey[:])
	conn.Write(header.Bytes())

	if !check_ok(conn) {
		return errors.New("Peer refused header")
	}

	// The server will want us to sign this. Proof of identity and all that.
	cookie := make([]byte, 20)
	net_recvall(cookie, conn)

	log.Debug("Cookie recieved, signing")

	sig := lp.Sign(cookie)
	conn.Write(sig)

	if !check_ok(conn) {
		return errors.New("Peer refused signature")
	}

	return nil
}
