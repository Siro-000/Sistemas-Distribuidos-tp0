import os
import signal
import socket
import logging

from .bet_communication import BetCommunication, LOAD_BETS, IpMapAgencyNumber
from .utils import has_won, load_bets, store_bets

TIMEOUT = 1 
AGENCY_NUMBER = int(os.getenv("NUM_AGENCY"))

class Server:
    def __init__(self, port, listen_backlog):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)# Initialize server socket
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(TIMEOUT)
        self.running = True
        
        self.lottery = False
        
        self.error = False
        self.amount_agency = 0 
        self.agency_winners = {1:[],2:[],3:[],4:[],5:[]}
        self._ip_map_agency_number = IpMapAgencyNumber()
        
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
                
                if self.amount_agency == AGENCY_NUMBER and not self.lottery: 
                    self.make_lottery()
                    self.lottery = True
            except socket.timeout:
                continue 
        
        self._server_socket.close()
    
    def make_lottery(self):
        for bet in load_bets(): 
            agency_id = int(bet.agency)
            if has_won(bet):
                self.agency_winners[agency_id].append(bet.document)
        logging.info("action: sorteo | result: success")
            
    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try: 
            bet_communication = BetCommunication(client_sock, self._ip_map_agency_number)
            operacion = bet_communication.recibe_operacion()
            if operacion == LOAD_BETS: 
                self.load_bets(bet_communication)
            else: 
                self.give_winners(bet_communication)
        except Exception as e: 
            self.error = True
            logging.error(f"action: handle client connection | result: fail | error: {e}")
        finally:
            bet_communication.close()

    def give_winners(self, bet_commuication: BetCommunication):
        if self.error: 
            bet_commuication.send_error()
            return
        elif self.amount_agency < AGENCY_NUMBER : 
            bet_commuication.send_wait()
            return
        else: 
            bet_commuication.confirm()
        
        bet_commuication.send_winners(self.agency_winners)
        
        
    def load_bets(self, bet_communication):
        try:
            bets, amount_bets, e = bet_communication.recibe_bet_batch()
            new_bets = amount_bets
             
            while bets is not None: 
                logging.info(f'action: apuesta_recibida  | result: success | cantidad: {new_bets}')
                store_bets(bets)
                bet_communication.confirm()
                        
                bets, new_bets, e = bet_communication.recibe_bet_batch()
                amount_bets += new_bets
                    
            if e: 
                self.error = True
                logging.error(f"action: receive_message | result: fail | cantidad: {new_bets} | Error: {e}")
            else:     
                logging.info(f'action: recibir apuestas  | result: success | cantidad: {amount_bets}')
                self.amount_agency += 1
                
                    
        except Exception as e:
            self.error = True
            logging.error(f"action: receive_message | result: fail | Error: {e}")
            bet_communication.send_error()

        

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
