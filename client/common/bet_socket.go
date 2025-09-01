package common

import (
	"encoding/binary"
	"encoding/json"
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

func (b *BetSocket) SendBet(bet *PostBetRequest) error {
	data, err := json.Marshal(bet)
	if err != nil {
		return fmt.Errorf("error al serializar JSON: %w", err)
	}

	// Prefijo de longitud (4 bytes, big-endian)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	// Enviar longitud + data
	if err := b.sendAll(lenBuf); err != nil {
		return err
	}
	if err := b.sendAll(data); err != nil {
		return err
	}
	return nil
}

func (b *BetSocket) RecibeConfirm() error {
	data, err := b.recvAll(4)
	if err != nil {
		return err
	}

	// Verificar que sean todos ceros
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
