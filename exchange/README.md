# Exchange integration — phases 14–21

This branch adds the exchange/DEX integration boundary on top of `wallet-core`.

## Components
- `android/app/.../exchange`: exchange UI/domain boundary.
- `exchange/contracts`: AMM, factory and router contracts.
- `exchange/docs`: security and deployment gates.

## Custody rule
The Android self-custody wallet may sign locally. The exchange backend never receives seed phrases or private keys.

## Production gate
The code is production-oriented but **not declared safe for mainnet** until compilation, unit/integration/fuzz/invariant tests, testnet validation, bytecode verification, threat modelling and an independent smart-contract/security audit are complete.
