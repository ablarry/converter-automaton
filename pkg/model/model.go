package model

import (
	"fmt"
	"github.com/ablarry/converter-automaton/pkg/util"
	"log"
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
		set["S"+initialState+"\\"+v] = []string{"S", initialState, "\\", v}
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

func (p *PushDownAutomaton) Build() {
	p.CollectStates()
	p.CreateFirstRule()
	p.CreateSecondRule()
	p.CreateThirdRule()
	p.CreateFourthRule()
}

func FormatRule(i int, v []string) (string, string) {

	var keyRule, equivalence string

	if i == 1 {
		rule := strings.Join(v, "")
		keyRule = rule[0:1]
		equivalence = rule[1:len(rule)]
	} else if i == 2 {
		keyRule = strings.Join(v[0:3], "")
		equivalence = v[3]
	} else if i == 3 {
		rule := strings.Join(v, "")
		keyRule = rule[0:3]
		equivalence = rule[3:len(rule)]
	} else if i == 4 {
		rule := strings.Join(v, "")
		keyRule = rule[0:3]
		equivalence = rule[3:len(rule)]
	}

	return keyRule, equivalence
}

func (p *PushDownAutomaton) String() string {
	var rules string

	rules += fmt.Sprintf("\nRules 1:\n")
	for _, v := range p.Rules1 {
		keyRule, equivalence := FormatRule(1, v)
		rules += fmt.Sprintf("\t %s --> %s \n", keyRule, equivalence)
	}
	rules += fmt.Sprintf("Rules 2:\n")
	for _, v := range p.Rules2 {
		keyRule, equivalence := FormatRule(2, v)
		rules += fmt.Sprintf("\t %s --> %s \n", keyRule, equivalence)
	}
	rules += fmt.Sprintf("Rules 3:\n")
	for _, v := range p.Rules3 {
		keyRule, equivalence := FormatRule(3, v)
		rules += fmt.Sprintf("\t %s --> %s \n", keyRule, equivalence)
	}
	rules += fmt.Sprintf("Rules 4:\n")
	for _, v := range p.Rules4 {
		keyRule, equivalence := FormatRule(4, v)
		rules += fmt.Sprintf("\t %s --> %s \n", keyRule, equivalence)
	}
	return rules
}

func (p *PushDownAutomaton) Find(s, s2 string) bool {
	if s2 == "" {
		// Iterate for Rules2 without cost
		for _, v := range p.Rules2 {
			keyRule, _ := FormatRule(2, v)
			rule := string(s[len(s)-1]) + keyRule
			if p.Find(s+keyRule, rule) {
				log.Println("Initial Rule 2 applied " + "input: " + s + " call: " + s + keyRule)
				return true
			}
		}
		return false
	} else {
		// verify if s is equal S
		for _, v := range p.Rules1 {
			_, ruleTransform := FormatRule(1, v)
			if s == ruleTransform {
				log.Println("Rule 1 applied " + "input: " + s + " complement:" + s2)
				return true
			}
		}

		// Validation Rule 3
		for _, v := range p.Rules3 {
			keyRule, ruleTransform := FormatRule(3, v)
			if strings.Contains(s, ruleTransform) {
				s_prev := s
				s = util.Rev(strings.Replace(util.Rev(s), util.Rev(ruleTransform), util.Rev(keyRule), 1))
				s2 = strings.Replace(s2, ruleTransform, keyRule, 1)
				if p.Find(s, s2) {
					log.Println("Rule 3 applied " + "input: " + s_prev + " call: " + s + " complement:" + s2)
					return true
				}
			}
		}

		// Validation Rule 4
		for _, v := range p.Rules4 {
			keyRule, ruleTransform := FormatRule(4, v)
			if strings.Contains(s, ruleTransform) {
				s_prev := s
				s = util.Rev(strings.Replace(util.Rev(s), util.Rev(ruleTransform), util.Rev(keyRule), 1))
				//s2 = strings.Replace(s2, ruleTransform, keyRule, 1)
				s2 = keyRule
				if p.Find(s, s2) {
					log.Println("Rule 4 applied " + "input: " + s_prev + " call: " + s + " complement:" + s2)
					return true
				}
			}
		}

		// Validation Rule 2
		for _, v := range p.Rules2 {
			keyRule, _ := FormatRule(2, v)
			r := util.Rev(strings.Replace(util.Rev(s), util.Rev(s2), util.Rev(keyRule+s2), 1))
			for _, v4 := range p.Rules4 {
				_, ruleTransform := FormatRule(4, v4)
				if strings.Contains(r, ruleTransform) {
					found := p.Find(r, keyRule)
					if found {
						log.Println("Rule 2 applied " + "input: " + s + " call: " + r + " complement:" + s2)
					}
					return found
				}
			}
			for _, v := range p.Rules3 {
				keyRule, ruleTransform := FormatRule(3, v)
				if strings.Contains(r, ruleTransform) {
					r_prev := r
					r = util.Rev(strings.Replace(util.Rev(r), util.Rev(ruleTransform), util.Rev(keyRule), 1))
					r2 := keyRule + s2
					//r2 = strings.Replace(r2, ruleTransform, keyRule, 1)
					if p.Find(r, r2) {
						log.Println("Rule 3 applied " + "input: " + r_prev + " call: " + r + " complement:" + r2)
						return true
					}
				}
			}

			r = util.Rev(strings.Replace(util.Rev(s), util.Rev(s2), util.Rev(s2+keyRule), 1))
			for _, v4 := range p.Rules4 {
				_, ruleTransform := FormatRule(4, v4)
				if strings.Contains(r, ruleTransform) {
					found := p.Find(r, keyRule)
					if found {
						log.Println("Rule 2 applied " + "input: " + s + " call: " + r + " complement:" + s2)
					}
					return found
				}
			}
			for _, v := range p.Rules3 {
				keyRule, ruleTransform := FormatRule(3, v)
				if strings.Contains(r, ruleTransform) {
					r_prev := r
					r = util.Rev(strings.Replace(util.Rev(r), util.Rev(ruleTransform), util.Rev(keyRule), 1))
					r2 := s2 + keyRule
					if p.Find(r, r2) {
						log.Println("Rule 3 applied " + "input: " + r_prev + " call: " + r + " complement:" + r2)
						return true
					}
				}
			}
		}
		return false
	}
}
