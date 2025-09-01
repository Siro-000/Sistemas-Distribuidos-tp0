import json
import struct

from .utils import Bet


class IpMapAgenciNumber:
    def __init__(self):
        self.ip_to_agency = {}
        self.next_agency = 1 
    
    def get_agency_number(self, ip: str) -> int:
        if ip not in self.ip_to_agency:
            self.ip_to_agency[ip] = self.next_agency
            self.next_agency += 1
        return self.ip_to_agency[ip]


class BetSocket:
    def __init__(self, socket):
        self._socket = socket
        self._ip_map_agenci_number = IpMapAgenciNumber()
    
    def _recv_all(self, n):
        """Lee exactamente n bytes o lanza error si la conexión se cierra antes."""
        data = b''
        while len(data) < n:
            packet = self._socket.recv(n - len(data))
            if not packet:
                raise ConnectionError("Socket cerrado antes de recibir todos los datos")
            data += packet
        return data
     
    def recibe_bet(self) -> Bet:
        # Leer prefijo de 4 bytes (longitud del mensaje)
        raw_len = self._recv_all(4)
        msg_len = struct.unpack('>I', raw_len)[0]

        # Leer mensaje completo
        data = self._recv_all(msg_len)

        # Parsear JSON
        bet_json = json.loads(data)

        # Crear objeto Bet
        return Bet(
            agency=self._ip_map_agenci_number.get_agency_number(
                self._socket.getpeername()[0]
            ),  
            first_name=bet_json["nombre"],
            last_name=bet_json["apellido"],
            document=bet_json["dni"],
            birthdate=bet_json["nacimiento"],
            number=bet_json["numero"]
        )
        
    def confirm(self): 
        """Envía confirmación de 4 bytes cero al cliente."""
        self._socket.sendall(bytes(4))
    
    def close(self):
        """Cierra la conexión del socket."""
        self._socket.close()
