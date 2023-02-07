## transaction hash

Transaction定义如下：
```go
// Transaction is an Ethereum transaction.
type Transaction struct {
	inner TxData //hash计算的主要数据   
	time  time.Time 
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}
```
交易的hash计算与获取用同一个函数来完成。
```go
// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
    // 如果已经写入hash，直接返回
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

    //计算hash
	var h common.Hash
    //判断交易类型，采用不同的方式hash
	if tx.Type() == LegacyTxType {//Legacy，EIP-2718议案提出之前使用的类型
		h = rlpHash(tx.inner)
	} else {//Access list transaction & EIP-1559
		h = prefixedRlpHash(tx.Type(), tx.inner)
	}
	tx.hash.Store(h)
	return h
}
```
这里以第二种为例说明：prefixedRlpHash传入的参数是
```go
// prefixedRlpHash writes the prefix into the hasher before rlp-encoding x.
// It's used for typed transactions.
func prefixedRlpHash(prefix byte, x interface{}) (h common.Hash) {
    // 使用的是sha3算法族
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
    // 写入交易的类型（uin8），作为前缀
	sha.Write([]byte{prefix})
    // 将tx.inner 使用rlp编码后，拼接到前缀后
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
```

其中的rlp编码以嵌套数组的形式存储数据，分为叶子与非叶子两种方式。

叶子长度超过一个字节需要添加前缀，前缀根据数据长度和和数据类型决定，非叶子节点也类似。
![image](https://raw.githubusercontent.com/HegelXu/mptProof/main/mptProof/RLP%20Encode.png)
> https://learnblockchain.cn/books/geth/part3/rlp.html

## DeriveSha

DeriveSha用于计算transaction.header 的ReceiptHash与TxHash，两者都实现了DerivableList interface，这里以计算TxHash为例：
```go
// DeriveSha creates the tree hashes of transactions and receipts in a block header.
func DeriveSha(list DerivableList, hasher TrieHasher) common.Hash {
	hasher.Reset()

	valueBuf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(valueBuf)

	var indexBuf []byte
    // 跳过第一个交易
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
        // 将txs的index编码
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
        // 对第i个交易进行rlp encode，不赘述
		value := encodeForDerive(list, i, valueBuf)
        // Update到TrieHasher上
		hasher.Update(indexBuf, value)
	}
    // 第一个交易
	if list.Len() > 0 {
		indexBuf = rlp.AppendUint64(indexBuf[:0], 0)
		value := encodeForDerive(list, 0, valueBuf)
		hasher.Update(indexBuf, value)
	}
    // 剩下的交易
	for i := 0x80; i < list.Len(); i++ {
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		value := encodeForDerive(list, i, valueBuf)
		hasher.Update(indexBuf, value)
	}
	return hasher.Hash()
}
```

这个函数另一个重要参数是hasher TrieHasher，Trie是eth用于组织数据的结构MPT，TrieHasher是用于计算DerivableList hash而定义的一个接口。
```go
type TrieHasher interface {
	Reset()
	Update([]byte, []byte)
	Hash() common.Hash
}

这里的Update()是在MPT树上插入节点，Hash是计算root hash。

```
一般来讲，eth上经常使用的有两种TrieHasher实现。
其中一种底层是由一个map与锁实现的，这种实现可以用于数据的put&get，包含了整棵trie：

```go
type Database struct {
	db   map[string][]byte
	lock sync.RWMutex
}
```
另外一种底层是StackTrie，它一般用于生成hash，不需要保存多余的数据，只需要最终结果：
```go
type StackTrie struct {
	owner    common.Hash    // the owner of the trie
	nodeType uint8          // node type (as in branch, ext, leaf)
	val      []byte         // value contained by this node if it's a leaf
	key      []byte         // key chunk covered by this (leaf|ext) node
	children [16]*StackTrie // list of children (for branch and exts)
	writeFn  NodeWriteFunc  // function for committing nodes, can be nil
}
```

