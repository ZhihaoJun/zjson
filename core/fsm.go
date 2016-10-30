package core

import (
	"log"
)

type TransitionFunc func(*FSM, int) error
type StateFunc func(*FSM, int) error

type FSM struct {
	state   int
	trans   map[int][]TransitionFunc
	stateFn map[int]StateFunc
}

func NewFSM() *FSM {
	return &FSM{
		trans:   map[int][]TransitionFunc{},
		stateFn: map[int]StateFunc{},
	}
}

func (fsm *FSM) AddTransitionFunc(state int, fn TransitionFunc) *FSM {
	if fsm.trans[state] == nil {
		fsm.trans[state] = []TransitionFunc{}
	}
	fsm.trans[state] = append(fsm.trans[state], fn)
	return fsm
}

func (fsm *FSM) SetStateFunc(state int, fn StateFunc) *FSM {
	fsm.stateFn[state] = fn
	return fsm
}

func (fsm *FSM) Tick() {
	funcs := fsm.trans[fsm.state]
	for _, fn := range funcs {
		err := fn(fsm, fsm.state)
		if err != nil {
			log.Println("[error occured during transition]", err)
			continue
		}
	}
}

func (fsm *FSM) UpdateState(state int) {
	fsm.state = state
	if fsm.stateFn[fsm.state] != nil {
		err := fsm.stateFn[fsm.state](fsm, fsm.state)
		if err != nil {
			log.Println("[error occured on state]", err)
		}
	}
}

func (fsm *FSM) Reset() {
	fsm.state = 0
}
