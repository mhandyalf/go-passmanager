# go-passmanager

Password-vault service. Authentication and user credential management live in
[`handy-auth`](https://github.com/mhandyalf/handy-auth).

## Environment

Copy `.env.example` to `.env`. `AUTH_SERVICE_URL` must point to a running
`handy-auth` instance. Every `/api/passwords` request forwards its bearer token
to `GET {AUTH_SERVICE_URL}/api/auth/validate` before accessing vault data.

The existing registration, login, and password-reset paths are retained as
compatibility proxies to `handy-auth`. New services may call `handy-auth`
directly; the current frontend can keep using the existing API base URL.
