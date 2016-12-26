# Protocol

## Messages
Zif communicates using ``Message``s. These are simple structs containing a
``Header``, and ``Content``. The ``Header`` is two bytes that indicates what
sort of message it is, and ``Content`` is the data sent along with the
``Message``.  These are serialized as JSON (for maximum support), and can be
gzipped.

## Headers
Below I have filled out a table of ``Header`` values, these are the ones
presently accepted by the Zif protocol. ``libzif/proto/protocolinfo.go`` also
contains these.

Name | Bytes
-----|------
ProtoZif | 0x7a66
ProtoVersion | 0x0000
ProtoHeader | 0x0000
ProtoOk | 0x0001
ProtoNo | 0x0002
ProtoTerminate | 0x0003
ProtoCookie | 0x0004
ProtoSig | 0x0005
ProtoPing | 0x0006
ProtoPong | 0x0007
ProtoDone | 0x0008
ProtoSearch | 0x0101
ProtoRecent | 0x0102
ProtoPopular | 0x0103
ProtoRequestHashList | 0x0104
ProtoRequestPiece | 0x0105
ProtoRequestAddPeer | 0x0106
ProtoEntry | 0x0200
ProtoPosts | 0x0201
ProtoHashList | 0x0202
ProtoPiece | 0x0203
ProtoPost | 0x0204
ProtoDhtQuery | 0x0300
ProtoDhtAnnounce | 0x0301
ProtoDhtFindClosest | 0x0302


## Handshaking
Before Zif requests can be made, handshaking must occur. This is to ensure that
both peers are who they claim to be, and to make sure everyone has the proper
keys.

From this point onwards, the "client" shall refer to the peer initiating a 
connection, and the "server" shall refer to the peer receiving a connection.

Before any handshaking can begin, a client must send the "zif" bytes, which are
``0x7a66``, and the "version" bytes, which presently are ``0x0000``. They are
sent over the network as Big Endian. The former indicates that this is a Zif
connection, the latter is protocol version - if protocol versions do not match,
then the connection can be dropped. At some future point I am interested in
adding protocol extension negotiation, but that is for the future.

Once we know this is a Zif connection, we can begin to share ``Message``s. First
of all, the client will send a ``Message`` with a ``Header`` of ``ProtoHeader``.
The ``Content`` field should be set to the client's public key.

If the server has any issues reading this message, then it will respond with
``ProtoNo``, and handshaking along with the connection will be terminated.
Otherwise, the server will respond with a ``ProtoOk``.

Once the server has the client's public key, it will generate the client's
address from it. It will then generate 20 bytes (CSPRNG please) and send these
to the client. The ``Message`` will have a ``Header`` of ``ProtoCookie`` and the
``Content`` field will be the bytes.

The client should then use its private key to sign the bytes, and send a
``Message`` with a ``Header`` of ``ProtoSig`` and ``Content`` of the signature.

The server recieves this signature, and will verify it. If it is correct, then
the client is sent a ``ProtoOk`` message, otherwise they are sent a ``ProtoNo``
message. As before, in the case of a ``ProtoNo``, handshaking and the connection
are both closed.

Once the signature has been verified, then we know that the client is in possesion
of the private key, and for our purposes they are who they say they are.

We now know that the client is who they claim to be, but handshaking is not yet
complete.

The server will now send either a ``ProtoNo`` or a ``ProtoOk``, and as before,
a ``ProtoNo`` leads to termination while ``ProtoOk`` is a signal to continue.
The server will then go through the above process as the "client", and the
client will go through the above process as the "server".

Once handshaking is done, both client and server should be verified.
