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

	account9 "github.com/filecoin-project/go-state-types/builtin/v9/account"
)

var _ State = (*state9)(nil)

func load9(store adt.Store, root cid.Cid) (State, error) {
	out := state9{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make9(store adt.Store, addr address.Address) (State, error) {
	out := state9{store: store}
	out.State = account9.State{Address: addr}
	return &out, nil
}

type state9 struct {
	account9.State
	store adt.Store
}

func (s *state9) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state9) GetState() interface{} {
	return &s.State
}

func (s *state9) ActorKey() string {
	return manifest.AccountKey
}

func (s *state9) ActorVersion() actorstypes.Version {
	return actorstypes.Version9
}

func (s *state9) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
