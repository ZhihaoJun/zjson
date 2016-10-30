package core

import (
	"log"
)

const (
	initState        = 0
	objectStartState = iota
	objectKeyState
	objectColonState
	commaState
	objectEndState
	arrayStartState
	arrayEndState
	stringState
	numberState
	boolState
	nullState
)

var stateToStr = map[int]string{
	initState:        "initState",
	objectStartState: "objectStartState",
	objectKeyState:   "objectKeyState",
	objectColonState: "objectColonState",
	commaState:       "commaState",
	objectEndState:   "objectEndState",
	arrayStartState:  "arrayStartState",
	arrayEndState:    "arrayEndState",
	stringState:      "stringState",
	numberState:      "numberState",
	boolState:        "boolState",
	nullState:        "nullState",
}

type Parser struct {
	fsm        *FSM
	tokens     []*Token
	current    int
	jsonStack  []*JSON
	lastPop    *JSON
	tokenStack []*Token
}

func (p *Parser) next() *Token {
	if p.current >= len(p.tokens) {
		return nil
	}
	p.current += 1
	return p.tokens[p.current-1]
}

func (p *Parser) currentToken() *Token {
	return p.tokens[p.current]
}

func (p *Parser) pushJSONStack(j *JSON) {
	p.jsonStack = append(p.jsonStack, j)
}

func (p *Parser) popJSONStack() *JSON {
	p.lastPop = p.jsonStack[len(p.jsonStack)-1]
	p.jsonStack = p.jsonStack[:len(p.jsonStack)-1]
	return p.lastPop
}

func (p *Parser) pushTokenStack(t *Token) {
	p.tokenStack = append(p.tokenStack, t)
}

func (p *Parser) popTokenStack() *Token {
	t := p.tokenStack[len(p.tokenStack)-1]
	p.tokenStack = p.tokenStack[:len(p.tokenStack)-1]
	return t
}

func (p *Parser) parent() *JSON {
	if len(p.jsonStack) == 0 {
		return nil
	}
	return p.jsonStack[len(p.jsonStack)-1]
}

func (p *Parser) Parse(tokens []*Token) *JSON {
	p.tokens = tokens
	p.current = 0
	for i := 0; i < len(tokens); i++ {
		log.Println("[current state]", stateToStr[p.fsm.state])
		p.fsm.Tick()
		p.next()
	}
	return p.lastPop
}

type ExecuteFunc func(p *Parser, fsm *FSM, state int) error

func (p *Parser) generateTransitionFunc(fn ExecuteFunc) TransitionFunc {
	return func(fsm *FSM, state int) error {
		return fn(p, fsm, state)
	}
}

func (p *Parser) objectStartState(fsm *FSM, _ int) error {
	m := &JSON{
		Type: JSONTypeMap,
		Map:  map[string]*JSON{},
	}
	if p.parent() != nil {
		switch p.parent().Type {
		case JSONTypeMap:
			p.setParentMap(m)
		case JSONTypeArray:
			p.parentArrAppend(m)
		}
	}
	p.pushJSONStack(m)
	return nil
}

func (p *Parser) arrayStartStateFunc(fsm *FSM, _ int) error {
	m := &JSON{
		Type: JSONTypeArray,
		Arr:  []*JSON{},
	}
	p.pushJSONStack(m)
	return nil
}

func (p *Parser) setParentMap(j *JSON) {
	key := p.popTokenStack()
	p.parent().Map[key.ToString()] = j
}

func (p *Parser) parentArrAppend(j *JSON) {
	p.parent().Arr = append(p.parent().Arr, j)
}

func (p *Parser) objectKeyState(fsm *FSM, _ int) error {
	t := p.currentToken()
	p.pushTokenStack(t)
	return nil
}

func (p *Parser) stringState(fsm *FSM, _ int) error {
	t := p.currentToken()
	j := &JSON{
		Type: JSONTypeString,
		Str:  t.ToString(),
	}
	switch p.parent().Type {
	case JSONTypeMap:
		p.setParentMap(j)
	case JSONTypeArray:
		p.parentArrAppend(j)
	}
	return nil
}

func (p *Parser) numberState(fsm *FSM, _ int) error {
	t := p.currentToken()
	j := &JSON{
		Type:   JSONTypeNumber,
		Number: t.ToNumber(),
	}
	switch p.parent().Type {
	case JSONTypeMap:
		p.setParentMap(j)
	case JSONTypeArray:
		p.parentArrAppend(j)
	}
	return nil
}

func (p *Parser) boolState(fsm *FSM, _ int) error {
	t := p.currentToken()
	j := &JSON{
		Type: JSONTypeBool,
		Bool: t.ToBool(),
	}
	switch p.parent().Type {
	case JSONTypeMap:
		p.setParentMap(j)
	case JSONTypeArray:
		p.parentArrAppend(j)
	}
	return nil
}

