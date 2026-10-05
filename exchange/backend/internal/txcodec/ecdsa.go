package txcodec

import (
 "crypto/sha256"
 "crypto/x509"
 "encoding/asn1"
 "encoding/pem"
 "errors"
 "math/big"
)

type ecdsaDER struct { R,S *big.Int }
func parseECDSADER(der []byte)(*ecdsaDER,error){var s ecdsaDER;rest,err:=asn1.Unmarshal(der,&s);if err!=nil||len(rest)!=0||s.R==nil||s.S==nil||s.R.Sign()<=0||s.S.Sign()<=0{return nil,errors.New("invalid ECDSA DER signature")};return &s,nil}

var secpP,_=new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F",16)
var secpN,_=new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141",16)
var secpGx,_=new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",16)
var secpGy,_=new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8",16)

type point struct{x,y *big.Int;inf bool}
func mod(v *big.Int)*big.Int{return new(big.Int).Mod(v,secpP)}
func add(a,b point)point{if a.inf{return b};if b.inf{return a};if a.x.Cmp(b.x)==0{if a.y.Cmp(b.y)==0{return dbl(a)};return point{inf:true}};m:=mod(new(big.Int).Mul(new(big.Int).Sub(b.y,a.y),new(big.Int).ModInverse(new(big.Int).Sub(b.x,a.x),secpP)));x:=mod(new(big.Int).Sub(new(big.Int).Sub(new(big.Int).Mul(m,m),a.x),b.x));y:=mod(new(big.Int).Sub(new(big.Int).Mul(m,new(big.Int).Sub(a.x,x)),a.y));return point{x:x,y:y}}
func dbl(a point)point{if a.inf||a.y.Sign()==0{return point{inf:true}};num:=new(big.Int).Mul(big.NewInt(3),new(big.Int).Mul(a.x,a.x));den:=new(big.Int).Lsh(new(big.Int).Set(a.y),1);inv:=new(big.Int).ModInverse(den,secpP);if inv==nil{return point{inf:true}};m:=mod(new(big.Int).Mul(num,inv));x:=mod(new(big.Int).Sub(new(big.Int).Mul(m,m),new(big.Int).Lsh(new(big.Int).Set(a.x),1)));y:=mod(new(big.Int).Sub(new(big.Int).Mul(m,new(big.Int).Sub(a.x,x)),a.y));return point{x:x,y:y}}
func mul(a point,k *big.Int)point{out:=point{inf:true};base:=a;n:=new(big.Int).Set(k);for n.Sign()>0{if n.Bit(0)==1{out=add(out,base)};base=dbl(base);n.Rsh(n,1)};return out}
func sqrtMod(v *big.Int)*big.Int{e:=new(big.Int).Add(secpP,big.NewInt(1));e.Rsh(e,2);r:=new(big.Int).Exp(v,e,secpP);if new(big.Int).Mul(r,r).Mod(new(big.Int).Mul(r,r),secpP).Cmp(new(big.Int).Mod(new(big.Int).Set(v),secpP))!=0{return nil};return r}
func parsePublicKey(raw []byte)(point,error){if p,_:=pem.Decode(raw);p!=nil{raw=p.Bytes};var spki struct{Algorithm struct{Algorithm asn1.ObjectIdentifier;Parameters asn1.RawValue} `asn1:"sequence"`;SubjectPublicKey asn1.BitString} ;if _,err:=asn1.Unmarshal(raw,&spki);err!=nil{return point{},err};b:=spki.SubjectPublicKey.Bytes;if len(b)!=33&&len(b)!=65{return point{},errors.New("unsupported public key encoding")};if b[0]==4&&len(b)==65{return point{x:new(big.Int).SetBytes(b[1:33]),y:new(big.Int).SetBytes(b[33:])},nil};if b[0]!=2&&b[0]!=3{return point{},errors.New("unsupported public key prefix")};x:=new(big.Int).SetBytes(b[1:]);rhs:=mod(new(big.Int).Add(new(big.Int).Exp(x,big.NewInt(3),secpP),big.NewInt(7)));y:=sqrtMod(rhs);if y==nil{return point{},errors.New("invalid secp256k1 public key")};if byte(y.Bit(0))!=b[0]&1{y=new(big.Int).Sub(secpP,y)};return point{x:x,y:y},nil}

func recoveryID(digest,der,pub []byte)(byte,error){sig,err:=parseECDSADER(der);if err!=nil{return 0,err};if sig.R.Cmp(secpN)>=0||sig.S.Cmp(secpN)>=0{return 0,errors.New("ECDSA scalar out of range")};expected,err:=parsePublicKey(pub);if err!=nil{return 0,err};z:=new(big.Int).SetBytes(digest);for rec:=0;rec<4;rec++{j:=rec/2;x:=new(big.Int).Mul(new(big.Int).Set(sig.R),big.NewInt(int64(j)));if x.Cmp(secpP)>=0{continue};rhs:=mod(new(big.Int).Add(new(big.Int).Exp(x,big.NewInt(3),secpP),big.NewInt(7)));y:=sqrtMod(rhs);if y==nil||int(y.Bit(0))!=(rec&1){continue};R:=point{x:x,y:y};if !mul(R,secpN).inf{continue};rinv:=new(big.Int).ModInverse(sig.R,secpN);if rinv==nil{continue};q:=mul(add(mul(R,sig.S),mul(point{x:secpGx,y:secpGy},new(big.Int).Neg(z))),rinv);if !q.inf&&q.x.Cmp(expected.x)==0&&q.y.Cmp(expected.y)==0{return byte(rec),nil}};return 0,errors.New("unable to recover ECDSA recovery id")}

func verifyDigest(pub, digest, der []byte)bool{sig,err:=parseECDSADER(der);if err!=nil{return false};q,err:=parsePublicKey(pub);if err!=nil{return false};if sig.R.Sign()<=0||sig.R.Cmp(secpN)>=0||sig.S.Sign()<=0||sig.S.Cmp(secpN)>=0{return false};z:=new(big.Int).SetBytes(digest);w:=new(big.Int).ModInverse(sig.S,secpN);if w==nil{return false};u1:=new(big.Int).Mul(z,w);u1.Mod(u1,secpN);u2:=new(big.Int).Mul(sig.R,w);u2.Mod(u2,secpN);p:=add(mul(point{x:secpGx,y:secpGy},u1),mul(q,u2));return !p.inf&&new(big.Int).Mod(p.x,secpN).Cmp(sig.R)==0}

func sha256Digest(b []byte)[]byte{h:=sha256.Sum256(b);return h[:]}
var _ = x509.MarshalPKIXPublicKey
