package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

const BYTES_LEN_NUMBER = 8
const BYTES_LEN_STRING = 4
const BYTES_LEN_CONFIRM = 1
const BYTES_AMOUNT_OF_BETS = 4

type BetCommunication struct {
	conn net.Conn
}

func NewBetCommunication(conn net.Conn) *BetCommunication {
	return &BetCommunication{conn: conn}
}

func (b *BetCommunication) recvAll(n int) ([]byte, error) {
	data := make([]byte, n)
	totalRead := 0
	for totalRead < n {
		readNow, err := b.conn.Read(data[totalRead:])
		if err != nil {
			return nil, fmt.Errorf("error al leer del socket: %w", err)
		}
		if readNow == 0 {
			return nil, fmt.Errorf("socket cerrado antes de recibir todos los datos")
		}
		totalRead += readNow
	}
	return data, nil
}

func (b *BetCommunication) sendAll(data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		sentNow, err := b.conn.Write(data[totalSent:])
		if err != nil {
			return fmt.Errorf("error al enviar datos: %w", err)
		}
		totalSent += sentNow
	}
	return nil
}

func encodeString(s string) []byte {
	buf := make([]byte, BYTES_LEN_STRING+len(s))
	binary.BigEndian.PutUint32(buf[:BYTES_LEN_STRING], uint32(len(s)))
	copy(buf[BYTES_LEN_STRING:], []byte(s))
	return buf
}

func (b *BetCommunication) SendBetBatch(bets []*PostBetRequest) error {
	buf := make([]byte, BYTES_AMOUNT_OF_BETS)
	binary.BigEndian.PutUint32(buf, uint32(len(bets)))

	payload := buf
	for _, bet := range bets {
		payload = append(payload, encodeString(bet.FirstName)...)
		payload = append(payload, encodeString(bet.LastName)...)
		payload = append(payload, encodeString(bet.Document)...)
		payload = append(payload, encodeString(bet.Birthdate)...)

		numBuf := make([]byte, BYTES_LEN_NUMBER)
		binary.BigEndian.PutUint64(numBuf, uint64(bet.Number))
		payload = append(payload, numBuf...)
	}

	return b.sendAll(payload)
}

func (b *BetCommunication) SendEndOfBatch() error {
	return b.sendAll(make([]byte, 4))
}

func (b *BetCommunication) RecibeConfirm() error {
	data, err := b.recvAll(BYTES_LEN_CONFIRM)
	if err != nil {
		return err
	}
	for _, v := range data {
		if v != 0 {
			return fmt.Errorf("confirmación inválida recibida: %v", data)
		}
	}
	return nil
}

func (b *BetCommunication) Close() {
	if b.conn != nil {
		b.conn.Close()
	}
}
