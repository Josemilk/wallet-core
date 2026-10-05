package com.trustwallet.core.exchange

import android.content.Context
import android.graphics.*
import android.graphics.drawable.ColorDrawable
import android.view.MotionEvent
import android.view.View
import kotlin.math.min

/**
 * Lightweight Canvas UI so the exchange dashboard does not add a new UI framework dependency.
 * The visual hierarchy follows the supplied design: navy background, purple/magenta accents,
 * neon green performance, glowing cards and a six-item bottom navigation.
 */
class ExchangeDashboardView(
    context: Context,
    private val actions: DashboardActionHandler
) : View(context) {

    private val bg = Color.rgb(5, 11, 35)
    private val panel = Color.rgb(10, 18, 48)
    private val panel2 = Color.rgb(13, 22, 57)
    private val white = Color.WHITE
    private val muted = Color.rgb(174, 181, 210)
    private val purple = Color.rgb(202, 116, 255)
    private val magenta = Color.rgb(239, 93, 229)
    private val green = Color.rgb(119, 255, 129)
    private val yellow = Color.rgb(255, 238, 58)
    private val red = Color.rgb(255, 91, 150)

    private val paint = Paint(Paint.ANTI_ALIAS_FLAG)
    private var tab = DashboardAction.PORTFOLIO
    private var selectedRange = DashboardAction.TIME_1D
    private var downX = 0f
    private var downY = 0f

    init {
        setBackgroundColor(bg)
        setLayerType(View.LAYER_TYPE_SOFTWARE, null)
    }

    private fun dp(v: Float): Float = v * resources.displayMetrics.density
    private fun text(size: Float, color: Int, bold: Boolean = false) {
        paint.style = Paint.Style.FILL
        paint.color = color
        paint.textSize = dp(size)
        paint.typeface = if (bold) Typeface.create(Typeface.DEFAULT, Typeface.BOLD) else Typeface.DEFAULT
        paint.clearShadowLayer()
    }

    private fun rr(c: Canvas, l: Float, t: Float, r: Float, b: Float, radius: Float, color: Int, stroke: Int? = null) {
        paint.style = if (stroke == null) Paint.Style.FILL else Paint.Style.STROKE
        paint.strokeWidth = dp(if (stroke == null) 1f else 1.2f)
        paint.color = stroke ?: color
        paint.setShadowLayer(if (stroke == null) 0f else dp(5f), 0f, 0f, stroke ?: color)
        c.drawRoundRect(dp(l), dp(t), dp(r), dp(b), dp(radius), dp(radius), paint)
        paint.clearShadowLayer()
    }

    override fun onDraw(c: Canvas) {
        super.onDraw(c)
        val scale = width / dp(390f)
        c.save()
        c.scale(scale, scale)
        drawBackground(c)
        when (tab) {
            DashboardAction.PORTFOLIO -> drawHome(c)
            DashboardAction.MARKETS -> drawMarkets(c)
            DashboardAction.TRADING -> drawTrading(c)
            DashboardAction.WALLET -> drawWallet(c)
            DashboardAction.EVENTS -> drawEvents(c)
            DashboardAction.HISTORY -> drawHistory(c)
            else -> drawHome(c)
        }
        drawBottomNav(c)
        c.restore()
    }

    private fun drawBackground(c: Canvas) {
        val g = LinearGradient(0f, 0f, 0f, dp(900f), bg, Color.rgb(8, 15, 43), Shader.TileMode.CLAMP)
        paint.shader = g
        c.drawRect(0f, 0f, width.toFloat(), height.toFloat(), paint)
        paint.shader = null
        paint.color = Color.rgb(22, 28, 74)
        paint.setShadowLayer(dp(18f), 0f, 0f, Color.rgb(46, 53, 140))
        c.drawCircle(dp(380f), dp(330f), dp(2f), paint)
        paint.clearShadowLayer()
    }

    private fun drawHeader(c: Canvas) {
        // profile
        paint.style = Paint.Style.STROKE
        paint.strokeWidth = dp(2f)
        paint.color = purple
        paint.setShadowLayer(dp(12f), 0f, 0f, purple)
        c.drawCircle(dp(48f), dp(50f), dp(22f), paint)
        paint.clearShadowLayer()
        c.drawCircle(dp(48f), dp(43f), dp(7f), paint)
        c.drawArc(dp(35f), dp(51f), dp(61f), dp(70f), 190f, 160f, false, paint)
        // notification
        paint.style = Paint.Style.FILL
        paint.color = panel2
        c.drawCircle(dp(320f), dp(50f), dp(22f), paint)
        text(24f, yellow)
        c.drawText("♧", dp(311f), dp(58f), paint)
        paint.color = yellow
        c.drawCircle(dp(331f), dp(39f), dp(4f), paint)
        // menu
        paint.color = panel2
        c.drawCircle(dp(370f), dp(50f), dp(22f), paint)
        text(23f, green)
        c.drawText("•••", dp(355f), dp(57f), paint)
    }

    private fun drawHome(c: Canvas) {
        drawHeader(c)
        text(18f, purple, true)
        c.drawText("Portfolio Balance", dp(28f), dp(115f), paint)
        paint.color = magenta
        paint.strokeWidth = dp(2f)
        c.drawRect(dp(28f), dp(121f), dp(105f), dp(123f), paint)
        text(42f, white, true)
        c.drawText("$24,381.42", dp(28f), dp(168f), paint)
        rr(c, 282f, 139f, 374f, 168f, 14f, Color.TRANSPARENT, green)
        text(15f, green, true)
        c.drawText("↗ +2.14% today", dp(291f), dp(159f), paint)

        rr(c, 16f, 190f, 374f, 465f, 12f, panel2, Color.rgb(65, 77, 155))
        text(18f, white)
        c.drawText("Market Overview", dp(30f), dp(220f), paint)
        drawRange(c, 250f, "1D", DashboardAction.TIME_1D, 290f)
        drawRange(c, 250f, "1W", DashboardAction.TIME_1W, 340f)
        drawRange(c, 250f, "1M", DashboardAction.TIME_1M, 390f)
        drawRange(c, 250f, "1Y", DashboardAction.TIME_1Y, 440f)
        drawRange(c, 250f, "ALL", DashboardAction.TIME_ALL, 490f)
        drawChart(c)

        text(13f, muted); c.drawText("Invested", dp(48f), dp(490f), paint)
        text(20f, muted); c.drawText("$20,250.00", dp(38f), dp(522f), paint)
        text(13f, muted); c.drawText("Profit", dp(169f), dp(490f), paint)
        text(20f, green, true); c.drawText("$4,131.42 ↗", dp(158f), dp(522f), paint)
        text(13f, muted); c.drawText("Day P&L", dp(290f), dp(490f), paint)
        text(20f, green, true); c.drawText("+$512.60 ↗", dp(278f), dp(522f), paint)

        text(21f, white, true); c.drawText("Top Cryptocurrencies", dp(28f), dp(574f), paint)
        text(14f, yellow, true); c.drawText("See all ›", dp(322f), dp(574f), paint)
        cryptoCard(c, 592f, "Bitcoin", "BTC", "$64,207.12", "+1.24%", green, "₿")
        cryptoCard(c, 680f, "Ethereum", "ETH", "$3,412.88", "-0.56%", red, "◆")
        cryptoCard(c, 768f, "Solana", "SOL", "$142.91", "+4.82%", yellow, "≋")
    }

    private fun drawRange(c: Canvas, y: Float, label: String, action: DashboardAction, x: Float) {
        val selected = selectedRange == action
        rr(c, x, y, x + if (label == "ALL") 44f else 42f, y + 28f, 14f, if (selected) green else Color.rgb(49, 70, 82))
        text(13f, if (selected) Color.BLACK else white, true)
        c.drawText(label, x + 9f, y + 19f, paint)
    }

    private fun drawChart(c: Canvas) {
        paint.style = Paint.Style.STROKE
        paint.strokeWidth = dp(1f)
        paint.color = Color.rgb(34, 52, 92)
        for (i in 0..4) c.drawLine(dp(28f), dp(275f + i * 42f), dp(362f), dp(275f + i * 42f), paint)
        paint.color = green
        paint.strokeWidth = dp(3f)
        paint.setShadowLayer(dp(8f), 0f, 0f, green)
        val p = Path()
        p.moveTo(dp(28f), dp(388f)); p.cubicTo(dp(60f), dp(402f), dp(61f), dp(350f), dp(92f), dp(365f))
        p.cubicTo(dp(122f), dp(380f), dp(125f), dp(326f), dp(155f), dp(352f))
        p.cubicTo(dp(186f), dp(379f), dp(197f), dp(289f), dp(232f), dp(302f))
        p.cubicTo(dp(264f), dp(315f), dp(268f), dp(364f), dp(302f), dp(340f))
        p.cubicTo(dp(327f), dp(321f), dp(335f), dp(358f), dp(362f), dp(292f))
        c.drawPath(p, paint); paint.clearShadowLayer()
        paint.style = Paint.Style.FILL
    }

    private fun cryptoCard(c: Canvas, y: Float, name: String, ticker: String, price: String, change: String, accent: Int, icon: String) {
        rr(c, 16f, y, 374f, y + 74f, 11f, panel, Color.rgb(41, 52, 106))
        paint.color = accent; paint.setShadowLayer(dp(12f),0f,0f,accent); c.drawCircle(dp(52f), dp(y+37f), dp(22f), paint); paint.clearShadowLayer()
        text(25f, white, true); c.drawText(icon, dp(42f), dp(y+45f), paint)
        text(18f, white, true); c.drawText(name, dp(80f), dp(y+32f), paint)
        text(13f, muted); c.drawText(ticker, dp(80f), dp(y+53f), paint)
        text(17f, white, true); c.drawText(price, dp(235f), dp(y+34f), paint)
        text(13f, accent, true); c.drawText(change, dp(318f), dp(y+54f), paint)
    }

    private fun drawMarkets(c: Canvas) {
        drawHeader(c); title(c, "Markets")
        marketRow(c, 150f, "BTC/USDT", "$64,207.12", "+1.24%", green)
        marketRow(c, 220f, "ETH/USDT", "$3,412.88", "-0.56%", red)
        marketRow(c, 290f, "SOL/USDT", "$142.91", "+4.82%", yellow)
        marketRow(c, 360f, "BNB/USDT", "$585.22", "+0.83%", green)
        marketRow(c, 430f, "XRP/USDT", "$0.61", "+2.07%", green)
    }

    private fun marketRow(c: Canvas, y: Float, symbol: String, price: String, change: String, accent: Int) {
        rr(c, 18f, y, 372f, y+56f, 10f, panel, Color.rgb(42,53,105))
        text(16f, white, true); c.drawText(symbol, dp(32f), dp(y+24f), paint)
        text(15f, white, true); c.drawText(price, dp(205f), dp(y+24f), paint)
        text(13f, accent, true); c.drawText(change, dp(300f), dp(y+44f), paint)
    }

    private fun drawTrading(c: Canvas) {
        drawHeader(c); title(c, "Trading")
        rr(c, 18f, 145f, 372f, 205f, 12f, panel2, purple)
        text(20f, white, true); c.drawText("BTC/USDT", dp(32f), dp(180f), paint)
        text(14f, muted); c.drawText("Order Book • Matching Engine", dp(32f), dp(199f), paint)
        text(14f, green, true); c.drawText("BUY", dp(45f), dp(250f), paint)
        text(14f, red, true); c.drawText("SELL", dp(300f), dp(250f), paint)
        text(15f, muted); c.drawText("Best ask", dp(45f), dp(285f), paint); c.drawText("Best bid", dp(250f), dp(285f), paint)
        text(19f, white, true); c.drawText("Live stream", dp(45f), dp(320f), paint)
        text(13f, green); c.drawText("Authenticated WebSocket", dp(45f), dp(342f), paint)
        actionCard(c, 390f, "Place order", "ExchangeRepository.placeOrder()", green)
        actionCard(c, 465f, "Cancel order", "ExchangeRepository.cancelOrder()", red)
    }

    private fun drawWallet(c: Canvas) {
        drawHeader(c); title(c, "Wallet")
        actionCard(c, 150f, "Self-custody", "Wallet Core • addresses()", purple)
        actionCard(c, 225f, "Deposit", "Blockchain watcher • confirmations", green)
        actionCard(c, 300f, "Withdraw", "Risk → queue → custody signer", yellow)
        actionCard(c, 375f, "Swap", "Router → AMM pool → signed transaction", magenta)
        text(13f, muted); c.drawText("Private keys remain inside Wallet Core.", dp(28f), dp(470f), paint)
    }

    private fun drawEvents(c: Canvas) {
        drawHeader(c); title(c, "Events")
        actionCard(c, 150f, "Market events", "Market-data stream / WebSocket", green)
        actionCard(c, 225f, "Account events", "Authenticated user channel", purple)
        actionCard(c, 300f, "Security events", "Append-only audit boundary", yellow)
    }

    private fun drawHistory(c: Canvas) {
        drawHeader(c); title(c, "History")
        actionCard(c, 150f, "Orders", "orders + matching", purple)
        actionCard(c, 225f, "Trades", "settlement_records", green)
        actionCard(c, 300f, "Ledger", "double-entry accounting", yellow)
        actionCard(c, 375f, "Deposits / withdrawals", "blockchain + custody state", magenta)
    }

    private fun title(c: Canvas, value: String) { text(25f, white, true); c.drawText(value, dp(28f), dp(115f), paint) }
    private fun actionCard(c: Canvas, y: Float, title: String, detail: String, accent: Int) {
        rr(c, 18f, y, 372f, y+58f, 11f, panel, accent)
        text(16f, white, true); c.drawText(title, dp(32f), dp(y+25f), paint)
        text(12f, muted); c.drawText(detail, dp(32f), dp(y+45f), paint)
    }

    private fun drawBottomNav(c: Canvas) {
        val y = 650f
        rr(c, 12f, y, 378f, 735f, 20f, Color.rgb(14,20,48), Color.rgb(49,57,102))
        val items = listOf(
            Triple("⌂", "Home", DashboardAction.PORTFOLIO),
            Triple("▥", "Markets", DashboardAction.MARKETS),
            Triple("⇄", "Trading", DashboardAction.TRADING),
            Triple("▣", "Wallet", DashboardAction.WALLET),
            Triple("▦", "Events", DashboardAction.EVENTS),
            Triple("◷", "History", DashboardAction.HISTORY)
        )
        val step = 60f
        items.forEachIndexed { i, item ->
            val x = 40f + i * step
            val active = tab == item.third
            text(22f, if (active) green else muted, true); c.drawText(item.first, dp(x), dp(y+34f), paint)
            text(10f, if (active) green else muted, active); c.drawText(item.second, dp(x-13f), dp(y+56f), paint)
        }
    }

    override fun onTouchEvent(event: MotionEvent): Boolean {
        val s = width / dp(390f)
        val x = event.x / s / resources.displayMetrics.density
        val y = event.y / s / resources.displayMetrics.density
        when (event.action) {
            MotionEvent.ACTION_DOWN -> { downX=x; downY=y; return true }
            MotionEvent.ACTION_UP -> {
                if (downY > 635f) {
                    val idx = ((x - 15f) / 60f).toInt().coerceIn(0,5)
                    val action = listOf(DashboardAction.PORTFOLIO, DashboardAction.MARKETS, DashboardAction.TRADING, DashboardAction.WALLET, DashboardAction.EVENTS, DashboardAction.HISTORY)[idx]
                    tab = action; actions.onDashboardAction(action); invalidate(); return true
                }
                if (tab == DashboardAction.PORTFOLIO) {
                    when {
                        y in 230f..285f -> { selectedRange = when { x < 330f && x >= 280f -> DashboardAction.TIME_1D; else -> selectedRange }; actions.onDashboardAction(selectedRange); invalidate() }
                        y in 580f..675f -> { when { y < 675f && x < 195f -> actions.onDashboardAction(DashboardAction.BTC); x < 300f -> actions.onDashboardAction(DashboardAction.ETH); else -> actions.onDashboardAction(DashboardAction.SOL) } }
                        y < 135f && x < 90f -> actions.onDashboardAction(DashboardAction.PROFILE)
                        y < 135f && x > 295f && x < 345f -> actions.onDashboardAction(DashboardAction.NOTIFICATIONS)
                        y < 135f && x >= 345f -> actions.onDashboardAction(DashboardAction.SETTINGS)
                    }
                }
                return true
            }
        }
        return true
    }
}
