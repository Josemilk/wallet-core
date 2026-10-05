package com.trustwallet.core.app

import android.graphics.Color
import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.trustwallet.core.exchange.DashboardAction
import com.trustwallet.core.exchange.DashboardActionHandler
import com.trustwallet.core.exchange.ExchangeDashboardView

class MainActivity : AppCompatActivity(), DashboardActionHandler {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        window.statusBarColor = Color.rgb(5, 11, 35)
        window.navigationBarColor = Color.rgb(5, 11, 35)
        setContentView(ExchangeDashboardView(this, this))
    }

    override fun onDashboardAction(action: DashboardAction) {
        // UI routing is deliberately separated from transport/authentication.
        // Feature screens map to the already-defined domains:
        // Trading -> ExchangeRepository + MarketStream
        // Wallet -> SelfCustodyWallet
        // Portfolio/History -> balances + ledger + settlement
        // Markets -> market-data stream
        // Events -> authenticated WebSocket/audit stream
        // The backend transport is injected here once API authentication/configuration is available.
        if (action in setOf(
                DashboardAction.PROFILE,
                DashboardAction.NOTIFICATIONS,
                DashboardAction.SETTINGS,
                DashboardAction.BTC,
                DashboardAction.ETH,
                DashboardAction.SOL
            )) {
            Toast.makeText(this, action.name, Toast.LENGTH_SHORT).show()
        }
    }
}
