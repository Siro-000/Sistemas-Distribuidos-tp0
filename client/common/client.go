package common

import (
	"context"
	"encoding/csv"
	"net"
	"os"
	"os/signal"
	"strconv"
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

func (c *Client) ReadBatches(filePath string, batchSize int) (<-chan []*PostBetRequest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	ch := make(chan []*PostBetRequest)

	go func() {
		defer close(ch)
		defer file.Close()

		reader := csv.NewReader(file)
		batch := []*PostBetRequest{}

		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					if len(batch) > 0 {
						ch <- batch
					}
					break
				}
				log.Criticalf("error leyendo CSV: %v", err)
				break
			}

			number, _ := strconv.ParseInt(record[4], 10, 64)
			bet := &PostBetRequest{
				FirstName: record[0],
				LastName:  record[1],
				Document:  record[2],
				Birthdate: record[3],
				Number:    number,
			}

			batch = append(batch, bet)
			if len(batch) >= batchSize {
				ch <- batch
				batch = []*PostBetRequest{}
			}
		}
	}()

	return ch, nil
}

func (c *Client) StartClientLoop(csvPath string, batchSize int) {
	c.listenSignals()
	problem := false

	if err := c.createClientBetSocket(); err != nil {
		return
	}

	batchCh, err := c.ReadBatches(csvPath, batchSize)
	if err != nil {
		log.Fatalf("error al leer CSV: %v", err)
	}

	for batch := range batchCh {
		if err := c.bet_socket.SendBetBatch(batch); err != nil {
			log.Criticalf("action: send bet batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			problem = true
			break
		}

		if err := c.bet_socket.RecibeConfirm(); err != nil {
			log.Criticalf("action: recibeConfrim batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			problem = true
			break
		} else {
			log.Infof("action: recibeConfrim batch | result: success | client_id: %v | cantidad: %d", c.config.ID, len(batch))
		}
	}

	if !problem {
		c.bet_socket.SendEndOfBatch()
		log.Infof("action: send all batch | result: success | client_id: %v", c.config.ID)
	}

	c.bet_socket.Close()

}
