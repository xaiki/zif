package libzif

// tcp server

import (
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	listener  net.Listener
	localPeer *LocalPeer
}

func (s *Server) Listen(addr string) {
	var err error

	s.listener, err = net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	log.Info("Listening on ", addr)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Error(err.Error())
		}

		log.Debug("Handshaking new connection")
		go s.Handshake(conn)
	}
}

func (s *Server) ListenStream(peer *Peer) {
	// Allowed to open 4 streams per second, bursting to three.
	limiter := NewLimiter(time.Second/4, 3, true)
	defer limiter.Stop()

	session := peer.GetSession()

	for {
		stream, err := session.Accept()
		limiter.Wait()

		if err != nil {
			if err == io.EOF {
				log.Info("Peer closed connection")
				s.localPeer.Peers.Remove(peer.ZifAddress.Encode())

				if peer.entry == nil {
					return
				}

				s.localPeer.Peers.Remove(fmt.Sprintf("%s:%d",
					peer.entry.PublicAddress, peer.entry.Port))
			} else {
				log.Error(err.Error())
			}

			return
		}

		log.Debug("Accepted stream (", session.NumStreams(), " total)")

		peer.AddStream(stream)

		go s.HandleStream(peer, stream)
	}
}

func (s *Server) HandleStream(peer *Peer, stream net.Conn) {
	log.Debug("Handling stream")

	cl := Client{stream, nil, nil}

	for {
		msg, err := cl.ReadMessage()

		if err != nil {
			log.Error(err.Error())
			return
		}
		msg.From = peer

		select {
		case s.localPeer.MsgChan <- *msg:
		default:
		}

		s.RouteMessage(msg)
	}
}

func (s *Server) RouteMessage(msg *Message) {
	var err error

	switch msg.Header {

	case ProtoDhtAnnounce:
		err = s.localPeer.HandleAnnounce(msg)
	case ProtoDhtQuery:
		err = s.localPeer.HandleQuery(msg)
	case ProtoSearch:
		err = s.localPeer.HandleSearch(msg)
	case ProtoRecent:
		err = s.localPeer.HandleRecent(msg)
	case ProtoPopular:
		err = s.localPeer.HandlePopular(msg)
	case ProtoRequestHashList:
		err = s.localPeer.HandleHashList(msg)
	case ProtoRequestPiece:
		err = s.localPeer.HandlePiece(msg)

	default:
		log.Error("Unknown message type")

	}

	if err != nil {
		log.Error(err.Error())
	}

}

func (s *Server) Handshake(conn net.Conn) {
	cl := Client{conn, nil, nil}

	header, err := handshake(cl, s.localPeer)
	addr := Address{}

	if err != nil {
		log.Error(err.Error())
		return
	}

	_, err = addr.Generate(header)

	if err != nil {
		log.Error(err.Error())
		return
	}

	peer := &Peer{}
	peer.SetTCP(ConnHeader{cl, header})
	_, err = peer.ConnectServer()

	if err != nil {
		log.Error(err.Error())
		return
	}

	s.localPeer.Peers.Set(peer.ZifAddress.Encode(), peer)

	go s.ListenStream(peer)
}

func (s *Server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}
