package repository

import "github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"

var r *repo

type repo struct {
	SSLCheckRepo sslcheck.SSLCheckRepository
}

func New() *repo {
	if r == nil {
		r = &repo{}
	}

	return r
}

func GetRepos() *repo {
	return r
}
