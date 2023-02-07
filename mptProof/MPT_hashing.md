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
    // 将inner 使用rlp编码后，拼接到前缀后
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
```


