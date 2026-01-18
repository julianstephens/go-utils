## v0.2.7
- adds greater than, less than, gte, lte, and equal assertion helpers to `tests`

## v0.2.6

- removes `Error.Code` from `httputil.response`
- adds support for additional details in `Error.Details` in `httputil.response`
- adds `ValidateUUID` to `validator`

## v0.2.5

- adds `PromptPassword` and `PromptPasswordWithValidation` to `cliutil` 
- adds `validator`

## v0.2.4

- adds `Exists` to `helpers`
- improves package documentation

## v0.2.3

- adds `security` package

## v0.2.2

- adds support for long-lived refresh tokens to `httputil.auth.JWTManager`
- adds cookie security helpers to `httputil.auth.JWTManager`

## v0.2.1

- adds `TokenExpiration` helpers to `httputil.auth.JWTManager`
- adds `NewEmpty` constructor for responder without hooks
- adds JWT middleware to `httputil.middleware`

## v0.2.0

- adds `tests`
- improves documentation

## v0.1.6

- fixes `httputil.response.ErrorWithStatus`

## v0.1.5

- refactors documentation
- fixes `httputil.response` error statuses
- fixes `dbutil` field mapper names

## v0.1.4

- adds `slices`
- adds `logger`
- adds `config`
- adds `generic`
- adds `jsonutil`
- adds `cliutil`
- adds `dbutil`

## v0.1.3

- adds `httputil.auth.utils`
- adds `httputil.middleware`
- removes unnecessary user identifier claims
- refactors `httputil.auth.GenerateToken` for simpler user identification

## v0.1.2

- fixes superfluous header write in `response.OK`
- adds `helpers.Default`
- adds `httputil.auth`

## v0.1.1

- moves `helpers` to root
- adds `httputil.request`

## v0.1.0

- adds mux response package
- adds common helpers
