package services

import "testing"

func TestHashCanonicalWithdrawTxNormalizesAddresses(t *testing.T) {
	storedTx := canonicalWithdrawTx{
		RequestId:       "787ef38b-e53c-4b9f-b596-6e96beaad6fd",
		BusinessTxId:    "yolo",
		ChainId:         "DappLinkBnbChain",
		From:            "0x0ddfaf0cbf9d926b73285af71395013f78b88bb5",
		To:              "0x1031876f710c558Eb86736EAd4b45A9aCB684677",
		Value:           "10000000000000000000",
		ContractAddress: "0x55d398326f99059fF775485246999027B3197955",
		TokenId:         "tokenId",
		TokenMeta:       "tokenMetaData",
	}
	requestTx := canonicalWithdrawTx{
		RequestId:       "787ef38b-e53c-4b9f-b596-6e96beaad6fd",
		BusinessTxId:    "yolo",
		ChainId:         "DappLinkBnbChain",
		From:            "0x0ddfaf0cbf9d926b73285af71395013f78b88bb5",
		To:              "0x1031876f710c558eb86736ead4b45a9acb684677",
		Value:           "10000000000000000000",
		ContractAddress: "0x55d398326f99059ff775485246999027b3197955",
		TokenId:         "tokenId",
		TokenMeta:       "tokenMetaData",
	}

	storedHash, err := hashCanonicalWithdrawTx(storedTx)
	if err != nil {
		t.Fatalf("hash stored tx: %v", err)
	}
	requestHash, err := hashCanonicalWithdrawTx(requestTx)
	if err != nil {
		t.Fatalf("hash request tx: %v", err)
	}

	if storedHash != requestHash {
		t.Fatalf("expected equal hashes after normalization, stored=%s request=%s", storedHash, requestHash)
	}
}
