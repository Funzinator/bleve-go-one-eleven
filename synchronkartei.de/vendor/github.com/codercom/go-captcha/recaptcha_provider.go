package captcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

//GoogleRecaptchaVerifyEndpoint contains the endpoint to POST to when verifying a captcha
const GoogleRecaptchaVerifyEndpoint = "https://www.google.com/recaptcha/api/siteverify"

type recaptchaProvider struct {
	siteKey   string
	secretKey string
}

//NewRecaptchaProvider returns a recaptcha provider
func NewRecaptchaProvider(siteKey string, secretKey string) CaptchaProvider {
	return recaptchaProvider{
		siteKey:   siteKey,
		secretKey: secretKey,
	}
}

func (rp recaptchaProvider) Generate() string {
	buf := &bytes.Buffer{}
	buf.WriteString("<script src='https://www.google.com/recaptcha/api.js'></script>")
	buf.WriteString(fmt.Sprintf(`<div class="g-recaptcha" data-sitekey="%s"></div>`, rp.siteKey))
	return buf.String()
}

func (rp recaptchaProvider) Verify(response string, remoteIP string) (ok bool, err error) {
	payload := url.Values{
		"secret":   {rp.secretKey},
		"response": {response},
	}
	if len(remoteIP) > 0 {
		payload.Add("remoteIP", remoteIP)
	}

	resp, err := http.PostForm(GoogleRecaptchaVerifyEndpoint, payload)
	if err != nil {
		return false, err
	}
	//fmt.Printf("Posted: %v\n", payload)
	defer resp.Body.Close()

	var jsonResp struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return false, err
	}
	//fmt.Printf("Response %+v\n", jsonResp)
	if len(jsonResp.ErrorCodes) > 0 {
		return false, fmt.Errorf("Error codes: %v", strings.Join(jsonResp.ErrorCodes, ","))
	}
	return jsonResp.Success, nil
}
