#!/bin/bash
set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$DIR/environments/dev.json"

echo "=== UMAG Gateway API Tests ==="
echo ""

echo "[1/9] Running Auth tests..."
newman run "$DIR/collections/01-auth.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[2/9] Running Cashier tests..."
newman run "$DIR/collections/02-cashier.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[3/9] Running Owner tests..."
newman run "$DIR/collections/03-owner.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[4/9] Running Agent tests..."
newman run "$DIR/collections/04-agent.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[5/9] Running Negative tests..."
newman run "$DIR/collections/05-negative.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[6/9] Running E2E Product Lifecycle tests..."
newman run "$DIR/collections/06-e2e-product-lifecycle.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[7/9] Running E2E Agent Lifecycle tests..."
newman run "$DIR/collections/07-e2e-agent-lifecycle.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[8/9] Running RBAC tests..."
newman run "$DIR/collections/08-rbac.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "[9/9] Running Data Integrity tests..."
newman run "$DIR/collections/09-data-integrity.json" \
  -e "$ENV_FILE" \
  --export-environment "$ENV_FILE" \
  --color on

echo ""
echo "=== All tests completed ==="
