# Production gates for phases 14–21

## Android
- Never log or serialize mnemonic/private-key material.
- Wallet Core signing stays behind `SelfCustodyWallet`.
- Verify transaction recipient, chain ID, token contract, amount and calldata before signing.
- Show exact/minimum output and slippage for swaps.
- Use authenticated WebSocket sessions for private exchange channels.

## Smart contracts
- Solidity 0.8.24, optimizer enabled.
- Use a tagged, audited OpenZeppelin release; do not copy library source into this repository.
- Add unit, fuzz and invariant tests.
- Verify pool/router/factory bytecode on the target explorer.
- Add deployment timelock/multisig policy for administrative ownership.
- Test pause/unpause and recovery procedures.
- Test fee accounting and LP share accounting.
- Test malicious ERC-20 behavior and non-standard return values.
- Test donation/sync scenarios and reserve invariants.
- Test price-impact and slippage protections.
- Complete an independent smart-contract audit before mainnet.

## Exchange integration
- CEX ledger remains separate from self-custody.
- No custodial key is shipped in the APK.
- DEX transactions are signed locally by the self-custody wallet.
- Backend must never trust client-supplied token metadata; resolve token addresses from an allowlist.

## Mainnet gate
This repository branch is not a declaration that the system is safe for real funds. Mainnet activation requires all gates above, successful testnet operation, reproducible builds and independent security review.
