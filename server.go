package zif

// tcp server

import (
	"bytes"
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

	var err error
	session := peer.GetSession()

	if session == nil {
		log.WithField("peer", peer.ZifAddress.Encode()).Info("No session, starting server")
		session, err = peer.ConnectServer()

		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for {
		stream, err := session.Accept()
		limiter.Wait()

		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Peer closed connection")
			} else {
				log.Error(err.Error())
			}

			s.localPeer.CheckSessions()

			return
		}

		log.Debug("Accepted stream (", session.NumStreams(), " total)")

		peer.AddStream(stream)

		go s.HandleStream(peer, stream)
	}
}

func (s *Server) HandleStream(peer *Peer, stream net.Conn) {
	log.Debug("Handling stream")
	msg := make([]byte, 2)
	for {
		err := net_recvall(msg, stream)

		if err != nil {
			if err.Error() == "EOF" {
				log.WithField("peer", peer.ZifAddress.Encode()).Info("Closed stream")
			} else {
				log.Error(err.Error())
			}

			peer.RemoveStream(stream)

			return
		}

		select {
		case s.localPeer.MsgChan <- msg:

		default:

		}
		if bytes.Equal(msg, proto_terminate) {
			peer.Terminate()
			log.Debug("Terminated connection with ", peer.ZifAddress.Encode())
			return
		}

		s.RouteMessage(msg, peer, stream)
	}
}

func (s *Server) RouteMessage(msg_type []byte, from *Peer, stream net.Conn) {
	//log.Debug("Routing message ", msg_type)

	if bytes.Equal(msg_type, proto_ping) {
		rep := Client{stream}
		rep.Pong()
	} else if bytes.Equal(msg_type, proto_pong) {
		log.Debug("Pong from ", from.ZifAddress.Encode())
	} else if bytes.Equal(msg_type, proto_dht_announce) {
		s.localPeer.HandleAnnounce(stream, from)
	} else if bytes.Equal(msg_type, proto_dht_query) {
		s.localPeer.HandleQuery(stream, from)
	} else if bytes.Equal(msg_type, proto_search) {
		s.localPeer.HandleSearch(stream, from)
	}
}

func (s *Server) Handshake(conn net.Conn) {
	header, err := handshake(conn, s.localPeer)

	if err != nil {
		log.Error(err.Error())
		return
	}

	var peer Peer
	peer.SetTCP(ConnHeader{conn, header})

	s.localPeer.AddPeer(&peer)

	go s.ListenStream(&peer)
}

func (s *Server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}
