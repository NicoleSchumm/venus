// FETCHED FROM LOTUS: builtin/cron/state.go.template

package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"

	actorstypes "github.com/filecoin-project/go-state-types/actors"
	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/venus/venus-shared/actors"
	"github.com/filecoin-project/venus/venus-shared/actors/adt"

	cron6 "github.com/filecoin-project/specs-actors/v6/actors/builtin/cron"
)

var _ State = (*state6)(nil)

func load6(store adt.Store, root cid.Cid) (State, error) {
	out := state6{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make6(store adt.Store) (State, error) {
	out := state6{store: store}
	out.State = *cron6.ConstructState(cron6.BuiltInEntries())
	return &out, nil
}

type state6 struct {
	cron6.State
	store adt.Store
}

func (s *state6) GetState() interface{} {
	return &s.State
}

func (s *state6) ActorKey() string {
	return manifest.CronKey
}

func (s *state6) ActorVersion() actorstypes.Version {
	return actorstypes.Version6
}

func (s *state6) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
