# Backend administration policy

The administration surface is backend-only. The Android exchange client does not expose admin routes, admin tokens, treasury controls, or custody controls.

## Roles

- `admin`: platform administration.
- `treasury_admin`: treasury operations.
- `events_admin`: event templates and scheduling.
- `auditor`: read-only audit and reconciliation.

Administrative actions require authenticated OIDC identity, RBAC authorization and an immutable audit record.

## Platform fees

The intended default platform fee is 10% (1000 basis points) for operation classes configured by the fee policy. The fee must be shown to the user before confirmation and must be posted to the ledger as a separate fee entry. The rate is configuration, not hard-coded business logic, so it can be changed through the protected backend policy surface.

## Pool / treasury

Pool balances are represented by dedicated ledger accounts. Adding assets to the pool and removing assets from the pool must be ledger transactions with idempotency keys and audit records. Production deployment should require independent approval for treasury outflows; the initiating administrator must not be able to approve their own outflow.

The personal destination address for a treasury outflow is not accepted from an Android client. It is supplied only through the protected backend administration surface and must pass risk and allow-list checks before custody signing.

## Unclaimed event prizes

A winning event entry is settled into a prize-escrow account first. The winner receives the prize by an explicit claim operation. A configurable 30-day claim deadline may be applied. After the deadline, an automated eligibility check can mark the prize forfeited only when the event terms and applicable law permit it. Any resulting transfer must be an auditable ledger transaction; administrators must not be able to arbitrarily debit an active customer account.

## Custody boundary

Treasury and prize operations ultimately use the same custody boundary as customer withdrawals: risk checks -> approval -> signing policy -> HSM/KMS signer -> broadcast -> confirmation. Private keys never enter the backend process or Android application.
