package com.trustwallet.core.exchange

/** Integer base units only. Never use Double/Float for monetary values. */
data class PriceLevel(
    val priceBaseUnits: Long,
    val quantityBaseUnits: Long
)

data class OrderBookSnapshot(
    val symbol: String,
    val sequence: Long,
    val bids: List<PriceLevel>,
    val asks: List<PriceLevel>
)

enum class OrderSide { BUY, SELL }

enum class OrderType { LIMIT, MARKET }

data class OrderRequest(
    val symbol: String,
    val side: OrderSide,
    val type: OrderType,
    val priceBaseUnits: Long?,
    val quantityBaseUnits: Long,
    val clientOrderId: String
)

data class SwapQuote(
    val chainId: Long,
    val tokenIn: String,
    val tokenOut: String,
    val amountIn: String,
    val amountOut: String,
    val minimumAmountOut: String,
    val priceImpactBps: Long,
    val router: String,
    val pool: String,
    val expiresAtEpochSeconds: Long
)
