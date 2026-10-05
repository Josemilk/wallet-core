package com.trustwallet.core.exchange

/** Actions exposed by the dashboard. Each action maps to an existing exchange/wallet domain. */
enum class DashboardAction {
    PROFILE, NOTIFICATIONS, SETTINGS,
    PORTFOLIO, MARKET_OVERVIEW, MARKET_DETAILS,
    MARKETS, TRADING, WALLET, EVENTS, HISTORY,
    TIME_1D, TIME_1W, TIME_1M, TIME_1Y, TIME_ALL,
    BTC, ETH, SOL
}

interface DashboardActionHandler {
    fun onDashboardAction(action: DashboardAction)
}

/**
 * Domain mapping used by the UI layer. Transport/auth implementations stay outside the view.
 * Trading -> ExchangeRepository/orderBook/placeOrder/cancelOrder.
 * Wallet -> SelfCustodyWallet/WalletAddress/SignedTransaction.
 * Portfolio/History -> PostgreSQL balances/ledger/settlement APIs.
 * Markets -> market-data stream and persisted ticks/candles.
 * Events -> authenticated WebSocket/audit/event stream.
 */
