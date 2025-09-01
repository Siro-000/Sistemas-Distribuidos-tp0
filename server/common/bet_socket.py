
import json
import struct

from server.common.utils import Bet

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
    
    def _recv_all(sock, n):
        """Lee exactamente n bytes o lanza error si la conexión se cierra antes."""
        data = b''
        while len(data) < n:
            packet = sock.recv(n - len(data))
            if not packet:
                raise ConnectionError("Socket cerrado antes de recibir todos los datos")
            data += packet
        return data
     
    def recibe_bet(self) ->Bet:
        raw_len = self.recv_all(self._socket, 4)
        if not raw_len:
            return
        
        msg_len = struct.unpack('>I', raw_len)[0]
        data = self.recv_all(self._socket, msg_len - len(data))
        
        bet = json.loads(data)
        
        return Bet(
            agency=self._ip_map_agenci_number.get_agency_number(self._socket.getpeername()[0]),  
            first_name=bet["nombre"],
            last_name=bet["apellido"],
            document=bet["dni"],
            birthdate=bet["nacimiento"],
            number=bet["numero"]
            )
        
    def confirm(self): 
        self._socket.sendall(bytes(4))
    
    def close(self):
        self._socket.close()