package model

import (
	"github.com/ablarry/converter-automaton/pkg/util"
	"strings"
)

// Transition struct represent a transition of PushDown Automaton
type Transition struct {
	InitialState string
	ReadElement  string
	PullElement  string
	PushElement  string
	FinalState   string
}

// MetaData struct represents info of PushDown Automaton
type MetaData struct {
	AcceptStates []string
}

// PushDownAutomaton struct represent PushDown Automaton
type PushDownAutomaton struct {
	States        []string
	StackElements []string
	Transitions   []*Transition
	MetaData      *MetaData
	Rules1        map[string][]string
	Rules2        map[string][]string
	Rules3        map[string][]string
	Rules4        map[string][]string
}

// CollectStates search all states an collect in p.States and collect elements in stack without repetition
func (p *PushDownAutomaton) CollectStates() {
	setStates := make(map[string]string)
	setStackElements := make(map[string]string)
	//Get states of PushDown
	for _, v := range p.Transitions {
		setStates[v.InitialState] = v.InitialState
		setStates[v.FinalState] = v.FinalState
		setStackElements[v.PushElement] = v.PushElement
		setStackElements[v.PullElement] = v.PullElement
	}
	p.States = make([]string, 0)
	for k, _ := range setStates {
		p.States = append(p.States, k)
	}
	p.StackElements = make([]string, 0)
	for k, _ := range setStackElements {
		p.StackElements = append(p.StackElements, k)
	}
}

// CreateFirstRule For each acceptState form rule S --> <initialState,lambda,State>
func (p *PushDownAutomaton) CreateFirstRule() *map[string][]string {
	initialState := (p.Transitions)[0].InitialState
	set := make(map[string][]string)
	for _, v := range p.MetaData.AcceptStates {
		set[v] = []string{"S", initialState, "\\", v}
	}
	p.Rules1 = set
	return &set
}

// CreateSecondRule For each state p in M form the rewrite rule <p,lambda,p> --> lambda
func (p *PushDownAutomaton) CreateSecondRule() *map[string][]string {
	if len(p.States) == 0 {
		p.CollectStates()
	}
	set := make(map[string][]string)
	for _, v := range p.States {
		set[v+"\\"+v+"\\"] = []string{v, "\\", v, "\\"}
	}
	p.Rules2 = set
	return &set
}

// CreateThirdRule For each transition (p,x,y;q,z) in M (where y is not lambda)
// generate a rewrote rule <p,y,r> --> z<q,z,r> for each state r in M
func (p *PushDownAutomaton) CreateThirdRule() *map[string][]string {
	if len(p.States) == 0 {
		p.CollectStates()
	}
	set := make(map[string][]string)
	for _, v := range p.Transitions {
		if v.PullElement != "\\" {
			for _, r := range p.States {
				set[v.InitialState+v.PullElement+r] = []string{v.InitialState, v.PullElement, r, v.ReadElement, v.FinalState, v.PushElement, r}
			}
		}
	}

	p.Rules3 = set
	return &set
}

// CreateFourthRule for each transition of the form (p,x,lambda;q,z), generate all rewrite rules of
// the form <p,w,r> --> x<q,z,k><k,w,r> where w is either a stack symbol or lamda, while k and r (possibly
// equal) are states of M.
func (p *PushDownAutomaton) CreateFourthRule() *map[string][]string {
	if len(p.States) == 0 {
		p.CollectStates()
	}
	set := make(map[string][]string)
	for _, v := range p.Transitions {
		if v.PullElement == "\\" {
			for _, r := range p.States {
				for _, w := range p.StackElements {
					for _, k := range p.States {
						set[v.InitialState+w+r+v.ReadElement+v.FinalState+v.PushElement+k+k+w+r] = []string{v.InitialState, w, r, v.ReadElement, v.FinalState, v.PushElement, k, k, w, r}
					}
				}
			}
		}
	}

	p.Rules4 = set
	return &set
}

// How to escape?
func (p *PushDownAutomaton) Find(s, s2 string, level *int) bool {
	if s2 == "" {
		// iterate for rules without cost
		for k := range p.Rules2 {
			keyRule := k[0:3]
			rule := string(s[len(s)-1]) + keyRule
			i := *level
			if p.Find(s+keyRule, rule, &i) {
				return true
			}
		}
		return false
	} else {
		// verify if s is equal S

		for _, v := range p.Rules3 {
			rule := strings.Join(v, "")
			keyRule := rule[0:3]
			ruleTransform := rule[3:len(rule)]
			if strings.HasSuffix(s, ruleTransform) {
				// replace v in s
				s = util.Rev(strings.Replace(util.Rev(s), util.Rev(ruleTransform), util.Rev(keyRule), 1))
				// replace v in s2
				s2 = strings.Replace(s2, ruleTransform, keyRule, 1)
				if p.Find(s, s2, level) {
					return true
				}
			}
		}
		for _, v := range p.Rules4 {
			rule := strings.Join(v, "")
			keyRule := rule[0:3]
			ruleTransform := rule[3:len(rule)]
			if strings.HasSuffix(s, ruleTransform) {
				// replace v in s
				s = util.Rev(strings.Replace(util.Rev(s), util.Rev(ruleTransform), util.Rev(keyRule), 1))
				// replace v in s2
				//s2 = strings.Replace(s2, ruleTransform, keyRule, 1)
				s2 = keyRule
				if p.Find(s, s2, level) {
					return true
				}
			}
		}
		for k := range p.Rules2 {
			keyRule := k[0:3]
			r := util.Rev(strings.Replace(util.Rev(s), util.Rev(s2), util.Rev(keyRule+s2), 1))
			r2 := keyRule + s2
			for _, v4 := range p.Rules4 {
				rule4 := strings.Join(v4, "")
				ruleTransform := rule4[3:len(rule4)]
				if strings.HasSuffix(r, ruleTransform) {
					return p.Find(r, r2, level)
				}
			}
			r = util.Rev(strings.Replace(util.Rev(s), util.Rev(s2), util.Rev(s2+keyRule), 1))
			r2 = s2 + keyRule
			for _, v4 := range p.Rules4 {
				rule4 := strings.Join(v4, "")
				ruleTransform := rule4[3:len(rule4)]
				if strings.HasSuffix(r, ruleTransform) {
					return p.Find(r, r2, level)
				}
			}
		}
		return false
	}
}
