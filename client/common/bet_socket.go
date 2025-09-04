package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

type BetSocket struct {
	conn net.Conn
}

func NewBetSocket(conn net.Conn) *BetSocket {
	return &BetSocket{conn: conn}
}

func (b *BetSocket) recvAll(n int) ([]byte, error) {
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

func (b *BetSocket) sendAll(data []byte) error {
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
	buf := make([]byte, 4+len(s))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(s)))
	copy(buf[4:], []byte(s))
	return buf
}

func (b *BetSocket) SendBetBatch(bets []*PostBetRequest) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(bets)))

	payload := buf
	for _, bet := range bets {
		payload = append(payload, encodeString(bet.FirstName)...)
		payload = append(payload, encodeString(bet.LastName)...)
		payload = append(payload, encodeString(bet.Document)...)
		payload = append(payload, encodeString(bet.Birthdate)...)

		numBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(numBuf, uint64(bet.Number))
		payload = append(payload, numBuf...)
	}

	return b.sendAll(payload)
}

func (b *BetSocket) SendEndOfBatch() error {
	return b.sendAll(make([]byte, 4))
}

func (b *BetSocket) RecibeConfirm() error {
	data, err := b.recvAll(1)
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

func (b *BetSocket) Close() {
	if b.conn != nil {
		b.conn.Close()
	}
}
