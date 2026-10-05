package txcodec

import("bytes";"crypto/sha256";"encoding/binary";"encoding/hex";"errors")

type BitcoinInput struct{TxID string;Vout uint32;ScriptCode []byte;Amount uint64;Sequence uint32}
type BitcoinOutput struct{Value uint64;ScriptPubKey []byte}
type BitcoinTx struct{Version int32;Inputs []BitcoinInput;Outputs []BitcoinOutput;LockTime uint32;SighashType uint32}

// BIP143Digest returns the SegWit v0 signature digest for one input.
// ScriptCode must be the canonical scriptCode for the input (for P2WPKH this is the P2PKH script).
func BIP143Digest(tx BitcoinTx,index int)([]byte,error){if index<0||index>=len(tx.Inputs){return nil,errors.New("input index out of range")};if len(tx.Inputs)==0||len(tx.Outputs)==0{return nil,errors.New("transaction must have inputs and outputs")};sighash:=tx.SighashType;if sighash==0{sighash=1};prev:=bytes.Buffer{};seq:=bytes.Buffer{};for _,in:=range tx.Inputs{b,err:=hex.DecodeString(in.TxID);if err!=nil||len(b)!=32{return nil,errors.New("invalid input txid")};for i:=31;i>=0;i--{prev.WriteByte(b[i])};binary.Write(&prev,binary.LittleEndian,in.Vout);binary.Write(&seq,binary.LittleEndian,in.Sequence)};outs:=bytes.Buffer{};for _,o:=range tx.Outputs{binary.Write(&outs,binary.LittleEndian,o.Value);putVarBytes(&outs,o.ScriptPubKey)};hashPrev:=dsha(prev.Bytes());hashSeq:=dsha(seq.Bytes());hashOut:=dsha(outs.Bytes());in:=tx.Inputs[index];b,err:=hex.DecodeString(in.TxID);if err!=nil{return nil,errors.New("invalid input txid")};p:=bytes.Buffer{};binary.Write(&p,binary.LittleEndian,tx.Version);p.Write(hashPrev);p.Write(hashSeq);for i:=31;i>=0;i--{p.WriteByte(b[i])};binary.Write(&p,binary.LittleEndian,in.Vout);putVarBytes(&p,in.ScriptCode);binary.Write(&p,binary.LittleEndian,in.Amount);binary.Write(&p,binary.LittleEndian,in.Sequence);p.Write(hashOut);binary.Write(&p,binary.LittleEndian,tx.LockTime);binary.Write(&p,binary.LittleEndian,sighash);return dsha(p.Bytes()),nil}
func dsha(b []byte)[]byte{a:=sha256.Sum256(b);c:=sha256.Sum256(a[:]);return c[:]}
func putVarBytes(b *bytes.Buffer,v []byte){putVarInt(b,uint64(len(v)));b.Write(v)}
func putVarInt(b *bytes.Buffer,n uint64){switch{case n<0xfd:b.WriteByte(byte(n));case n<=0xffff:b.WriteByte(0xfd);binary.Write(b,binary.LittleEndian,uint16(n));case n<=0xffffffff:b.WriteByte(0xfe);binary.Write(b,binary.LittleEndian,uint32(n));default:b.WriteByte(0xff);binary.Write(b,binary.LittleEndian,n)}}
