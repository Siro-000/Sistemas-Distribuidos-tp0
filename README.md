# TP0: Docker + Comunicaciones + Concurrencia

## Solucion ejercisio 7 
Para resolucion del ejercisio 7 el protocolo de mensajes se modifico un poco.

Ahora el ejercisio propone dos operacion para el cliente por asi decirlo, `load_bets` y `get_winners`. Entonces para avisarle al sevidor que operacionse quiere hacer, se manda primero un mensaje de un byte, que representa el CODIGO de la operacion, y una vez enviado, recien ahi se ejecuta la operacion.

El `load_bets` es igual a lo que se hacia en el ejercisio 6.  

El `get_winners`, el cliente le pide al servidor los ganadores. El servidor le puede responder `CONFIRM`, `ERROR` o `WAIT` (esto respeta los codigos de respuesta de un byte que veniamos teneniedo, pero ahora se agrego `WAIT`). Si responde `WAIT`, el cliente cierra la conexion, hace sleep y despues de un tiempo vuelve a preguntar. 
El `ERROR` esta por si hubo algun fallo con alguno de los otros agencias y por lo tanto no se puede hacer la loteria. 

Desde la vista del servidor, una vez todas las agencias suben sus apuestas. Se hace el sorteo y se guarda por cada agencia los documentos ganadores. 

Cuando un cliente le solicita los ganadores, se le manda todos en un solo mensaje. El paquete se constuye con 4 bytes al princio para indicar la cantidad de ganadores, luego con el formato de mandar strings que se venia teniendo de mandar 4 bytes del tamaño y luego el mensaje. 

