package mptproof

import (
	"errors"
)

type MPTProof struct {
	Index []byte
	Value []byte
	Proof ProofMap
}

type ProofMap struct {
	db map[string][]byte
}

func New() ProofMap {
	return ProofMap{
		db: make(map[string][]byte),
	}
}

func (pm *ProofMap) Put(key []byte, value []byte) error {
	if pm.db == nil {
		return errors.New("proofMap db is nil")
	}
	copiedBytes := make([]byte, len(value))
	copy(copiedBytes, value)
	pm.db[string(key)] = copiedBytes
	return nil
}

func (pm *ProofMap) Delete(key []byte) error {

	if pm.db == nil {
		return errors.New("proofMap db is nil")
	}
	delete(pm.db, string(key))
	return nil
}

func (pm *ProofMap) Has(key []byte) (bool, error) {

	if pm.db == nil {
		return false, errors.New("proofMap db is nil")
	}
	_, ok := pm.db[string(key)]
	return ok, nil
}

func (pm *ProofMap) Get(key []byte) ([]byte, error) {

	if pm.db == nil {
		return nil, errors.New("proofMap db is nil")
	}
	if entry, ok := pm.db[string(key)]; ok {
		copiedBytes := make([]byte, len(entry))
		copy(copiedBytes, entry)
		return copiedBytes, nil
	}
	return nil, errors.New("value not found")
}
