# go-siwe-ankr-demo

A minimal **Go backend** demo that implements **Sign-In With Ethereum (SIWE)** and issues a **JWT** for protected API access, then reads data from Ethereum via **Ankr RPC**.

## Features
- ✅ SIWE auth flow (nonce → message signature → verification)
- ✅ JWT session + auth middleware for protected routes
- ✅ Ethereum RPC integration (chain id + balance)
- ✅ Tiny single-file UI (`/`) to test with MetaMask

## Tech
- Go + chi (HTTP routing)
- SIWE (message parsing + signature verification)
- JWT (token issuance/verification)
- go-ethereum (RPC client)
- Ankr as the Ethereum provider

## Quick Start
1. Create `.env` from `.env.example` and fill in your values.
2. Run:
   ```bash
   go run ./cmd/api