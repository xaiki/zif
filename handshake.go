package zif

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

/*func handshake(conn net.Conn, lp *LocalPeer) (ProtocolHeader, error) {
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

	log.WithFields(log.Fields{"peer": pHeader.zifAddress.Encode()}).Info("Incoming connection")

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

	log.WithFields(log.Fields{"peer": pHeader.zifAddress.Encode()}).Info("Verified")

	return pHeader, nil
}
*/
func handshake_send(cl Client, lp *LocalPeer) error {
	log.Debug("Handshaking with ", cl.conn.RemoteAddr().String())

	header := Message{
		Header:  ProtoHeader,
		Content: lp.PublicKey,
	}

	cl.WriteMessage(header)

	msg, err := cl.ReadMessage()

	if err != nil {
		return err
	}

	if !msg.Ok() {
		return errors.New("Peer refused header")
	}

	msg, err = cl.ReadMessage()

	if err != nil {
		return err
	}

	if !msg.Ok() {
		return errors.New("Peer refused header")
	}

	log.Debug("Cookie recieved, signing")

	sig := lp.Sign(msg.Content)

	msg = Message{
		Header:  ProtoSig,
		Content: sig,
	}

	cl.WriteMessage(msg)

	msg, err = cl.ReadMessage()

	if err != nil {
		return err
	}

	if !msg.Ok() {
		return errors.New("Peer refused signature")
	}

	finish

	return nil
}
