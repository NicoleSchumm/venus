// FETCHED FROM LOTUS: builtin/account/state.go.template

package account

import (
	"fmt"

	actorstypes "github.com/filecoin-project/go-state-types/actors"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/manifest"
	"github.com/filecoin-project/venus/venus-shared/actors"
	"github.com/filecoin-project/venus/venus-shared/actors/adt"

	account4 "github.com/filecoin-project/specs-actors/v4/actors/builtin/account"
)

var _ State = (*state4)(nil)

func load4(store adt.Store, root cid.Cid) (State, error) {
	out := state4{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make4(store adt.Store, addr address.Address) (State, error) {
	out := state4{store: store}
	out.State = account4.State{Address: addr}
	return &out, nil
}

type state4 struct {
	account4.State
	store adt.Store
}

func (s *state4) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state4) GetState() interface{} {
	return &s.State
}

func (s *state4) ActorKey() string {
	return manifest.AccountKey
}

func (s *state4) ActorVersion() actorstypes.Version {
	return actorstypes.Version4
}

func (s *state4) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
