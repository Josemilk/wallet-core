package com.trustwallet.core.exchange

data class ExchangeScreenState(
    val symbol: String = "BTC/USDT",
    val orderBook: OrderBookSnapshot? = null,
    val selectedSide: OrderSide = OrderSide.BUY,
    val quantityBaseUnits: String = "",
    val priceBaseUnits: String = "",
    val openOrders: List<String> = emptyList(),
    val error: String? = null,
    val loading: Boolean = false
)

data class DexScreenState(
    val chainId: Long? = null,
    val tokenIn: String = "",
    val tokenOut: String = "",
    val amountIn: String = "",
    val quote: SwapQuote? = null,
    val slippageBps: Long = 50,
    val loading: Boolean = false,
    val error: String? = null
)
