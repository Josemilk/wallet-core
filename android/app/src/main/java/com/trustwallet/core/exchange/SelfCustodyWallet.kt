package com.trustwallet.core.exchange

/**
 * Narrow boundary around Wallet Core.
 * Implementations must keep seed/private-key material inside the wallet layer.
 */
data class WalletAddress(val chain: String, val address: String)

data class UnsignedTransaction(
    val chain: String,
    val to: String,
    val valueBaseUnits: String,
    val dataHex: String? = null
)

data class SignedTransaction(val chain: String, val rawTransaction: String)

interface SelfCustodyWallet {
    fun addresses(): List<WalletAddress>
    fun sign(transaction: UnsignedTransaction): SignedTransaction
}

/**
 * Exchange code should only receive public addresses and signed transactions.
 * No API in this package accepts a mnemonic or private key.
 */
