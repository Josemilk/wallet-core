# Events carousel

The Events section uses the same visual language as the exchange dashboard: deep navy background, purple/magenta highlights, neon green positive states, yellow accents, soft-glow cards and translucent bottom navigation.

## Carousel order

1. Crypto Price Prediction — chart icon — select an asset and predict the published target condition.
2. Football Result — football icon — predict the published fixture result.
3. Basketball Result — basketball icon — predict the published game result.
4. Tennis Match — tennis icon — predict the published match winner.
5. Seasonal Skill Quiz — quiz icon — answer a prepared seasonal question set.

Each card displays: title, short description, event countdown, entry price, prize amount and prize asset. The client reads these values from the authenticated Events API; no hard-coded monetary values are used.

## Admin model

The administrator selects a prebuilt template and supplies the event schedule/data. Reward configuration is isolated to `entry_price`, `prize_amount` and `prize_asset`. Template type, validation rules, result source and settlement rules are not editable per event.

Paid participation and prize distribution must remain disabled until the deployment's legal/compliance gate explicitly enables the relevant event type and jurisdiction. The backend must not infer legal eligibility from the client UI.

## Dashboard integration

`GET /v1/dashboard` returns authenticated portfolio, market quotes and event summaries. The Android client should use the same authenticated session for Home, Markets, History, Events, Trading and Wallet. Market quotes are sourced server-side so provider credentials never enter the APK.
