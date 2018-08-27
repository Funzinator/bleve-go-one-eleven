package captcha

//CaptchaProvider is a service that generates captcha codes
type CaptchaProvider interface {
	//Generate returns an HTML string containing the captcha challenge
	Generate() string
	//Verify verifies a challenge response
	//RemoteIP may be empty
	Verify(challenge string, remoteIP string) (ok bool, err error)
}
