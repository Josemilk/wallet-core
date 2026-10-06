package txcodec

import (
    "encoding/hex"
    "testing"
)

func FuzzParseECDSADER(f *testing.F) {
    for _, seed := range []string{"3006020101020101","3006020100020101","00",""} {
        b, _ := hex.DecodeString(seed)
        f.Add(b)
    }
    f.Fuzz(func(t *testing.T, b []byte) {
        _, _ = parseECDSADER(b)
    })
}

func FuzzDecodeHex(f *testing.F) {
    f.Add("0x")
    f.Add("0x00ff")
    f.Add("deadbeef")
    f.Fuzz(func(t *testing.T, s string) {
        _, _ = decodeHex(s)
    })
}

func FuzzBIP143Digest(f *testing.F) {
    f.Add("00")
    f.Add("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")
    f.Fuzz(func(t *testing.T, txid string) {
        tx := BitcoinTx{
            Version: 2,
            SighashType: 1,
            Inputs: []BitcoinInput{{TxID: txid, Vout: 0, ScriptCode: []byte{0x76,0xa9,0x14,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0x88,0xac}, Amount: 1, Sequence: 0xffffffff}},
            Outputs: []BitcoinOutput{{Value: 1, ScriptPubKey: []byte{0x51}}},
        }
        _, _ = BIP143Digest(tx, 0)
    })
}
