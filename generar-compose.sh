#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Uso: $0 <archivo_salida> <cantidad_clientes>"
    exit 1
fi

ARCHIVO=$1
CANTIDAD=$2

echo "Generando $ARCHIVO con $CANTIDAD clientes..."
python3 mi-generador.py "$ARCHIVO" "$CANTIDAD"
echo "Archivo $ARCHIVO generado correctamente."