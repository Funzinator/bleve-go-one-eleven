package captcha

import "fmt"

//DryCaptchaSecret is the secret used by the dry captcha provider
var DryCaptchaSecret = "TEST"

type dryCaptchaProvider struct {
}

func NewDryCaptchaProvider() CaptchaProvider {
	return dryCaptchaProvider{}
}

func (cs dryCaptchaProvider) Generate() string {
	return fmt.Sprintf("<h1> Enter %v </h1>", DryCaptchaSecret)
}

func (cs dryCaptchaProvider) Verify(challenge string, remoteIP string) (ok bool, err error) {
	return (challenge == DryCaptchaSecret), nil
}
