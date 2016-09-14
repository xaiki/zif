package main

import "errors"

// At the moment just query for the closest known peer
const ResolveListSize = 1

// This takes a Zif address as a string and attempts to resolve it to an entry.
// This may be fast, may be a little slower. Will recurse its way through as
// many Queries as needed, getting closer to the target until either it cannot
// be found or is found.
// Cannot be found if a Query returns nothing, in this case the address does not
// exist on the DHT. Otherwise we should get to a peer that either has the entry,
// or one that IS the peer we are hunting.

// Takes a string as the API will just be passing a Zif address as a string.
// May well change, I'm unsure really.
func (lp *LocalPeer) Resolve(addr string) (*Entry, error) {
	address := DecodeAddress(addr)

	// First, find the closest peers in our routing table.
	// Satisfying if we already have the address :D
	closest := lp.RoutingTable.FindClosest(address, ResolveListSize)[0]

	for {
		// Check the current closest known peers. First iteration this will be
		// the ones from our routing table.
		if closest == nil {
			return nil, errors.New("Address could not be resolved")
			// The first in the slice is the closest, if we have this entry in our table
			// then this will be it.
		} else if closest.ZifAddress.Equals(&address) {
			return closest, nil
		}

		var peer *Peer

		// If the peer is not already connected, then connect.
		if peer = lp.GetPeer(closest.ZifAddress.Encode()); peer == nil {

			peer = NewPeer(lp)
			err := peer.Connect(closest.PublicAddress)

			if err != nil {
				return nil, err
			}

			_, err = peer.ConnectClient()

			if err != nil {
				return nil, err
			}
		}

		client, results, err := peer.Query(closest.ZifAddress.Encode())
		closest = &results[0]
		defer client.Close()

		if err != nil {
			return nil, errors.New("Peer query failed: " + err.Error())
		}
	}

	return nil, errors.New("Address could not be resolved")
}
