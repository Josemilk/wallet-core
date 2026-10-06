package custody

import("context";"encoding/json";"errors";"io";"net/http";"os";"strings";"time")

type GoogleTokenSource struct{HTTP *http.Client}
func(s GoogleTokenSource)Token(ctx context.Context)(string,error){
 if v:=strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"));v!=""{return v,nil}
 client:=s.HTTP;if client==nil{client=&http.Client{Timeout:5*time.Second}}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,"http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token",nil);if err!=nil{return "",err};req.Header.Set("Metadata-Flavor","Google")
 resp,err:=client.Do(req);if err!=nil{return "",err};defer resp.Body.Close();if resp.StatusCode/100!=2{return "",errors.New("google metadata token request failed")};raw,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if err!=nil{return "",err};var out struct{AccessToken string `json:"access_token"`};if err:=json.Unmarshal(raw,&out);err!=nil{return "",err};if out.AccessToken==""{return "",errors.New("google metadata returned empty access token")};return out.AccessToken,nil
}
