// a few network helpers

package main

import "net"

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

func check_ok() {

}
