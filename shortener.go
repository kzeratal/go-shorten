package main

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"
)

type URLShortener struct{}

func (s *URLShortener) shorten(str string) (string, error) {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum((nil))), nil
}