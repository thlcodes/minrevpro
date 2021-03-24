#!/bin/sh
PORT=${PORT:-8080}
DEBUG=${DEBUG:-false}
./app --port "$PORT" --host "$HOST" --target "$TARGET" --secret "$SECRET" --debug=$DEBUG --basepath "$BASEPATH"