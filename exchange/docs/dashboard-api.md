# Dashboard API

The Android dashboard must call the exchange backend, never CoinMarketCap directly.

## Authentication

Send `Authorization: Bearer <OIDC/JWT>` to `GET /v1/dashboard`.

The backend validates RS256 signature, issuer, audience, expiry and key id against the configured JWKS endpoint. Configure `OIDC_ISSUER`, `OIDC_AUDIENCE` and `OIDC_JWKS_URL` in the backend deployment. For Firebase Authentication, the issuer/audience must match the Firebase project's documented ID-token values and Google Secure Token JWKS endpoint.

## Market data

The backend uses CoinMarketCap server-side. Set `CMC_API_KEY` for the authenticated API; when it is absent, the provider can use CoinMarketCap's documented keyless `/public-api` endpoints for development/evaluation. The API key is never shipped in the APK.

Production asset IDs should be CoinMarketCap IDs rather than symbols. The current provider requests `/v3/cryptocurrency/quotes/latest` and converts to USD.

## Live updates

REST provides the initial dashboard snapshot. The existing authenticated WebSocket/market-data stream is the live update channel for ticks, order-book changes, account events and event notifications. The Android client should reconcile sequence numbers and refetch the REST snapshot after a gap.

## Required environment

`OIDC_ISSUER`
`OIDC_AUDIENCE`
`OIDC_JWKS_URL`
`CMC_API_KEY`

No secret belongs in the Android APK or repository.
