import signal
import socket
import logging

from .bet_socket import BetSocket
from .utils import store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        self.running = True
        
        signal.signal(signal.SIGTERM, self._handle_sigterm)

    def _handle_sigterm(self, signum, frame):
        logging.info("action: received SIGTERM | result: in_progress")
        self.running = False
        
    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self.running:
            try: 
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except socket.timeout:
                continue 
        
        self._server_socket.close()
                
    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet_socket = BetSocket(client_sock)
            bet = bet_socket.recibe_bet()
            
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            
            bet_socket.confirm()
            
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            bet_socket.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
