package api

import (
	"math/big"

	"github.com/google/uuid"
)

type ModelTable struct {
	Slots []ModelSlot `json:"slots"`
}

type ModelSlot struct {
	ID      uuid.UUID   `json:"id"`
	Bind    string      `json:"bind"`
	Service string      `json:"service"`
	Peers   []ModelPeer `json:"peers"`
}

type ModelPeer struct {
	ID        uuid.UUID       `json:"id"`
	Username  string          `json:"username"`
	Password  string          `json:"password"`
	Bandwidth *ModelBandwidth `json:"bandwidth,omitempty"`
}

type ModelBandwidth struct {
	RX    int `json:"rx"`
	TX    int `json:"tx"`
	MinRX int `json:"min_rx"`
	MinTX int `json:"min_tx"`
}

type ModelMetrics struct {
	Deltas  []ModelDelta `json:"deltas"`
	Service ModelService `json:"service"`
}

type ModelDelta struct {
	SlotID uuid.UUID `json:"slot_id"`
	PeerID uuid.UUID `json:"peer_id"`
	RX     big.Int   `json:"rx"`
	TX     big.Int   `json:"tx"`
}

type ModelService struct {
	RunID      uuid.UUID              `json:"run_id"`
	Uptime     int64                  `json:"uptime"`
	DataVolume ModelServiceDataVolume `json:"data_volume"`
}

type ModelServiceDataVolume struct {
	TotalRX big.Int `json:"total_rx"`
	TotalTX big.Int `json:"total_tx"`
}
