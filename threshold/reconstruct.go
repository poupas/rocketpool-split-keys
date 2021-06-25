package threshold

import (
	"fmt"

	"github.com/herumi/bls-eth-go-binary/bls"
)

// Reconstructs the public key from a map of secret key shares
func ReconstructPublicKey(keyShares map[uint64]*bls.SecretKey) (*bls.PublicKey, error) {
	var ids []bls.ID
	var pubKeys []bls.PublicKey

	for i, keyShare := range keyShares {
		id := bls.ID{}
		err := id.SetDecString(fmt.Sprintf("%d", i))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
		pubKey := keyShare.GetPublicKey()
		pubKeys = append(pubKeys, *pubKey)
	}

	aggregatePubKey := bls.PublicKey{}
	err := aggregatePubKey.Recover(pubKeys, ids)
	return &aggregatePubKey, err
}
