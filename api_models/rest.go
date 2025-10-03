package api_models

import (
	"math/big"

	"github.com/google/uuid"
)

type Table struct {
	Slots []Slot `json:"slots"`
}

type Slot struct {
	ID      uuid.UUID `json:"id"`
	Bind    string    `json:"bind"`
	Service string    `json:"service"`
	Peers   []Peer    `json:"peers"`
}

type Peer struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Bandwidth *Bandwidth `json:"bandwidth,omitempty"`
}

type Bandwidth struct {
	RX    int `json:"rx"`
	TX    int `json:"tx"`
	MinRX int `json:"min_rx"`
	MinTX int `json:"min_tx"`
}

type Metrics struct {
	Deltas  []Delta `json:"deltas"`
	Service Service `json:"service"`
}

type Delta struct {
	SlotID uuid.UUID `json:"slot_id"`
	PeerID uuid.UUID `json:"peer_id"`
	RX     big.Int   `json:"rx"`
	TX     big.Int   `json:"tx"`
}

type Service struct {
	RunID      uuid.UUID         `json:"run_id"`
	Uptime     int64             `json:"uptime"`
	DataVolume ServiceDataVolume `json:"data_volume"`
}

type ServiceDataVolume struct {
	TotalRX big.Int `json:"total_rx"`
	TotalTX big.Int `json:"total_tx"`
}
