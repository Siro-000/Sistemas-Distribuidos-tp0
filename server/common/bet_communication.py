from .utils import Bet
from typing import Optional, Tuple

ERROR_CODE = 1 
CONFIRM_MESSAJE = bytes(1)
BYTES_AMOUNT_OF_BETS = 4
BYTES_LEN_STRING = 4 
BYTES_LEN_NUMBER = 8
LOAD_BETS = 0 
GET_WINNERS = 1 

class IpMapAgenciNumber:
    def __init__(self):
        self.ip_to_agency = {}
        self.next_agency = 1 
    
    def get_agency_number(self, ip: str) -> int:
        if ip not in self.ip_to_agency:
            self.ip_to_agency[ip] = self.next_agency
            self.next_agency += 1
        return self.ip_to_agency[ip]

class BetCommunication:
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
        # Leer 4 bytes big-endian
        raw_len = self._recv_all(BYTES_LEN_STRING)
        str_len = int.from_bytes(raw_len, "big")
        
        raw_str = self._recv_all(str_len)
        return raw_str.decode("utf-8")

    def recibe_bet(self) -> Bet:
        first_name = self._recv_string()
        last_name = self._recv_string()
        document = self._recv_string()
        birthdate = self._recv_string()

        # leer número (int64, 8 bytes)
        number_bytes = self._recv_all(BYTES_LEN_NUMBER)
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

    def recibe_bet_batch(self) -> Tuple[list[Bet], int,  Optional[Exception]]:
        """Recibe un batch de apuestas"""
        try: 
            raw_len = self._recv_all(BYTES_AMOUNT_OF_BETS)
            batch_len = int.from_bytes(raw_len, "big")
        except Exception as e:
            return None, 0, e
        
        if batch_len == 0:
            return (None, 0, None)
        
        bets = []
        for i in range(batch_len):
            try:
                bets.append(self.recibe_bet())
            except Exception as e:
                return None, i, e  
                  
        return (bets, batch_len, None)

    def confirm_batch(self):
        self._socket.sendall(CONFIRM_MESSAJE)

    def send_error_batch(self):
        self._socket.sendall(bytes([ERROR_CODE]))
        
    def confirm(self):
        self._socket.sendall(CONFIRM_MESSAJE)
    
    def send_error(self):
        self._socket.sendall(bytes([ERROR_CODE]))
        
    def close(self):
        self._socket.close()

    
    def recibe_operacion(self): 
        byte = self._recv_all(1)
        return byte[0]
    
    def send_wait(self): 
        self._socket.sendall(bytes([2]))
    
    def send_winners(self, agency_winners: dict): 
        agency = self._ip_map_agenci_number.get_agency_number(
                self._socket.getpeername()[0]) #la agencia de la comunicaion actual
        
        winners = agency_winners[agency] #the documents winners of the agency
        
        
        
    
