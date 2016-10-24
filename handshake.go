package zif

import (
	"errors"

	"golang.org/x/crypto/ed25519"

	log "github.com/sirupsen/logrus"
)

func handshake(cl Client, lp *LocalPeer) (ed25519.PublicKey, error) {
	header, err := handshake_recieve(cl)

	if err != nil {
		return header, err
	}

	if lp == nil {
		return header, errors.New("Handshake passed nil LocalPeer")
	}

	err = handshake_send(cl, lp)

	if err != nil {
		return header, err
	}

	return header, nil
}

func handshake_recieve(cl Client) (ed25519.PublicKey, error) {
	check := func(e error) bool {
		if e != nil {
			log.Error(e.Error())
			cl.Close()
			return true
		}

		return false
	}

	header, err := cl.ReadMessage()

	if check(err) {
		cl.WriteMessage(Message{Header: ProtoNo})
		return nil, err
	}

	cl.WriteMessage(Message{Header: ProtoOk})

	address := Address{}
	address.Generate(header.Content)

	log.WithFields(log.Fields{"peer": address.Encode()}).Info("Incoming connection")

	// Send the client a cookie for them to sign, this proves they have the
	// private key, and it is highly unlikely an attacker has a signed cookie
	// cached.
	cookie, err := CryptoRandBytes(20)
	if check(err) {
		return nil, err
	}

	cl.WriteMessage(Message{Header: ProtoCookie, Content: cookie})

	sig, err := cl.ReadMessage()

	if check(err) {
		return nil, err
	}

	verified := ed25519.Verify(header.Content, cookie, sig.Content)

	if !verified {
		log.Error("Failed to verify peer ", address.Encode())
		cl.WriteMessage(Message{Header: ProtoNo})
		cl.Close()
		return nil, errors.New("Signature not verified")
	}

	cl.WriteMessage(Message{Header: ProtoOk})

	log.WithFields(log.Fields{"peer": address.Encode()}).Info("Verified")

	return header.Content, nil
}

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
		log.Error(err.Error())
		return err
	}

	log.Info("Cookie recieved, signing")

	sig := lp.Sign(msg.Content)

	msg = &Message{
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

	return nil
}
