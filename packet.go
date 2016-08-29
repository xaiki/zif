package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"strconv"
	"strings"
)

const PacketSize = 2 + ed25519.PublicKeySize + 8 + 8 + 20 + 8

// Packet binary format
// ============
// protocol message
// public key
// udp server port
// tcp server port
// cookie
// size of the following
// any other content, may be nothing
// ============

type Packet struct {
	protoMsg  [2]byte
	publicKey [ed25519.PublicKeySize]byte
	udpPort   [8]byte
	tcpPort   [8]byte
	cookie    [20]byte
	dataSize  [8]byte
	data      []byte
}

func (p *Packet) PublicKey() ed25519.PublicKey {
	return ed25519.PublicKey(p.publicKey[:])
}

func (p *Packet) UDPPort() int64 {
	port, _ := binary.Varint(p.udpPort[:])

	return port
}

func (p *Packet) TCPPort() int64 {
	port, _ := binary.Varint(p.tcpPort[:])

	return port
}

func (p *Packet) DataSize() int64 {
	size, _ := binary.Varint(p.dataSize[:])

	return size
}

func (p *Packet) SetData(data []byte) {
	p.data = data
	binary.PutVarint(p.dataSize[:], int64(len(p.data)))
}

func (p *Packet) Bytes() []byte {
	packet := make([]byte, 0, PacketSize)

	packet = append(packet, p.protoMsg[:]...)
	packet = append(packet, p.publicKey[:]...)
	packet = append(packet, p.udpPort[:]...)
	packet = append(packet, p.tcpPort[:]...)
	packet = append(packet, p.cookie[:]...)
	packet = append(packet, p.dataSize[:]...)
	packet = append(packet, p.data...)

	return packet
}

func PacketFromBytes(bytes []byte) Packet {
	var packet Packet

	copy(packet.protoMsg[:], bytes[:2])
	copy(packet.publicKey[:], bytes[2:ed25519.PublicKeySize])
	copy(packet.udpPort[:], bytes[2+ed25519.PublicKeySize:ed25519.PublicKeySize+10])
	copy(packet.tcpPort[:], bytes[10+ed25519.PublicKeySize:ed25519.PublicKeySize+18])
	copy(packet.cookie[:], bytes[18+ed25519.PublicKeySize:ed25519.PublicKeySize+38])
	copy(packet.dataSize[:], bytes[38+ed25519.PublicKeySize:46+ed25519.PublicKeySize])

	packet.data = make([]byte, packet.DataSize())
	copy(packet.data[:], bytes[PacketSize:PacketSize+packet.DataSize()])

	return packet
}

func MakePacket(proto_msg, cookie []byte, from *LocalPeer) Packet {
	var packet Packet

	copy(packet.protoMsg[:], proto_msg)
	copy(packet.publicKey[:], from.publicKey)

	udp_port, err := strconv.Atoi(strings.Split(from.DHTAddress, ":")[1])
	tcp_port, err := strconv.Atoi(strings.Split(from.RouterAddress, ":")[1])

	if err != nil {
		panic(err)
	}

	binary.PutVarint(packet.udpPort[:], int64(udp_port))
	binary.PutVarint(packet.tcpPort[:], int64(tcp_port))

	copy(packet.cookie[:], RandBytes(20))

	return packet
}

func RandBytes(size uint) []byte {
	cookie := make([]byte, 20)
	_, err := rand.Read(cookie)

	// Then append a cookie that we expect to be in the pong.
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil
	}

	return cookie
}
