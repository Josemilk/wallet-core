# Production custody and blockchain connectivity

## Cloud KMS / HSM

The backend uses `custody.CloudKMSSigner`. Set `GOOGLE_OAUTH_ACCESS_TOKEN` only for controlled local integration tests. In production, run the service on GCP with Workload Identity/service-account credentials so `GoogleTokenSource` obtains a short-lived token from the metadata server.

Configure the KMS key reference as the `KeyRef` stored in the withdrawal queue. Use an asymmetric signing key whose protection level is `HSM`. The chain-specific `Digestor` must serialize and hash the exact unsigned transaction according to that chain. The chain-specific `SignatureAssembler` must attach the KMS DER signature using that chain's canonical signature rules. The generic signer deliberately does not hash arbitrary raw transaction text.

Private keys, seed phrases and decrypted key blobs are never accepted by the custody interface.

## EVM node

Set an HTTPS JSON-RPC URL in the service configuration and construct `blockchain.EVMRPC{URL: ...}`. The adapter supports `eth_blockNumber`, `eth_getBlockByNumber`, `eth_getTransactionByHash` and `eth_sendRawTransaction`. Wrap it with `EVMChainReader` for the confirmation engine.

## Bitcoin Core node

Set the Bitcoin Core RPC URL and credentials. `blockchain.BitcoinRPC` supports `getblockchaininfo`, `getblockhash`, `getblock`, `getrawtransaction` and `sendrawtransaction`. Wrap it with `BitcoinChainReader` for confirmation tracking.

## Safety gate

Mainnet credentials and production KMS key references must be supplied through the deployment secret manager, never committed to Git. The repository only contains adapters and configuration boundaries. Testnet must be the first environment in which the complete signing/broadcast/confirmation path is enabled.
