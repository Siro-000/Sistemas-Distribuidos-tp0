# TP0: Docker + Comunicaciones + Concurrencia

## Solucion ejercisio 5 
Para resolucion del ejercisio 5 el protocolo que se implemento es el siguiente: 

Consite en esquema de pasaje de mensajes mixto

El cliente empieza directamente enviandole la solicitud de apuesta, que tiene los campos nombre(str), apellido(str), documento(str), birhdate(str) y numero(int64). 

Para los campos str se envia primero en 4 bytes su tamaño, para luego leer la cantidad de bytes. 

El campo numero int64 tiene tamaño fijo. 

El servidor lee los primeros 4 bytes donde saca el tamaño del campo, lee el campo y repite. Del ultimo campo ya sabe el tamaño porque es fijo. 

Si sale todo okey, el servidor manda una confirmacion de 1 byte de todos 0. Y el cliente no cierra sesion hasta recibir la confirmacion. 

En el codigo actual no se hace nada, pero en un esenario de tolerancia a fallos ese codigo de confirmacion podria usarse para recibir errores y voler a enviar la apuesta.