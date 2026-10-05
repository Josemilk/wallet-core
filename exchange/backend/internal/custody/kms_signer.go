package custody

import("bytes";"context";"encoding/base64";"encoding/json";"errors";"fmt";"hash/crc32";"io";"net/http";"strconv";"strings";"time")

type TokenSource interface{Token(ctx context.Context)(string,error)}
type Digestor interface{Digest(ctx context.Context,network,rawTransaction string)([]byte,error)}
type SignatureAssembler interface{Assemble(ctx context.Context,network,rawTransaction string,derSignature []byte,publicKeyPEM []byte)(string,error)}
type PublicKeySource interface{PublicKey(ctx context.Context,keyRef string)([]byte,error)}
type CloudKMSSigner struct{Endpoint,KeyVersion string;Tokens TokenSource;Digestor Digestor;Assembler SignatureAssembler;PublicKeys PublicKeySource;HTTP *http.Client;RequireHSM bool}

// CloudKMSSigner never handles private key material. For ECDSA transaction signing,
// the returned signature is verified against the public key before the assembler
// is allowed to construct a broadcastable transaction.
func(s *CloudKMSSigner)Sign(ctx context.Context,req SignRequest)(string,error){
 if s==nil||s.Tokens==nil||s.Digestor==nil||s.Assembler==nil||s.PublicKeys==nil{return "",errors.New("KMS signer dependencies missing")};if req.KeyRef==""||req.Network==""||req.RawTransaction==""{return "",errors.New("invalid signing request")}
 digest,err:=s.Digestor.Digest(ctx,req.Network,req.RawTransaction);if err!=nil{return "",err};if len(digest)!=32{return "",errors.New("transaction digest must be 32 bytes")}
 token,err:=s.Tokens.Token(ctx);if err!=nil{return "",err};if token==""{return "",errors.New("empty KMS access token")};endpoint:=s.Endpoint;if endpoint==""{endpoint="https://cloudkms.googleapis.com"};key:=req.KeyRef;if s.KeyVersion!=""{key=s.KeyVersion};url:=strings.TrimRight(endpoint,"/")+"/v1/"+key+":asymmetricSign"
 crc:=crc32.Checksum(digest,crc32.MakeTable(crc32.Castagnoli));body,_:=json.Marshal(map[string]any{"digest":map[string]string{"sha256":base64.StdEncoding.EncodeToString(digest)},"digestCrc32c":strconv.FormatUint(uint64(crc),10)})
 client:=s.HTTP;if client==nil{client=&http.Client{Timeout:15*time.Second}};httpReq,err:=http.NewRequestWithContext(ctx,http.MethodPost,url,bytes.NewReader(body));if err!=nil{return "",err};httpReq.Header.Set("Authorization","Bearer "+token);httpReq.Header.Set("Content-Type","application/json")
 resp,err:=client.Do(httpReq);if err!=nil{return "",err};defer resp.Body.Close();raw,readErr:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if readErr!=nil{return "",readErr};if resp.StatusCode/100!=2{return "",fmt.Errorf("KMS signing failed: status=%d",resp.StatusCode)}
 var out struct{Signature string `json:"signature"`;SignatureCRC string `json:"signatureCrc32c"`;Name string `json:"name"`;ProtectionLevel string `json:"protectionLevel"`;VerifiedDigestCRC bool `json:"verifiedDigestCrc32c"`};if err:=json.Unmarshal(raw,&out);err!=nil{return "",err};if out.Name!=""&&out.Name!=key{return "",errors.New("KMS returned an unexpected key version")};if !out.VerifiedDigestCRC{return "",errors.New("KMS did not verify digest integrity")};if s.RequireHSM&&out.ProtectionLevel!="HSM"&&out.ProtectionLevel!="HSM_SINGLE_TENANT"{return "",errors.New("configured KMS key is not HSM protected")};if out.Signature==""||out.SignatureCRC==""{return "",errors.New("KMS returned incomplete signature integrity fields")};sig,err:=base64.StdEncoding.DecodeString(out.Signature);if err!=nil{return "",err};expected,err:=strconv.ParseUint(out.SignatureCRC,10,32);if err!=nil{return "",errors.New("invalid KMS signature CRC")};if uint64(crc32.Checksum(sig,crc32.MakeTable(crc32.Castagnoli)))!=expected{return "",errors.New("KMS signature integrity check failed")}
 pub,err:=s.PublicKeys.PublicKey(ctx,key);if err!=nil{return "",err};if len(pub)==0{return "",errors.New("KMS public key unavailable")};return s.Assembler.Assemble(ctx,req.Network,req.RawTransaction,sig,pub)
}

type CloudKMSPublicKeySource struct{Endpoint string;Tokens TokenSource;HTTP *http.Client}
func(s *CloudKMSPublicKeySource)PublicKey(ctx context.Context,keyRef string)([]byte,error){if s==nil||s.Tokens==nil||keyRef==""{return nil,errors.New("KMS public-key source unavailable")};token,err:=s.Tokens.Token(ctx);if err!=nil{return nil,err};endpoint:=s.Endpoint;if endpoint==""{endpoint="https://cloudkms.googleapis.com"};u:=strings.TrimRight(endpoint,"/")+"/v1/"+keyRef+":getPublicKey";client:=s.HTTP;if client==nil{client=&http.Client{Timeout:10*time.Second}};req,err:=http.NewRequestWithContext(ctx,http.MethodGet,u,nil);if err!=nil{return nil,err};req.Header.Set("Authorization","Bearer "+token);resp,err:=client.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();body,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if err!=nil{return nil,err};if resp.StatusCode/100!=2{return nil,fmt.Errorf("KMS public key failed: status=%d",resp.StatusCode)};var out struct{Pem string `json:"pem"`;Name string `json:"name"`};if err:=json.Unmarshal(body,&out);err!=nil{return nil,err};if out.Name!=""&&out.Name!=keyRef{return nil,errors.New("unexpected KMS public-key version")};if out.Pem==""{return nil,errors.New("empty KMS public key")};return []byte(out.Pem),nil}
