# Zif

A distributed torrent sharing and indexing network.

# Why?

I believe in anyone being able to share any information they wish, with anyone - without censorship getting in the way. Bittorrent is an excellent file-transfer protocol, though the torrent/magnet link discovery has been centralised and flawed since its inception.

This is my attempt at a solution.

# What even is this?

- a queryable database of torrents and metadata
- a p2p network for the distribution of torrent info hashes and metadata

Peer discovery is done via an implementation of Kademlia, a DHT similar to the one Bittorrent uses. Peers are assigned addresses similar to Bitcoin addresses, except using Ed25519 and SHA3.

# Sounds cool, when can I use it?

Sometime in the future. As of the time of writing, I have written a large portion of the DHT and a fair amount of the database, there's just still a lot of work left to do to glue the two together. Plus I need to write a nice graphical client!
