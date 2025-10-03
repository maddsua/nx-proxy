package main

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/maddsua/nx-proxy/api_models"
)

type Service interface {
	ID() string
	BindAddrString() string
	ListenAndServe() error
	SetPeers(users []api_models.Peer) error
	Close() error
}

func NewServiceHub() *serviceHub {
	return &serviceHub{bindMap: map[string]Service{}}
}

type serviceHub struct {
	bindMap map[string]Service
	mtx     sync.Mutex
}

func (hub *serviceHub) ApplySlots(slots []api_models.Slot) error {

	hub.mtx.Lock()
	defer hub.mtx.Unlock()

	refreshSet := map[string]struct{}{}

	for _, slot := range slots {

		var nextService Service
		switch slot.Service {
		case "socks":
			//	todo: implement
			nextService = &dummyService{id: slot.ID, addr: slot.Bind}
		case "http":
			//	todo: implement
			nextService = &dummyService{id: slot.ID, addr: slot.Bind}
		default:
			slog.Error("ServiceHub: ApplySlots: Unsupported service",
				slog.String("type", slot.Service))
			continue
		}

		bindAddr := nextService.BindAddrString()
		refreshSet[bindAddr] = struct{}{}

		if service, has := hub.bindMap[bindAddr]; has {

			if IsSameService(service, nextService) {

				if err := service.SetPeers(slot.Peers); err != nil {
					slog.Error("ServiceHub: ApplySlots: Set peers",
						slog.String("slot", service.BindAddrString()),
						slog.String("id", service.ID()),
						slog.String("err", err.Error()))
				}

				continue
			}

			slog.Debug("ServiceHub: ApplySlots: Replacing slot service",
				slog.String("slot", service.BindAddrString()),
				slog.String("id", service.ID()))

			if err := service.Close(); err != nil {
				slog.Error("ServiceHub: ApplySlots: Replace slot: Service shutdown failed",
					slog.String("slot", service.BindAddrString()),
					slog.String("id", service.ID()),
					slog.String("err", err.Error()))
				continue
			}
		}

		if err := nextService.ListenAndServe(); err != nil {
			slog.Error("ServiceHub: ApplySlots: Service start",
				slog.String("slot", nextService.BindAddrString()),
				slog.String("id", nextService.ID()),
				slog.String("err", err.Error()))
			continue
		}

		hub.bindMap[bindAddr] = nextService
	}

	for key, service := range hub.bindMap {
		if _, has := refreshSet[key]; !has {

			if err := service.Close(); err != nil {
				slog.Error("ServiceHub: ApplySlots: Service shutdown failed",
					slog.String("slot", service.BindAddrString()),
					slog.String("id", service.ID()),
					slog.String("err", err.Error()))
				continue
			}

			delete(hub.bindMap, key)
		}
	}

	return nil
}

func (hub *serviceHub) Close() error {
	return hub.ApplySlots(nil)
}

func IsSameService(old Service, new Service) bool {

	var fingerprint = func(val Service) string {
		return fmt.Sprintf("%s=%T", val, val.ID())
	}

	return fingerprint(old) == fingerprint(new)
}

type dummyService struct {
	id   uuid.UUID
	addr string
}

func (s *dummyService) ID() string {
	return s.id.String()
}

func (s *dummyService) BindAddrString() string {
	return s.addr + "/tcp"
}

func (s *dummyService) ListenAndServe() error {
	return nil
}

func (s *dummyService) SetPeers(users []api_models.Peer) error {
	return nil
}

func (s *dummyService) Close() error {
	return nil
}
