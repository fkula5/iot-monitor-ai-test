#!/bin/bash
(cd device-service && go run main.go) &
DEVICE_PID=$!
(cd auth-service && go run main.go) &
AUTH_PID=$!
(cd telemetry-service && go run main.go) &
TELEMETRY_PID=$!
(cd gateway && go run main.go) &
GATEWAY_PID=$!
(cd edge-simulator && go run main.go) &
EDGE_PID=$!

echo "All services started."
wait $DEVICE_PID $AUTH_PID $TELEMETRY_PID $GATEWAY_PID $EDGE_PID
