package captcha

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestCaptchaProvider(t *testing.T) {
	bindAddr := "localhost:23222"
	if len(os.Getenv("RECAPTCHA_SITEKEY")) == 0 || len(os.Getenv("RECAPTCHA_SECRETKEY")) == 0 {
		t.Errorf(("Set RECAPTCHA_SITEKEY and RECAPTCHA_SECRETKEY environment variables\n"))
		t.FailNow()
	}

	fmt.Printf("Go to http://%v and verify the captcha\n", bindAddr)

	recaptcha := NewRecaptchaProvider(os.Getenv("RECAPTCHA_SITEKEY"), os.Getenv("RECAPTCHA_SECRETKEY"))

	exitChan := make(chan int)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintf(w, `
        <form method="POST" action="/verify-captcha">
            %s
            <button type="submit">Submit</button>
        </form>
        
        `,
			recaptcha.Generate(),
		)
	})

	http.HandleFunc("/verify-captcha", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved request to /verify-captcha")
		response := r.PostFormValue("g-recaptcha-response")
		verified, err := recaptcha.Verify(response, "")

		if err != nil {
			t.Errorf("Error verifying captcha: %v", err)
		}

		if !verified {
			t.Errorf("Captcha is not correct")
		}

		exitChan <- 1
	})

	server := &http.Server{
		Addr:    bindAddr,
		Handler: http.DefaultServeMux,
	}

	go func() {
		t.Errorf("Error starting server: %v", server.ListenAndServe())
	}()
	<-exitChan
}
