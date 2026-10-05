package custody

import (
 "bytes"
 "context"
 "encoding/base64"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "strings"
 "time"
)

type TokenSource interface { Token(ctx context.Context)(string,error) }
type Digestor interface { Digest(ctx context.Context,network,rawTransaction string)([]byte,error) }
type SignatureAssembler interface { Assemble(ctx context.Context,network,rawTransaction string,derSignature []byte)(string,error) }

type CloudKMSSigner struct { Endpoint,KeyVersion string; Tokens TokenSource; Digestor Digestor; Assembler SignatureAssembler; HTTP *http.Client }

// CloudKMSSigner calls Google Cloud KMS asymmetricSign. Configure the KMS key
// with protection_level=HSM for hardware-backed custody. Private key material
// never enters this process. Digestor and Assembler must be chain-specific.
func(s *CloudKMSSigner) Sign(ctx context.Context,req SignRequest)(string,error){
 if s==nil||s.Tokens==nil||s.Digestor==nil||s.Assembler==nil{return "",errors.New("KMS signer dependencies missing")};if req.KeyRef==""||req.Network==""||req.RawTransaction==""{return "",errors.New("invalid signing request")}
 digest,err:=s.Digestor.Digest(ctx,req.Network,req.RawTransaction);if err!=nil{return "",err};if len(digest)==0{return "",errors.New("empty transaction digest")}
 token,err:=s.Tokens.Token(ctx);if err!=nil{return "",err};if token==""{return "",errors.New("empty KMS access token")}
 endpoint:=s.Endpoint;if endpoint==""{endpoint="https://cloudkms.googleapis.com"};key:=req.KeyRef;if s.KeyVersion!=""{key=s.KeyVersion};url:=strings.TrimRight(endpoint,"/")+"/v1/"+key+":asymmetricSign"
 body,_:=json.Marshal(map[string]any{"digest":map[string]string{"sha256":base64.StdEncoding.EncodeToString(digest)}})
 client:=s.HTTP;if client==nil{client=&http.Client{Timeout:15*time.Second}}
 httpReq,err:=http.NewRequestWithContext(ctx,http.MethodPost,url,bytes.NewReader(body));if err!=nil{return "",err};httpReq.Header.Set("Authorization","Bearer "+token);httpReq.Header.Set("Content-Type","application/json")
 resp,err:=client.Do(httpReq);if err!=nil{return "",err};defer resp.Body.Close();raw,readErr:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if readErr!=nil{return "",readErr};if resp.StatusCode/100!=2{return "",fmt.Errorf("KMS signing failed: status=%d body=%s",resp.StatusCode,string(raw))}
 var out struct{Signature string `json:"signature"`};if err:=json.Unmarshal(raw,&out);err!=nil{return "",err};if out.Signature==""{return "",errors.New("KMS returned empty signature")};sig,err:=base64.StdEncoding.DecodeString(out.Signature);if err!=nil{return "",err};return s.Assembler.Assemble(ctx,req.Network,req.RawTransaction,sig)
}
