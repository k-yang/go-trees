package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

type Node struct {
	Hash      string
	Direction string
}

const (
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
)

func SHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	bz := h.Sum(nil)
	return hex.EncodeToString(bz)
}

func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

/**
 * Finds the index of the hash in the leaf hash list of the Merkle tree
 * and verifies if it's a left or right child by checking if its index is
 * even or odd. If the index is even, then it's a left child, if it's odd,
 * then it's a right child.
 * @param {string} hash
 * @param {Array<Array<string>>} merkleTree
 * @returns {string} direction
 */
func GetLeafNodeDirectionInMerkleTree(hash string, merkleTree [][]string) string {
	hashIndex := -1
	for i, h := range merkleTree[0] {
		if h == hash {
			hashIndex = i
			break
		}
	}
	if hashIndex%2 == 0 {
		return LEFT
	}
	return RIGHT
}

/**
 * If the hashes length is not even, then it copies the last hashes and adds it to the
 * end of the array, so it can be hashed with itself.
 * @param {Array<string>} hashes
 */

func EnsureEven(hashes []string) []string {
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}
	return hashes
}

/**
 * Generates the merkle root of the hashes passed through the parameter.
 * Recursively concatenates pair of hash hashes and calculates each sha256 hash of the
 * concatenated hashes until only one hash is left, which is the merkle root, and returns it.
 * @param {Array<string>} hashes
 * @returns merkleRoot
 */
func GenerateMerkleRoot(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}
	hashes = EnsureEven(hashes)
	var combinedHashes []string

	for i := 0; i < len(hashes); i += 2 {
		hashPairConcatenated := hashes[i] + hashes[i+1]
		hash := SHA256([]byte(hashPairConcatenated))
		combinedHashes = append(combinedHashes, string(hash))
	}
	if len(combinedHashes) == 1 {
		return combinedHashes[0]
	}
	return GenerateMerkleRoot(combinedHashes)
}

/**
 * Calculates the merkle root using the merkle proof by concatenating each pair of
 * hash hashes with the correct tree branch direction (left, right) and calculating
 * the sha256 hash of the concatenated pair, until the merkle root hash is generated
 * and returned.
 * The first hash needs to be in the first position of this array, with its
 * corresponding tree branch direction.
 * @param {Array<node>} merkleProof
 * @returns {string} merkleRoot
 */
func GetMerkleRootFromMerkleProof(merkleProof []Node) string {
	if len(merkleProof) == 0 {
		return ""
	}

	var merkleRootFromProof string = merkleProof[0].Hash
	for i := 1; i < len(merkleProof); i++ {
		if merkleProof[i].Direction == RIGHT {
			hash := SHA256([]byte(merkleRootFromProof + merkleProof[i].Hash))
			merkleRootFromProof = string(hash)
		} else {
			hash := SHA256([]byte(merkleProof[i].Hash + merkleRootFromProof))
			merkleRootFromProof = string(hash)
		}
	}

	return merkleRootFromProof
}

/**
 * Creates a merkle tree, recursively, from the provided hash hashes, represented
 * with an array of arrays of hashes/nodes. Where each array in the array, or hash list,
 * is a tree level with all the hashes/nodes in that level.
 * In the array at position tree[0] (the first array of hashes) there are
 * all the original hashes.
 * In the array at position tree[1] there are the combined pair or sha256 hashes of the
 * hashes in the position tree[0], and so on.
 * In the last position (tree[tree.length - 1]) there is only one hash, which is the
 * root of the tree, or merkle root.
 * @param {Array<string>} hashes
 * @returns {Array<Array<string>>} merkleTree
 */
func GenerateMerkleTree(hashes []string) [][]string {
	if len(hashes) == 0 {
		return [][]string{}
	}
	tree := [][]string{hashes}
	var generate func([]string, [][]string) [][]string
	generate = func(hashes []string, tree [][]string) [][]string {
		if len(hashes) == 1 {
			return tree
		}
		hashes = EnsureEven(hashes)
		var combinedHashes []string
		for i := 0; i < len(hashes); i += 2 {
			hashPairConcatenated := hashes[i] + hashes[i+1]
			hash := SHA256([]byte(hashPairConcatenated))
			combinedHashes = append(combinedHashes, string(hash))
		}
		tree = append(tree, combinedHashes)
		// fmt.Println("tree: ", tree)
		return generate(combinedHashes, tree)
	}
	return generate(hashes, tree)
}

/**
 * Generates the merkle proof by first creating the merkle tree,
 * and then finding the hash index in the tree and calculating if it's a
 * left or right child (since the hashes are calculated in pairs,
 * hash at index 0 would be a left child, hash at index 1 would be a right child.
 * Even indices are left children, odd indices are right children),
 * then it finds the sibling node (the one needed to concatenate and hash it with the child node)
 * and adds it to the proof, with its direction (left or right)
 * then it calculates the position of the next node in the next level, by
 * dividing the child index by 2, so this new index can be used in the next iteration of the
 * loop, along with the level.
 * If we check the result of this representation of the merkle tree, we notice that
 * The first level has all the hashes, an even number of hashes.
 * All the levels have an even number of hashes, except the last one (since is the
 * merkle root)
 * The next level have half or less hashes than the previous level, which allows us
 * to find the hash associated with the index of a previous hash in the next level in constant time.
 * Then we simply return this merkle proof.
 * @param {string} hash
 * @param {Array<string>} hashes
 * @returns {Array<node>} merkleProof
 */
func GenerateMerkleProof(hash string, hashes []string) []Node {
	if hash == "" || len(hashes) == 0 {
		return nil
	}
	tree := GenerateMerkleTree(hashes)
	merkleProof := []Node{{hash, GetLeafNodeDirectionInMerkleTree(hash, tree)}}
	hashIndex := -1
	for i, h := range tree[0] {
		if h == hash {
			hashIndex = i
			break
		}
	}

	for level := 0; level < len(tree)-1; level++ {
		isLeftChild := hashIndex%2 == 0
		siblingDirection := LEFT
		if isLeftChild {
			siblingDirection = RIGHT
		}
		siblingIndex := hashIndex - 1
		if isLeftChild {
			siblingIndex = hashIndex + 1
		}
		siblingNode := Node{tree[level][siblingIndex], siblingDirection}
		merkleProof = append(merkleProof, siblingNode)
		hashIndex = hashIndex / 2
	}
	return merkleProof
}
