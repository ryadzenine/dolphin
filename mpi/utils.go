package mpi

import (
	"bytes"
	"encoding/gob"
	"io"
)

func encodeData(toSend map[string]Versionable) (*bytes.Buffer, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(toSend)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

func decodeData(data io.Reader) (states map[string]Versionable, err error) {
	dec := gob.NewDecoder(data)
	err = dec.Decode(&states)
	if err != nil {
		return nil, err
	}
	return
}
func getNxtPrev(me string, hosts []string) (nxt string, prev string) {
	for i, v := range hosts {
		if v == me {
			nxt = hosts[(i+1)%len(hosts)]
			j := (i - 1) % len(hosts)
			if j < 0 {
				j = j + len(hosts)
			}
			prev = hosts[j]
		}
	}
	return
}
