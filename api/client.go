package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	apiEndpoint = "https://judgeapi.u-aizu.ac.jp"
	datEndpoint = "https://judgedat.u-aizu.ac.jp"
	userAgent   = "aotjool"
)

type Client struct {
	Endpoint  *url.URL
	UserAgent string

	httpClient *http.Client
	cookieJar  *cookiejar.Jar

	Auth   *AuthService
	Submit *SubmitService
	Status *StatusService
	Test   *TestService
}

func NewClient() (*Client, error) {
	apiURL, err := url.ParseRequestURI(apiEndpoint)
	if err != nil {
		return nil, err
	}

	datURL, err := url.ParseRequestURI(datEndpoint)
	if err != nil {
		return nil, err
	}

	cookieJar, err := newCookieJar(apiURL, datURL)
	if err != nil {
		return nil, err
	}

	var client = &Client{
		Endpoint:  nil,
		UserAgent: userAgent,

		httpClient: &http.Client{
			Jar:     cookieJar,
			Timeout: time.Duration(10) * time.Second,
		},
		cookieJar: cookieJar,
	}

	client.Auth = &AuthService{client}
	client.Submit = &SubmitService{client}
	client.Status = &StatusService{client}
	client.Test = &TestService{client}

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, endpoint, method, path string, payload interface{}) (*http.Request, error) {
	parsedURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	c.Endpoint = parsedURL

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

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErrors util.ApiErrors
		err = json.NewDecoder(resp.Body).Decode(&apiErrors)
		if err != nil {
			return fmt.Errorf("http error: %s", resp.Status)
		} else {
			return apiErrors
		}
	}

	if v == nil {
		return nil
	} else {
		return json.NewDecoder(resp.Body).Decode(v)
	}
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

func newCookieJar(urls ...*url.URL) (*cookiejar.Jar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies, err := loadCookies()
	if err != nil {
		return nil, err
	}

	if cookies != nil {
		for _, u := range urls {
			jar.SetCookies(u, cookies)
		}
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
