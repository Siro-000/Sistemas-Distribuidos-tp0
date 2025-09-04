# TP0: Docker + Comunicaciones + Concurrencia

## Solucion ejercisio 6 
Para resolucion del ejercisio 6 el protocolo de mensajes se modifico un poco.

Ahora el cliente manda de a batch, es decir, en de avarias apuestas por mensajes. 

El cliente va a leer desde el csv una cierta cantidad de bets, con las cuales va a construir un batch. Para hacerlo, en los primeros 4 bytes pone la cantidad de apuestas. Luego completa iterativamente con apuestas, siguiendo la linea de serializacion de apuestas del ejercisio anterior.

El servidor va a ir leyendo del canal los batch, comenzando con los primeros 4 bytes para saber la cantidad de bets que hay, para luego des-estructurar de forma iterativa como se venia haciendo. 
Si se recibe bien, el servidor manda una confirmacion al cliente de un byte de 0s. Si hay algun error, le avisa y se corta la comunicacion. 

Para avisarle el cliente de que no hay mas batch para leer, el cliente le manda un batch al servidor de tamaño. El servidor lo recibe y sabe que hay que terminar la comunicacion. 
