package custody

import("bytes";"context";"encoding/base64";"encoding/json";"errors";"fmt";"hash/crc32";"io";"net/http";"strconv";"strings";"time")

type TokenSource interface{Token(ctx context.Context)(string,error)}
type Digestor interface{Digest(ctx context.Context,network,rawTransaction string)([]byte,error)}
type SignatureAssembler interface{Assemble(ctx context.Context,network,rawTransaction string,derSignature []byte)(string,error)}
type CloudKMSSigner struct{Endpoint,KeyVersion string;Tokens TokenSource;Digestor Digestor;Assembler SignatureAssembler;HTTP *http.Client;RequireHSM bool}

// CloudKMSSigner calls Google Cloud KMS asymmetricSign. With RequireHSM=true,
// the response must report HSM or HSM_SINGLE_TENANT protection. Private key
// material never enters this process. Digestor and Assembler are chain-specific.
func(s *CloudKMSSigner)Sign(ctx context.Context,req SignRequest)(string,error){
 if s==nil||s.Tokens==nil||s.Digestor==nil||s.Assembler==nil{return "",errors.New("KMS signer dependencies missing")};if req.KeyRef==""||req.Network==""||req.RawTransaction==""{return "",errors.New("invalid signing request")}
 digest,err:=s.Digestor.Digest(ctx,req.Network,req.RawTransaction);if err!=nil{return "",err};if len(digest)==0{return "",errors.New("empty transaction digest")}
 token,err:=s.Tokens.Token(ctx);if err!=nil{return "",err};if token==""{return "",errors.New("empty KMS access token")};endpoint:=s.Endpoint;if endpoint==""{endpoint="https://cloudkms.googleapis.com"};key:=req.KeyRef;if s.KeyVersion!=""{key=s.KeyVersion};url:=strings.TrimRight(endpoint,"/")+"/v1/"+key+":asymmetricSign"
 crc:=crc32.Checksum(digest,crc32.MakeTable(crc32.Castagnoli));body,_:=json.Marshal(map[string]any{"digest":map[string]string{"sha256":base64.StdEncoding.EncodeToString(digest)},"digestCrc32c":strconv.FormatUint(uint64(crc),10)})
 client:=s.HTTP;if client==nil{client=&http.Client{Timeout:15*time.Second}};httpReq,err:=http.NewRequestWithContext(ctx,http.MethodPost,url,bytes.NewReader(body));if err!=nil{return "",err};httpReq.Header.Set("Authorization","Bearer "+token);httpReq.Header.Set("Content-Type","application/json")
 resp,err:=client.Do(httpReq);if err!=nil{return "",err};defer resp.Body.Close();raw,readErr:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if readErr!=nil{return "",readErr};if resp.StatusCode/100!=2{return "",fmt.Errorf("KMS signing failed: status=%d body=%s",resp.StatusCode,string(raw))}
 var out struct{Signature string `json:"signature"`;Name string `json:"name"`;ProtectionLevel string `json:"protectionLevel"`;VerifiedDigestCRC bool `json:"verifiedDigestCrc32c"`};if err:=json.Unmarshal(raw,&out);err!=nil{return "",err};if out.Name!=""&&out.Name!=key{return "",errors.New("KMS returned an unexpected key version")};if !out.VerifiedDigestCRC{return "",errors.New("KMS did not verify digest integrity")};if s.RequireHSM&&out.ProtectionLevel!="HSM"&&out.ProtectionLevel!="HSM_SINGLE_TENANT"{return "",errors.New("configured KMS key is not HSM protected")};if out.Signature==""{return "",errors.New("KMS returned empty signature")};sig,err:=base64.StdEncoding.DecodeString(out.Signature);if err!=nil{return "",err};return s.Assembler.Assemble(ctx,req.Network,req.RawTransaction,sig)
}
