# go-passmanager

Password-vault service. Authentication and user credential management live in
[`handy-auth`](https://github.com/mhandyalf/handy-auth).

## Environment

Copy `.env.example` to `.env`. `AUTH_SERVICE_URL` must point to a running
`handy-auth` instance. Every `/api/passwords` request forwards its bearer token
to `GET {AUTH_SERVICE_URL}/api/auth/validate` before accessing vault data.

The existing endpoint paths for registration, login, and password reset have
moved to `handy-auth`; clients should send those requests to the auth service.
