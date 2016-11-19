// Kademlia

package zif

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

const MaxBucketSize = 20

type DhtFile struct {
	entryCount int
	entries    [][]Entry
}

type DHTSave struct {
	Buckets     [][]*Entry
	LongBuckets [][]*Entry
}

type RoutingTable struct {
	LocalAddress Address
	Buckets      []*list.List
	LongBuckets  []*list.List
}

func (rt *RoutingTable) Setup(addr Address) {
	rt.LocalAddress = addr

	rt.Buckets = make([]*list.List, len(rt.LocalAddress.Bytes)*8)
	rt.LongBuckets = make([]*list.List, len(rt.LocalAddress.Bytes)*8)

	for i := 0; i < len(rt.LocalAddress.Bytes)*8; i++ {
		rt.Buckets[i] = list.New()
		rt.LongBuckets[i] = list.New()
	}
}

func (rt *RoutingTable) Save(filename string) error {
	save := DHTSave{
		make([][]*Entry, len(rt.LocalAddress.Bytes)*8),
		make([][]*Entry, len(rt.LocalAddress.Bytes)*8),
	}

	for n, b := range rt.Buckets {
		for i := b.Front(); i != nil; i = i.Next() {
			save.Buckets[n] = append(save.Buckets[n], i.Value.(*Entry))
		}
	}

	for n, b := range rt.LongBuckets {
		for i := b.Front(); i != nil; i = i.Next() {
			save.LongBuckets[n] = append(save.LongBuckets[n], i.Value.(*Entry))
		}
	}

	json, err := json.Marshal(save)

	log.Info(string(json))

	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = ioutil.WriteFile(filename, json, 0600)

	if err != nil {
		return err
	}

	return nil
}

func LoadRoutingTable(path string, addr Address) (*RoutingTable, error) {
	var ret RoutingTable
	ret.Setup(addr)
	var save DHTSave

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Info("Creating new routing table")
		return &ret, nil
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &save)

	if err != nil {
		return nil, err
	}

	for n, i := range save.Buckets {
		for _, j := range i {
			ret.Buckets[n].PushBack(j)
		}
	}

	for n, i := range save.LongBuckets {
		for _, j := range i {
			ret.LongBuckets[n].PushBack(j)
		}
	}

	return &ret, nil
}

func (rt *RoutingTable) NumPeers() int {
	return BucketSize(rt.LongBuckets) + BucketSize(rt.Buckets)
}

func BucketSize(bucket []*list.List) int {
	count := 0

	for _, b := range bucket {
		for i := b.Front(); i != nil; i = i.Next() {
			count += 1
		}
	}

	return count
}

func (rt *RoutingTable) UpdateBucket(buckets []*list.List, entry Entry) bool {
	if len(entry.ZifAddress.Bytes) < AddressBinarySize {
		return false
	}

	zero_count := entry.ZifAddress.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := buckets[zero_count]

	// TODO: Ping peers, starting from back. If none reply, remove them.
	// Ensures only active peers are stored.
	if bucket.Len() == MaxBucketSize {
		return false
	}

	var foundEntry *list.Element = nil
	for i := bucket.Front(); i != nil; i = i.Next() {
		if i.Value.(*Entry).ZifAddress.Equals(&entry.ZifAddress) {
			foundEntry = i
		}
	}

	if foundEntry == nil {
		bucket.PushFront(&entry)
	} else {
		bucket.MoveToFront(foundEntry)
	}

	return true
}

func (rt *RoutingTable) Update(entry Entry) bool {
	var success bool

	closest := rt.FindClosest(entry.ZifAddress, 1)
	if len(closest) == 1 {
		closest_entry := closest[0]

		dist_closest := closest_entry.ZifAddress.Xor(&entry.ZifAddress)
		dist_this := rt.LocalAddress.Xor(&entry.ZifAddress)

		if dist_this.Less(dist_closest) {
			success = rt.UpdateBucket(rt.LongBuckets, entry)
		}
	} else if len(closest) == 0 {
		success = rt.UpdateBucket(rt.LongBuckets, entry)
	}

	success = rt.UpdateBucket(rt.Buckets, entry)

	return success
}

func copyToEntrySlice(slice *[]*Entry, begin *list.Element, count int) {

	for i := begin; i != nil && len(*slice) < count; i = i.Next() {
		*slice = append(*slice, i.Value.(*Entry))
	}

}

func (rt *RoutingTable) FindClosestInBuckets(buckets []*list.List, target Address, count int) []*Entry {
	if len(target.Bytes) != AddressBinarySize {
		return nil
	}

	ret := make([]*Entry, 0, count)

	bucket_num := target.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := buckets[bucket_num]

	copyToEntrySlice(&ret, bucket.Front(), count)

	// If the bucket is not filled, look the the buckets either side.
	for i := 1; (bucket_num-i >= 0 || bucket_num+i <= len(target.Bytes)*8) &&
		len(ret) < count; i++ {

		if bucket_num-i >= 0 {
			bucket = buckets[bucket_num-i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

		if bucket_num+i < len(target.Bytes)*8 {
			bucket = buckets[bucket_num+i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

	}

	for _, e := range ret {
		e.distance = *e.ZifAddress.Xor(&target)
	}

	sort.Sort(Entries(ret))

	return ret
}

func (rt *RoutingTable) FindClosest(target Address, count int) []*Entry {
	entries := make([]*Entry, 0)

	entries = append(entries, rt.FindClosestInBuckets(rt.Buckets, target, count)...)
	entries = append(entries, rt.FindClosestInBuckets(rt.Buckets, target, count)...)

	sort.Sort(Entries(entries))

	// then remove duplicates, as the two bucket lists may contain the same
	// entries

	if len(entries) < 1 {
		return entries
	}

	last := entries[0]
	j := 1
	for i := 1; i < len(entries); i++ {
		if entries[i].ZifAddress.Equals(&last.ZifAddress) {
			continue
		}

		entries[j] = entries[i]
		last = entries[i]
		j++
	}

	return entries[:j]
}
