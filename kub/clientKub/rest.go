package clientKub

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/pkg/errors"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36"
)

type RestClient struct {
	BaseUrl   string
	Token     string
	csrfToken string

	client       *http.Client
	responseAuth *ResponseAuth
	cookie       []*http.Cookie
	jar          *Jar
}

//easyjson
type ResponseAuth struct {
	JweToken string   `json:"jweToken"`
	Errors   []string `json:"errors"`
}

//easyjson
type Token struct {
	Token string `json:"token"`
}

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func NewRestClient(url, token string, forceInsecure bool) *RestClient {
	jar := NewJar()

	return &RestClient{
		BaseUrl: url,
		Token:   token,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: forceInsecure},
			},
			Timeout: time.Duration(2 * time.Second),
			Jar:     jar,
		},
		responseAuth: &ResponseAuth{},
	}
}

func (rc *RestClient) CsrfLogin() (string, error) {
	req, err := http.NewRequest("GET", rc.BaseUrl+"/api/v1/csrftoken/login", strings.NewReader(`{"token":"`+rc.Token+`"}`))
	if err != nil {
		return "", err
	}
	resp, err := rc.client.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.Bytes()

	type CsrfToken struct {
		Token string `json:"token"`
	}

	var csrfToken CsrfToken

	err = json.Unmarshal(newStr, &csrfToken)
	if err != nil {
		return "", err
	}

	resp.Body.Close()
	rc.csrfToken = csrfToken.Token
	return csrfToken.Token, err
}

func (rc *RestClient) Login(csrfToken string) error {
	req, err := http.NewRequest(
		"POST",
		rc.BaseUrl+"/api/v1/login",
		strings.NewReader(`{"token":"`+rc.Token+`"}`),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-csrf-token", csrfToken)
	resp, err := rc.client.Do(req)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.Bytes()

	respAuth := &ResponseAuth{}

	err = json.Unmarshal(newStr, respAuth)
	if err != nil {
		return err
	}
	if respAuth.JweToken == "" {
		err = respAuth.UnmarshalJSON(newStr)
		if err != nil {
			return err
		}
	}

	if respAuth.JweToken == "" {
		return errors.New("failet to unmarshal jwt")
	}

	rc.responseAuth = respAuth
	resp.Body.Close()
	return nil
}

type LoginStatus struct {
	HeaderPresent bool `json:"headerPresent"`
	HttpsMode     bool `json:"httpsMode"`
	TokenPresent  bool `json:"tokenPresent"`
}

func (rc *RestClient) Status() (*LoginStatus, error) {
	var (
		err error
	)

	loginStatus := &LoginStatus{}

	req, err := rc.request("GET", "/api/v1/login/status")
	if err != nil {
		return loginStatus, err
	}

	err = json.Unmarshal(req, loginStatus)
	if err != nil {
		return loginStatus, err
	}

	return loginStatus, err
}

func (rc *RestClient) CsrfToken() error {
	var (
		err error
	)

	rq, err := rc.request("GET", "/api/v1/csrftoken/token")
	if err != nil {
		return err
	}

	type CsrfToken struct {
		Token string `json:"token"`
	}

	var csrfToken CsrfToken

	err = json.Unmarshal(rq, &csrfToken)
	if err != nil {
		return err
	}

	rc.csrfToken = csrfToken.Token

	return nil
}

func (rc *RestClient) UpdateRefreshToken() error {
	xx, err := rc.responseAuth.MarshalJSON()
	if err != nil {
		return err
	}

	zz := bytes.NewReader(xx)

	req, err := http.NewRequest("POST", rc.BaseUrl+"/api/v1/token/refresh", zz)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-csrf-token", rc.csrfToken)
	//req.Header.Set("jwetoken", rc.responseAuth.JweToken)
	resp, err := rc.client.Do(req)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bufByte := buf.Bytes()
	resp.Body.Close()

	respAuth := &ResponseAuth{}
	err = respAuth.UnmarshalJSON(bufByte)
	if err != nil {
		return err
	}
	rc.responseAuth = respAuth
	return nil
}

func (rc *RestClient) OverView() ([]byte, error) {
	var bufByte []byte
	req, err := http.NewRequest("GET", rc.BaseUrl+"/api/v1/overview?filterBy=&itemsPerPage=15&name=&page=1&sortBy=d,creationTimestamp", strings.NewReader(`{"token":"`+rc.Token+`"}`))
	if err != nil {
		return bufByte, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-csrf-token", rc.csrfToken)
	//req.Header.Set("jwetoken", rc.responseAuth.JweToken)

	resp, err := rc.client.Do(req)
	if err != nil {
		return bufByte, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bufByte = buf.Bytes()
	return bufByte, nil
}

func (rc *RestClient) Pod(nameSpace string) ([]byte, error) {
	return rc.request("GET", "/api/v1/pod/deploy?filterBy=&itemsPerPage=1000&name=&page=1&sortBy=d,creationTimestamp")
}

func (rc *RestClient) Deployment(nameSpace string) ([]byte, error) {
	return rc.request("GET", "/api/v1/deployment/deploy?filterBy=&itemsPerPage=1000&name=&page=1&sortBy=d,creationTimestamp")
}

func (rc *RestClient) Scale(nameDep string, scaleBy uint64) ([]byte, error) {
	return rc.request("PUT", fmt.Sprintf("/api/v1/scale/deployment/deploy/%v?scaleBy=%v", nameDep, scaleBy))
}

func (rc *RestClient) request(method, urlPath string) ([]byte, error) {
	var (
		bufByte []byte
		err     error
	)

	token := &Token{
		Token: rc.Token,
	}

	bt, err := token.MarshalJSON()
	if err != nil {
		return bufByte, err
	}

	//
	req, err := http.NewRequest(method, rc.BaseUrl+urlPath, bytes.NewReader(bt))
	if err != nil {
		return bufByte, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-csrf-token", rc.csrfToken)
	req.Header.Set("jwetoken", rc.responseAuth.JweToken)

	//for _, c := range rc.cookie {
	//	req.AddCookie(c)
	//}

	resp, err := rc.client.Do(req)
	if err != nil {
		return bufByte, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bufByte = buf.Bytes()
	resp.Body.Close()
	return bufByte, err
}
