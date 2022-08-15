package exec

import (
	"github.com/wooyang2018/corechain-sdk/code"
	"github.com/wooyang2018/corechain/protos"
)

var (
	_ code.Iterator = (*Iterator)(nil)
)

const MAX_ITERATOR_CAP = 100

// Iterator is a private struct
type Iterator struct {
	buf          []*protos.IteratorItem // current buffer of the kv items
	curBuf       *protos.IteratorItem   // pointer of current position
	curIdx       int                    // next index
	c            *contractContext       // where we can get the kv items
	end          bool
	err          error
	start, limit []byte
}

// newkvIterator return a code.Iterator
func newKvIterator(c *contractContext, start, limit []byte) code.Iterator {
	return &Iterator{
		start: start,
		limit: limit,
		c:     c,
	}
}

// load loads the data from xbrigde, called when buf is empty, maintains the curIdx and starter
func (ki *Iterator) load() {
	//clean the buf at beginning
	ki.buf = ki.buf[0:0]
	req := &protos.IteratorRequest{
		Start:  ki.start,
		Limit:  ki.limit,
		Cap:    MAX_ITERATOR_CAP,
		Header: &ki.c.header,
	}
	resp := new(protos.IteratorResponse)
	if err := ki.c.bridgeCallFunc("NewIterator", req, resp); err != nil {
		ki.err = err
		return
	}

	if len(resp.Items) == 0 {
		ki.end = true
		return
	}
	lastKey := resp.Items[len(resp.Items)-1].Key
	newStartKey := make([]byte, len(lastKey)+1)
	copy(newStartKey, lastKey)
	newStartKey[len(lastKey)] = 1

	ki.curIdx = 0
	ki.buf = resp.Items
	ki.start = newStartKey
}

func (ki *Iterator) Key() []byte {
	return ki.curBuf.Key
}

func (ki *Iterator) Value() []byte {
	return ki.curBuf.Value
}

func (ki *Iterator) Next() bool {
	if ki.end || ki.err != nil {
		return false
	}
	//永远保证有数据
	if ki.curIdx >= len(ki.buf) {
		ki.load()
		if ki.end || ki.err != nil {
			return false
		}
	}

	ki.curBuf = ki.buf[ki.curIdx]
	ki.curIdx += 1
	return true
}

func (ki *Iterator) Error() error {
	return ki.err
}

func (ki *Iterator) Close() {}
