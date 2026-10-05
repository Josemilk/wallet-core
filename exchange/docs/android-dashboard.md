# Android exchange dashboard wiring

The Android home screen follows the supplied visual design and is intentionally separated from transport/authentication.

| UI element | Domain / code boundary |
| --- | --- |
| Portfolio Balance | PostgreSQL `balances`, accounts, ledger and settlement projections |
| Market Overview chart | `market_ticks`, `market_candles`, market-data stream |
| 1D / 1W / 1M / 1Y / ALL | market-data query/window selector |
| Top Cryptocurrencies / See all | market-data market list and market detail |
| BTC / ETH / SOL cards | market symbol selection -> `ExchangeRepository.orderBook()` / market stream |
| Profile | authenticated account/profile boundary |
| Notifications | authenticated event/WebSocket stream |
| More | settings/security boundary |
| Home | portfolio projection |
| Markets | market-data streams and candles |
| Trading | `ExchangeRepository`, order book, matching engine, risk and settlement |
| Wallet | `SelfCustodyWallet`, deposits, withdrawals and swap router |
| Events | authenticated WebSocket + audit/event stream |
| History | orders, trades, ledger, deposits and withdrawals |

The UI does not contain private keys, custodial signing keys, or database credentials. Self-custody operations remain behind `SelfCustodyWallet`; custodial withdrawals remain behind the backend custody signer.

The current Canvas dashboard is the visual/navigation layer. Live values must be supplied by an authenticated repository/data provider before production activation; the values drawn by the design are presentation defaults, not account balances.
