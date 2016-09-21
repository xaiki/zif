package zif

import (
	"bytes"
	"testing"
)

func CreateEntry() Entry {
	var lp LocalPeer
	lp.GenerateKey()
	lp.ZifAddress.Generate(lp.publicKey)

	var e Entry
	e.ZifAddress = lp.ZifAddress

	return e
}

func TestRoutingTableUpdate(t *testing.T) {
	var lp LocalPeer

	lp.GenerateKey()
	lp.ZifAddress.Generate(lp.publicKey)
	lp.RoutingTable.Setup(lp.ZifAddress)

	e1 := CreateEntry()
	lp.RoutingTable.Update(e1)

	zero_count := e1.ZifAddress.Xor(&lp.RoutingTable.LocalAddress).LeadingZeroes()

	if lp.RoutingTable.Buckets[zero_count].Len() == 0 {
		t.Error("Routing table bucket not updated")
	}

	lp.RoutingTable.Update(e1)

	if lp.RoutingTable.Buckets[zero_count].Len() > 1 {
		t.Error("Routing table added same entry more than once")
	}

	// Add another entry, ensure it is added correctly
	e2 := CreateEntry()
	lp.RoutingTable.Update(e2)

	zero_count = e2.ZifAddress.Xor(&lp.RoutingTable.LocalAddress).LeadingZeroes()

	if lp.RoutingTable.Buckets[zero_count].Len() == 0 {
		t.Error("Routing table bucket not updated")
	}

	// Ensure entries are properly moved to the front of the list
	lp.RoutingTable.Update(e1)
	zero_count = e1.ZifAddress.Xor(&lp.RoutingTable.LocalAddress).LeadingZeroes()

	addr := lp.RoutingTable.Buckets[zero_count].Front().Value.(*Entry).ZifAddress

	if !bytes.Equal(addr.Bytes, e1.ZifAddress.Bytes) {
		t.Error("Routing table bucket not updated")
	}
}

func TestRoutingTableFindClosest(t *testing.T) {
	var lp LocalPeer

	lp.GenerateKey()
	lp.ZifAddress.Generate(lp.publicKey)
	lp.RoutingTable.Setup(lp.ZifAddress)

	e1 := CreateEntry()
	e2 := CreateEntry()
	lp.RoutingTable.Update(e1)
	lp.RoutingTable.Update(e2)

	// Even though I ask for 10, will only get two.
	// This forces the RoutingTable to look in all buckets, testing iteration.
	closest := lp.RoutingTable.FindClosest(e1.ZifAddress, 10)

	if len(closest) != 2 {
		t.Error("Routing table results are of incorrect length")
	}

	if !bytes.Equal(closest[0].ZifAddress.Bytes, e1.ZifAddress.Bytes) {
		t.Error("Routing table search not correctly performed")
	}

	// then test looking for an entry *not* in the table, to get a close match
	e3 := CreateEntry()
	closest = lp.RoutingTable.FindClosest(e3.ZifAddress, 10)

	if len(closest) != 2 {
		t.Error("Routing table results are of incorrect length")
	}

	// make sure the list I get is correctly sorted, ie things at a lower index
	// are closer to what I'm looking for than those at a higher index
	dist1 := e3.ZifAddress.Xor(&closest[0].ZifAddress)
	dist2 := e3.ZifAddress.Xor(&closest[1].ZifAddress)

	if !dist1.Less(dist2) {
		t.Error("Routing table search results not correctly ordered")
	}
}
