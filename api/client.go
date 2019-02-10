package api

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ken-tunc/aojtool/util"
)

var (
	cachePath = filepath.Join(util.HomeDir(), ".cache", "aojtool", "cookies")
	endpoint  = "https://judgeapi.u-aizu.ac.jp"
	userAgent = "aotjool"
)

type Client struct {
	Endpoint  *url.URL
	UserAgent string

	httpClient *http.Client
	cookieJar  *cookiejar.Jar
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

	return client, nil
}

func (c *Client) Cookies() []*http.Cookie {
	return c.cookieJar.Cookies(c.Endpoint)
}

func (c *Client) SaveCookies() error {
	path, err := util.EnsurePath(cachePath)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = gob.NewEncoder(buf).Encode(c.Cookies())
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveCookies() error {
	exist, err := util.Exists(cachePath)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	return os.Remove(cachePath)
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
	exist, err := util.Exists(cachePath)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	var cookies []*http.Cookie

	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	if err = gob.NewDecoder(buf).Decode(&cookies); err != nil {
		return nil, err
	}

	return cookies, nil
}
