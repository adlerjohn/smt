package smt

import (
	"hash"
)

// BadProofError is returned when an invalid Merkle proof is supplied.
type BadProofError struct{}

func (e *BadProofError) Error() string {
	return "bad proof"
}

// DeepSparseMerkleSubTree is a deep Sparse Merkle subtree for working on only a few leafs.
type DeepSparseMerkleSubTree struct {
	*SparseMerkleTree
}

// NewDeepSparseMerkleSubTree creates a new deep Sparse Merkle subtree on an empty MapStore.
func NewDeepSparseMerkleSubTree(ms MapStore, hasher hash.Hash, root []byte) *DeepSparseMerkleSubTree {
	smt := &SparseMerkleTree{
		th: *newTreeHasher(hasher),
		ms: ms,
	}

	smt.SetRoot(root)

	return &DeepSparseMerkleSubTree{SparseMerkleTree: smt}
}

// AddBranch adds a branch to the tree.
// These branches are generated by smt.ProveForRoot.
// If the proof is invalid, a BadProofError is returned.
func (dsmst *DeepSparseMerkleSubTree) AddBranch(proof SparseMerkleProof, key []byte, value []byte) error {
	result, updates := verifyProofWithUpdates(proof, dsmst.Root(), key, value, dsmst.th.hasher)
	if !result {
		return &BadProofError{}
	}

	for _, update := range updates {
		err := dsmst.ms.Set(update[0], update[1])
		if err != nil {
			return err
		}
	}

	return nil
}
