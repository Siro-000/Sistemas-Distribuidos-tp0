#!/bin/bash
# validar-echo-server.sh

SERVER_CONTAINER="server"

SERVER_PORT=$(grep "^SERVER_PORT" server/config.ini | cut -d '=' -f2 | tr -d '[:space:]')

TEST_MSG="Hello"

# obtener la red que Compose creó
COMPOSE_NETWORK=$(docker network ls | grep testing_net | awk '{print $1}')

RESULT=$(docker run --rm --network $COMPOSE_NETWORK busybox sh -c "\
  echo '$TEST_MSG' | nc $SERVER_CONTAINER $SERVER_PORT -w 2")

# Comparamos la respuesta con el mensaje enviado
if [ "$RESULT" = "$TEST_MSG" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
