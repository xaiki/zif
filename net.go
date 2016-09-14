// a few network helpers

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

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

// len returns an int
// so why uint64?
// One day this protocol may be implemented in not-Go, and I'd just rather not
// be constrained to signed integers that could be either 32 or 64 bit :)
// this way it is known what is going on
func net_sendlength(conn net.Conn, length uint64) error {
	length_b := make([]byte, 8)
	binary.PutUvarint(length_b, length)

	_, err := conn.Write(length_b)

	return err
}

func net_recvlength(conn net.Conn) (uint64, error) {
	length_b := make([]byte, 8)
	err := net_recvall(length_b, conn)

	if err != nil {
		return 0, err
	}

	length, _ := binary.Uvarint(length_b)

	return length, nil
}

func check_ok(conn net.Conn) bool {
	buf := make([]byte, 2)

	net_recvall(buf, conn)

	return bytes.Equal(buf, proto_ok)
}

func recieve_entry(conn net.Conn) (Entry, error) {
	length_b := make([]byte, 8)
	net_recvall(length_b, conn)
	length, _ := binary.Varint(length_b)

	if length > EntryLengthMax {
		return Entry{}, errors.New("Peer entry larger than max")
	}

	entry_json := make([]byte, length)
	net_recvall(entry_json, conn)

	entry, err := JsonToEntry(entry_json)

	err = ValidateEntry(&entry)

	return entry, err
}

func listen_stream(peer *Peer) {
	var err error
	session := peer.GetSession()

	if session == nil {
		log.Info("Peer has no active session, starting server")
		session, err = peer.ConnectServer()

		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for {
		stream, err := session.Accept()

		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Peer closed connection")
			} else {
				log.Error(err.Error())
			}

			peer.localPeer.CheckSessions()

			return
		}

		log.Debug("Accepted stream (", session.NumStreams(), " total)")

		peer.AddStream(stream)

		go handle_stream(peer, stream)
	}
}

func handle_stream(peer *Peer, stream net.Conn) {
	log.Debug("Handling stream")
	msg := make([]byte, 2)
	for {
		err := net_recvall(msg, stream)

		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Closed stream from ", peer.ZifAddress.Encode())
			} else {
				log.Error(err.Error())
			}

			peer.RemoveStream(stream)

			return
		}

		if bytes.Equal(msg, proto_terminate) {
			peer.Terminate()
			log.Debug("Terminated connection with ", peer.ZifAddress.Encode())
			return
		}

		RouteMessage(msg, peer, stream)
	}
}
