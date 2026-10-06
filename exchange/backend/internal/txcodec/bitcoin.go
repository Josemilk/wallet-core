package txcodec

import("bytes";"crypto/sha256";"encoding/binary";"encoding/hex";"errors";"golang.org/x/crypto/ripemd160")

type BitcoinInput struct{TxID string;Vout uint32;ScriptCode []byte;Amount uint64;Sequence uint32}
type BitcoinOutput struct{Value uint64;ScriptPubKey []byte}
type BitcoinTx struct{Version int32;Inputs []BitcoinInput;Outputs []BitcoinOutput;LockTime uint32;SighashType uint32}

// BIP143Digest returns the SegWit v0 signature digest for one input.
// ScriptCode must be the canonical scriptCode for the input (for P2WPKH this is the P2PKH script).
func BIP143Digest(tx BitcoinTx,index int)([]byte,error){if index<0||index>=len(tx.Inputs){return nil,errors.New("input index out of range")};if len(tx.Inputs)==0||len(tx.Outputs)==0{return nil,errors.New("transaction must have inputs and outputs")};sighash:=tx.SighashType;if sighash==0{sighash=1};prev:=bytes.Buffer{};seq:=bytes.Buffer{};for _,in:=range tx.Inputs{b,err:=hex.DecodeString(in.TxID);if err!=nil||len(b)!=32{return nil,errors.New("invalid input txid")};for i:=31;i>=0;i--{prev.WriteByte(b[i])};binary.Write(&prev,binary.LittleEndian,in.Vout);binary.Write(&seq,binary.LittleEndian,in.Sequence)};outs:=bytes.Buffer{};for _,o:=range tx.Outputs{binary.Write(&outs,binary.LittleEndian,o.Value);putVarBytes(&outs,o.ScriptPubKey)};hashPrev:=dsha(prev.Bytes());hashSeq:=dsha(seq.Bytes());hashOut:=dsha(outs.Bytes());in:=tx.Inputs[index];b,err:=hex.DecodeString(in.TxID);if err!=nil{return nil,errors.New("invalid input txid")};p:=bytes.Buffer{};binary.Write(&p,binary.LittleEndian,tx.Version);p.Write(hashPrev);p.Write(hashSeq);for i:=31;i>=0;i--{p.WriteByte(b[i])};binary.Write(&p,binary.LittleEndian,in.Vout);putVarBytes(&p,in.ScriptCode);binary.Write(&p,binary.LittleEndian,in.Amount);binary.Write(&p,binary.LittleEndian,in.Sequence);p.Write(hashOut);binary.Write(&p,binary.LittleEndian,tx.LockTime);binary.Write(&p,binary.LittleEndian,sighash);return dsha(p.Bytes()),nil}
func dsha(b []byte)[]byte{a:=sha256.Sum256(b);c:=sha256.Sum256(a[:]);return c[:]}
func putVarBytes(b *bytes.Buffer,v []byte){putVarInt(b,uint64(len(v)));b.Write(v)}
func putVarInt(b *bytes.Buffer,n uint64){switch{case n<0xfd:b.WriteByte(byte(n));case n<=0xffff:b.WriteByte(0xfd);binary.Write(b,binary.LittleEndian,uint16(n));case n<=0xffffffff:b.WriteByte(0xfe);binary.Write(b,binary.LittleEndian,uint32(n));default:b.WriteByte(0xff);binary.Write(b,binary.LittleEndian,n)}}


// BitcoinSegWitV0Assembler assembles a P2WPKH witness after a KMS signature.
type BitcoinSegWitV0Assembler struct{}
func(a BitcoinSegWitV0Assembler)Assemble(tx BitcoinTx,index int,derSignature,pubKey []byte)(string,error){
 if index<0||index>=len(tx.Inputs){return "",errors.New("input index out of range")}; in:=tx.Inputs[index]
 if len(in.ScriptCode)!=25||in.ScriptCode[0]!=0x76||in.ScriptCode[1]!=0xa9||in.ScriptCode[2]!=0x14||in.ScriptCode[23]!=0x88||in.ScriptCode[24]!=0xac{return "",errors.New("only P2PKH scriptCode for P2WPKH is supported")}
 if len(pubKey)!=33&&len(pubKey)!=65{return "",errors.New("invalid secp256k1 public key length")}
 h:=sha256.Sum256(pubKey); rh:=ripemd160.New(); _,_=rh.Write(h[:]); if !bytes.Equal(rh.Sum(nil),in.ScriptCode[3:23]){return "",errors.New("public key does not match input scriptCode")}
 digest,err:=BIP143Digest(tx,index);if err!=nil{return "",err};if !verifyDigest(pubKey,digest,derSignature){return "",errors.New("KMS signature does not verify against BIP143 digest")}
 if _,err:=parseECDSADER(derSignature);err!=nil{return "",err}; sh:=tx.SighashType;if sh==0{sh=1}; witnessSig:=append(append([]byte(nil),derSignature...),byte(sh)); witness:=[][]byte{witnessSig,append([]byte(nil),pubKey...)}
 return "0x"+hex.EncodeToString(serializeSegWitTx(tx,index,witness)),nil
}
func serializeSegWitTx(tx BitcoinTx,signedIndex int,witness [][]byte)[]byte{var b bytes.Buffer;binary.Write(&b,binary.LittleEndian,tx.Version);b.WriteByte(0);b.WriteByte(1);putVarInt(&b,uint64(len(tx.Inputs)));for _,in:=range tx.Inputs{prev,_:=hex.DecodeString(in.TxID);for i:=31;i>=0;i--{b.WriteByte(prev[i])};binary.Write(&b,binary.LittleEndian,in.Vout);b.WriteByte(0);binary.Write(&b,binary.LittleEndian,in.Sequence)};putVarInt(&b,uint64(len(tx.Outputs)));for _,o:=range tx.Outputs{binary.Write(&b,binary.LittleEndian,o.Value);putVarBytes(&b,o.ScriptPubKey)};for i:=range tx.Inputs{if i==signedIndex{putVarInt(&b,uint64(len(witness)));for _,item:=range witness{putVarBytes(&b,item)}}else{b.WriteByte(0)}};binary.Write(&b,binary.LittleEndian,tx.LockTime);return b.Bytes()}
