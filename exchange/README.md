# Exchange platform architecture

This branch integrates the existing Wallet Core Android/self-custody code with an exchange backend and DEX boundary.

## Production components
01 architecture: `exchange/backend`, `exchange/database`, `exchange/contracts`
02 custody: self-custody signing remains local; custodial signing uses `KeyRef` only
03 PostgreSQL + ledger: `exchange/database/migrations`
04 users/accounts: users and account tables
05 blockchain gateway: `backend/internal/blockchain`
06 custodial wallet: `backend/internal/custody`
07 deposits: database model + blockchain gateway boundary
08 withdrawals: database model + custody signer boundary
09 order book: deterministic in-memory engine, PostgreSQL persistence boundary
10 matching engine: price/time priority implementation
11 settlement: explicit settlement validation boundary
12 market data: ticks/candles model
13 WebSocket: authenticated streaming remains a deployment integration gate
14 Android Exchange UI: existing exchange package
15 Wallet Core / self-custody: signing boundary
16–20 DEX: contracts/factory/pool/router/LP
21 DEX UI: Android DEX state
22 risk engine: pre-trade deterministic limits
23 AI trading: strategy signal boundary
24 backtesting: strategy result model
25 paper trading: same order/risk interfaces without settlement
26 automated trading: signals never bypass risk gates
27 smart order router: venue quote selection
28 admin: privileged authenticated service, never APK secrets
29 monitoring: health/readiness endpoints and deployment observability boundary
30 security audit: mandatory pre-mainnet gate
31 testnet: mandatory integration environment
32 mainnet: disabled until all gates pass

## API/secret policy
No exchange API secret, RPC credential, custodial private key, mnemonic, seed phrase or KMS credential belongs in Android or source control. Production values are injected through the deployment secret manager.

## Mainnet policy
There is deliberately no automatic mainnet switch. Mainnet requires signed deployment approval after testnet, security audit, reconciliation and operational readiness gates are complete.