func (p *Parser) nullState(fsm *FSM, _ int) error {
	j := &JSON{
		Type: JSONTypeNull,
	}
	switch p.parent().Type {
	case JSONTypeMap:
		p.setParentMap(j)
	case JSONTypeArray:
		p.parentArrAppend(j)
	}
	return nil
}

func (p *Parser) endState(fsm *FSM, _ int) error {
	p.popJSONStack()
	return nil
}

func (p *Parser) arrayStartState(fsm *FSM, _ int) error {
	m := &JSON{
		Type: JSONTypeArray,
		Arr:  []*JSON{},
	}
	if p.parent() != nil {
		switch p.parent().Type {
		case JSONTypeMap:
			p.setParentMap(m)
		case JSONTypeArray:
			p.parentArrAppend(m)
		}
	}
	p.pushJSONStack(m)
	return nil
}

func NewJSONParser() *Parser {
	p := &Parser{}
	fsm := NewFSM()

	initStateOut := p.generateTransitionFunc(func(p *Parser, fsm *FSM, state int) error {
		t := p.currentToken()
		switch t.Type {
		case TokenTypeLeftBracket:
			fsm.UpdateState(objectStartState)
		case TokenTypeLeftSquareBracket:
			fsm.UpdateState(arrayStartState)
		}
		return nil
	})
	objectStartOut := p.generateTransitionFunc(func(p *Parser, fsm *FSM, _ int) error {
		t := p.currentToken()
		switch t.Type {
		case TokenTypeString:
			fsm.UpdateState(objectKeyState)
		case TokenTypeRightBracket:
			fsm.UpdateState(objectEndState)
		}
		return nil
	})
	objectKeyOut := p.generateTransitionFunc(func(p *Parser, fsm *FSM, _ int) error {
		t := p.currentToken()
		if t.Type == TokenTypeColon {
			fsm.UpdateState(objectColonState)
		}
		return nil
	})
	valueIn := p.generateTransitionFunc(func(p *Parser, fsm *FSM, _ int) error {
		t := p.currentToken()
		switch t.Type {
		case TokenTypeString:
			fsm.UpdateState(stringState)
		case TokenTypeNumber:
			fsm.UpdateState(numberState)
		case TokenTypeBool:
			fsm.UpdateState(boolState)
		case TokenTypeNull:
			fsm.UpdateState(nullState)
		case TokenTypeLeftBracket:
			fsm.UpdateState(objectStartState)
		case TokenTypeLeftSquareBracket:
			fsm.UpdateState(arrayStartState)
		}
		return nil
	})
	valueOut := p.generateTransitionFunc(func(p *Parser, fsm *FSM, _ int) error {
		t := p.currentToken()
		switch t.Type {
		case TokenTypeRightBracket:
			fsm.UpdateState(objectEndState)
		case TokenTypeRightSquareBracket:
			fsm.UpdateState(arrayEndState)
		case TokenTypeComma:
			fsm.UpdateState(commaState)
		}
		return nil
	})
	commaOut := p.generateTransitionFunc(func(p *Parser, fsm *FSM, state int) error {
		switch p.parent().Type {
		case JSONTypeMap:
			fsm.UpdateState(objectKeyState)
		case JSONTypeArray:
			return valueIn(fsm, state)
		}
		return nil
	})

	fsm.AddTransitionFunc(initState, initStateOut)
	fsm.AddTransitionFunc(objectStartState, objectStartOut)
	fsm.AddTransitionFunc(objectKeyState, objectKeyOut)
	fsm.AddTransitionFunc(objectColonState, valueIn)
	fsm.AddTransitionFunc(stringState, valueOut)
	fsm.AddTransitionFunc(numberState, valueOut)
	fsm.AddTransitionFunc(boolState, valueOut)
	fsm.AddTransitionFunc(nullState, valueOut)
	fsm.AddTransitionFunc(objectEndState, valueOut)
	fsm.AddTransitionFunc(arrayEndState, valueOut)
	fsm.AddTransitionFunc(commaState, commaOut)
	fsm.AddTransitionFunc(arrayStartState, valueIn)

	fsm.SetStateFunc(objectStartState, p.objectStartState)
	fsm.SetStateFunc(objectKeyState, p.objectKeyState)
	fsm.SetStateFunc(stringState, p.stringState)
	fsm.SetStateFunc(numberState, p.numberState)
	fsm.SetStateFunc(boolState, p.boolState)
	fsm.SetStateFunc(nullState, p.nullState)
	fsm.SetStateFunc(objectEndState, p.endState)
	fsm.SetStateFunc(arrayEndState, p.endState)
	fsm.SetStateFunc(arrayStartState, p.arrayStartState)

	fsm.Reset()
	p.fsm = fsm
	return p
}
