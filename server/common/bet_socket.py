from .utils import Bet
from typing import Tuple

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
    
    def _recv_all(self, n: int) -> bytes:
        data = b''
        while len(data) < n:
            packet = self._socket.recv(n - len(data))
            if not packet:
                raise ConnectionError("Socket cerrado antes de recibir todos los datos")
            data += packet
        return data

    def _recv_string(self) -> str:
        # Leer longitud de 4 bytes big-endian
        raw_len = self._recv_all(4)
        str_len = int.from_bytes(raw_len, "big")
        raw_str = self._recv_all(str_len)
        return raw_str.decode("utf-8")

    def recibe_bet(self) -> Bet:
        first_name = self._recv_string()
        last_name = self._recv_string()
        document = self._recv_string()
        birthdate = self._recv_string()

        # leer número (int64, 8 bytes)
        number_bytes = self._recv_all(8)
        number = int.from_bytes(number_bytes, "big")

        return Bet(
            agency=self._ip_map_agenci_number.get_agency_number(
                self._socket.getpeername()[0]
            ),
            first_name=first_name,
            last_name=last_name,
            document=document,
            birthdate=birthdate,
            number=number
        )

    def recibe_bet_batch(self) -> Tuple[list[Bet], int]:
        """Recibe un batch de apuestas"""
        raw_len = self._recv_all(4)
        batch_len = int.from_bytes(raw_len, "big")
        
        if batch_len == 0:
            return (None, 0)
        
        bets = []
        for _ in range(batch_len):
            bets.append(self.recibe_bet())  
        return (bets, batch_len)

    def confirm_batch(self):
        self._socket.sendall(bytes(1))

    def send_error_batch(self):
        self._socket.sendall(bytes([1]))
        
    def confirm(self):
        self._socket.sendall(bytes(1))
    
    def send_error(self):
        self._socket.sendall(bytes([1]))
        
    def close(self):
        self._socket.close()

