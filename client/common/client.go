package common

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const WAIT_TIME = 1
const BUFER_SIZE = 1

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config     ClientConfig
	bet_socket *BetSocket
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) listenSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Infof("action: received signal %v | result: in_progress", sig)
		c.cancel()
		c.bet_socket.Close()
	}()
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientBetSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	c.bet_socket = NewBetSocket(conn)
	return nil
}

func (c *Client) StartClientLoop(bet PostBetRequest) {
	c.listenSignals()
	if err := c.createClientBetSocket(); err != nil {
		return
	}

	log.Infof("\nla apuesta es %v\n", bet)
	if err := c.bet_socket.SendBet(&bet); err != nil {
		log.Criticalf(
			"action: send bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if err := c.bet_socket.RecibeConfirm(); err != nil {
		log.Criticalf(
			"action: recibe confirm | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	c.bet_socket.Close()
	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Document, bet.Number)
}

// There is an autoincremental msgID to identify every message sent
// 	// Messages if the message amount threshold has not been surpassed
// Loop:
// 	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
// 		select {
// 		case <-c.ctx.Done():
// 			break Loop
// 		default:

// 			if err := c.createClientSocket(); err != nil {
// 				return
// 			}

// 			// TODO: Modify the send to avoid short-write
// 			fmt.Fprintf(c.conn, "[CLIENT %v] Message N°%v\n", c.config.ID, msgID)

// 			msg, err := bufio.NewReader(c.conn).ReadString('\n')
// 			c.conn.Close()

// 			if err != nil {
// 				if errors.Is(err, net.ErrClosed) {
// 					// signal
// 					continue
// 				}
// 				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
// 					c.config.ID, err)
// 				return
// 			}

// 			log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
// 				c.config.ID, msg)

// 			time.Sleep(c.config.LoopPeriod)
// 		}
// 	}
// 	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
// }
