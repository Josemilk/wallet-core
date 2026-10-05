package txcodec

import("errors";"math/big";"golang.org/x/crypto/sha3")

type EVM1559Tx struct{ChainID,Nonce,MaxPriorityFeePerGas,MaxFeePerGas,GasLimit *big.Int;To []byte;Value *big.Int;Data []byte;AccessList []AccessTuple}
type AccessTuple struct{Address []byte;StorageKeys [][]byte}

type EVMLegacyTx struct{Nonce,GasPrice,GasLimit,Value *big.Int;To []byte;Data []byte;ChainID *big.Int}

func EIP1559Digest(tx EVM1559Tx)([]byte,error){if err:=validateEVM(tx);err!=nil{return nil,err};items:=[][]byte{intBytes(tx.ChainID),intBytes(tx.Nonce),intBytes(tx.MaxPriorityFeePerGas),intBytes(tx.MaxFeePerGas),intBytes(tx.GasLimit),fixedBytes(tx.To,20),intBytes(tx.Value),tx.Data,encodeAccessList(tx.AccessList)};payload:=rlpListRaw(items);return keccak(append([]byte{0x02},payload...)),nil}
func LegacyEIP155Digest(tx EVMLegacyTx)([]byte,error){if tx.Nonce==nil||tx.GasPrice==nil||tx.GasLimit==nil||tx.Value==nil||tx.ChainID==nil{return nil,errors.New("missing legacy transaction field")};items:=[][]byte{intBytes(tx.Nonce),intBytes(tx.GasPrice),intBytes(tx.GasLimit),fixedBytes(tx.To,20),intBytes(tx.Value),tx.Data,intBytes(tx.ChainID),nil,nil};return keccak(rlpListRaw(items)),nil}
func validateEVM(tx EVM1559Tx)error{if tx.ChainID==nil||tx.Nonce==nil||tx.MaxPriorityFeePerGas==nil||tx.MaxFeePerGas==nil||tx.GasLimit==nil||tx.Value==nil{return errors.New("missing EIP-1559 field")};if tx.ChainID.Sign()<0||tx.Nonce.Sign()<0||tx.MaxPriorityFeePerGas.Sign()<0||tx.MaxFeePerGas.Sign()<0||tx.GasLimit.Sign()<0||tx.Value.Sign()<0{return errors.New("negative EIP-1559 field")};if len(tx.To)!=0&&len(tx.To)!=20{return errors.New("EVM recipient must be 20 bytes")};if tx.MaxFeePerGas.Cmp(tx.MaxPriorityFeePerGas)<0{return errors.New("max fee below priority fee")};for _,a:=range tx.AccessList{if len(a.Address)!=20{return errors.New("access-list address must be 20 bytes")};for _,k:=range a.StorageKeys{if len(k)!=32{return errors.New("access-list storage key must be 32 bytes")}}};return nil}
func encodeAccessList(list []AccessTuple)[]byte{items:=make([][]byte,0,len(list));for _,a:=range list{keys:=make([][]byte,0,len(a.StorageKeys));for _,k:=range a.StorageKeys{keys=append(keys,rlpBytes(k))};items=append(items,rlpListRaw([][]byte{a.Address,rlpListEncoded(keys)}))};return rlpListEncoded(items)}
func intBytes(v *big.Int)[]byte{if v==nil||v.Sign()==0{return nil};return v.Bytes()}
func fixedBytes(v []byte,n int)[]byte{if len(v)==0{return nil};if len(v)!=n{return v};return v}
func keccak(b []byte)[]byte{h:=sha3.NewLegacyKeccak256();h.Write(b);return h.Sum(nil)}
func rlpListRaw(items [][]byte)[]byte{enc:=make([][]byte,0,len(items));for _,v:=range items{enc=append(enc,rlpBytes(v))};return rlpListEncoded(enc)}
func rlpListEncoded(items [][]byte)[]byte{n:=0;for _,v:=range items{n+=len(v)};prefix:=rlpLen(0xc0,n);out:=make([]byte,0,len(prefix)+n);out=append(out,prefix...);for _,v:=range items{out=append(out,v...)};return out}
func rlpBytes(v []byte)[]byte{if len(v)==1&&v[0]<0x80{return []byte{v[0]}};if len(v)<=55{return append([]byte{byte(0x80+len(v))},v...)};p:=rlpLen(0xb7,len(v));return append(p,v...)}
func rlpLen(base,n int)[]byte{if n<=55{return []byte{byte(base+n)}};tmp:=big.NewInt(int64(n)).Bytes();out:=[]byte{byte(base+55+len(tmp))};return append(out,tmp...)}
