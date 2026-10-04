package com.trustwallet.core.exchange

interface ExchangeRepository {
    suspend fun orderBook(symbol: String): OrderBookSnapshot
    suspend fun placeOrder(request: OrderRequest): String
    suspend fun cancelOrder(orderId: String)
    suspend fun swapQuote(chainId: Long, tokenIn: String, tokenOut: String, amountIn: String): SwapQuote
}

interface MarketStream {
    fun subscribe(symbol: String, listener: (OrderBookSnapshot) -> Unit)
    fun unsubscribe(symbol: String, listener: (OrderBookSnapshot) -> Unit)
}
