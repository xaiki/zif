// a few network helpers

package main

import (
	"bytes"
	"net"
)

func net_recvall(buf []byte, conn net.Conn) error {
	read := 0

	for read < len(buf) {
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
