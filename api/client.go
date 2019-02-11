package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ken-tunc/aojtool/util"
)

var cookiesCache = filepath.Join(util.CacheDir, "cookies")

const (
	endpoint  = "https://judgeapi.u-aizu.ac.jp"
	userAgent = "aotjool"
)

type Client struct {
	Endpoint  *url.URL
	UserAgent string

	httpClient *http.Client
	cookieJar  *cookiejar.Jar

	Auth   *AuthService
	Submit *SubmitService
	Status *StatusService
}

func NewClient() (*Client, error) {
	parsedURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	cookieJar, err := newCookieJar(parsedURL)
	if err != nil {
		return nil, err
	}

	var client = &Client{
		Endpoint:  parsedURL,
		UserAgent: userAgent,

		httpClient: &http.Client{
			Jar:     cookieJar,
			Timeout: time.Duration(10) * time.Second,
		},
		cookieJar: cookieJar,
	}

	client.Auth = &AuthService{client: client}
	client.Submit = &SubmitService{Client: client}
	client.Status = &StatusService{Client: client}

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	ref := &url.URL{Path: path}
	u := c.Endpoint.ResolveReference(ref)

	var buf io.ReadWriter
	if payload != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v == nil {
		return resp, nil
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Cookies() []*http.Cookie {
	return c.cookieJar.Cookies(c.Endpoint)
}

func (c *Client) SaveCookies() error {
	byteCookies, err := util.Serialize(c.Cookies())
	if err != nil {
		return err
	}

	err = util.WriteBytes(byteCookies, cookiesCache)
	if err != nil {
		return err
	}

	err = os.Chmod(cookiesCache, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveCookies() error {
	exist, err := util.Exists(cookiesCache)
	if err != nil {
		return err
	}

	if exist {
		return os.Remove(cookiesCache)
	} else {
		return nil
	}
}

func newCookieJar(u *url.URL) (*cookiejar.Jar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies, err := loadCookies()
	if err != nil {
		return nil, err
	}

	if cookies != nil {
		jar.SetCookies(u, cookies)
	}

	return jar, nil
}

func loadCookies() ([]*http.Cookie, error) {
	exist, err := util.Exists(cookiesCache)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, nil
	}
	var cookies []*http.Cookie

	data, err := ioutil.ReadFile(cookiesCache)
	if err != nil {
		return nil, err
	}

	if err = util.Deserialize(data, &cookies); err != nil {
		return nil, err
	}

	return cookies, nil
}
