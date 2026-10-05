# Blockchain signing pipeline

Implemented building blocks:

```text
unsigned transaction
      ↓
chain-specific digest codec
      ↓
Cloud KMS / HSM asymmetricSign
      ↓
chain-specific signature assembler
      ↓
RPC sendRawTransaction
```

## Bitcoin

`txcodec.BitcoinTx` implements BIP143 SegWit v0 digest generation. It requires the exact input `scriptCode`, UTXO amount, sequence, outputs, locktime and sighash type. This is deliberately explicit: the signer must never infer a UTXO amount or script.

## EVM

`txcodec.EVM1559Tx` implements EIP-1559 typed transaction RLP and Keccak-256 signing digest. `txcodec.EVMLegacyTx` implements EIP-155 signing digest. Access-list validation is included.

## Not yet silently assumed

A complete EVM sender still needs recovery-id derivation/verification from the configured KMS public key before assembling `r,s,yParity`. A complete Bitcoin sender still needs deterministic UTXO selection, fee policy, exact output scripts and public-key/witness assembly. These values must come from explicit chain policy and custody configuration, not guessed defaults.

Therefore the KMS signer remains an explicit `Digestor` + `SignatureAssembler` boundary. The repository must not mark Mainnet custody as enabled until those chain-specific assemblers have vector tests and Testnet broadcast/confirmation tests.
