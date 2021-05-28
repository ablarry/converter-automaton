package model

import (
	"github.com/ablarry/converter-automaton/pkg/mapper"
	"github.com/ablarry/converter-automaton/pkg/model"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCollectStates(t *testing.T) {
	transitions, metadata, err := mapper.MapperFileToPA("./example2.pd")
	assert.Nil(t, err, "Error not excepcted")
	p := &model.PushDownAutomaton{Transitions: transitions, MetaData: metadata}
	p.CollectStates()
	reflect.DeepEqual(p.States, []string{"g", "h", "f"})
}

func TestPushDownAutomaton_CreateFirstRule(t *testing.T) {
	type fields struct {
		Transitions []*model.Transition
		MetaData    *model.MetaData
		Rules1      map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *map[string][]string
	}{
		{name: "Creation First Rule - Ok",
			fields: fields{Transitions: []*model.Transition{&model.Transition{InitialState: "f"}},
				MetaData: &model.MetaData{AcceptStates: []string{"h"}},
			},
			want: &map[string][]string{
				"h": {"S", "f", "\\", "h"},
			},
		},
		{name: "Creation First Rule two symbols of accept states - Ok",
			fields: fields{Transitions: []*model.Transition{&model.Transition{InitialState: "f"}},
				MetaData: &model.MetaData{AcceptStates: []string{"h", "i"}},
			},
			want: &map[string][]string{
				"h": {"S", "f", "\\", "h"},
				"i": {"S", "f", "\\", "i"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &model.PushDownAutomaton{
				Transitions: tt.fields.Transitions,
				MetaData:    tt.fields.MetaData,
				Rules1:      tt.fields.Rules1,
			}
			if got := p.CreateFirstRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFirstRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestPushDownAutomaton_CreateSecondRule(t *testing.T) {
	type fields struct {
		States      []string
		Transitions []*model.Transition
		MetaData    *model.MetaData
		Rules1      map[string][]string
		Rules2      map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *map[string][]string
	}{
		{name: "Creation Second Rule - Ok",
			fields: fields{Transitions: []*model.Transition{
				&model.Transition{"g", "c", "c", "\\", "h"},
				&model.Transition{"f", "c", "\\", "c", "g"},
			},
				MetaData: &model.MetaData{AcceptStates: []string{"h"}},
			},
			want: &map[string][]string{
				"f" + "\\" + "f" + "\\": {"f", "\\", "f", "\\"},
				"g" + "\\" + "g" + "\\": {"g", "\\", "g", "\\"},
				"h" + "\\" + "h" + "\\": {"h", "\\", "h", "\\"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &model.PushDownAutomaton{
				States:      tt.fields.States,
				Transitions: tt.fields.Transitions,
				MetaData:    tt.fields.MetaData,
				Rules1:      tt.fields.Rules1,
				Rules2:      tt.fields.Rules2,
			}
			if got := p.CreateSecondRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSecondRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPushDownAutomaton_CreateThirdRule(t *testing.T) {
	type fields struct {
		States      []string
		Transitions []*model.Transition
		MetaData    *model.MetaData
		Rules1      map[string][]string
		Rules2      map[string][]string
		Rules3      map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *map[string][]string
	}{
		{name: "Creation Third Rule - Ok",
			fields: fields{Transitions: []*model.Transition{
				&model.Transition{"g", "c", "c", "\\", "h"},
				&model.Transition{"f", "c", "\\", "c", "g"},
			},
				MetaData: &model.MetaData{AcceptStates: []string{"h"}},
			},
			want: &map[string][]string{
				"gcf": {"g", "c", "f", "c", "h", "\\", "f"},
				"gcg": {"g", "c", "g", "c", "h", "\\", "g"},
				"gch": {"g", "c", "h", "c", "h", "\\", "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &model.PushDownAutomaton{
				States:      tt.fields.States,
				Transitions: tt.fields.Transitions,
				MetaData:    tt.fields.MetaData,
				Rules1:      tt.fields.Rules1,
				Rules2:      tt.fields.Rules2,
				Rules3:      tt.fields.Rules3,
			}
			if got := p.CreateThirdRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateThirdRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPushDownAutomaton_CreateFourthRule(t *testing.T) {
	type fields struct {
		States        []string
		StackElements []string
		Transitions   []*model.Transition
		MetaData      *model.MetaData
		Rules1        map[string][]string
		Rules2        map[string][]string
		Rules3        map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *map[string][]string
	}{
		{name: "Creation Fourth Rule - Ok",
			fields: fields{Transitions: []*model.Transition{
				&model.Transition{"f", "c", "\\", "c", "g"},
				&model.Transition{"g", "b", "\\", "\\", "g"},
				&model.Transition{"g", "c", "c", "\\", "h"},
			},
				MetaData: &model.MetaData{AcceptStates: []string{"h"}},
			},
			want: &map[string][]string{

				"f" + "\\" + "f" + "c" + "g" + "c" + "f" + "f" + "\\" + "f": {"f", "\\", "f", "c", "g", "c", "f", "f", "\\", "f"},
				"f" + "\\" + "f" + "c" + "g" + "c" + "g" + "g" + "\\" + "f": {"f", "\\", "f", "c", "g", "c", "g", "g", "\\", "f"},
				"f" + "\\" + "f" + "c" + "g" + "c" + "h" + "h" + "\\" + "f": {"f", "\\", "f", "c", "g", "c", "h", "h", "\\", "f"},
				"f" + "\\" + "g" + "c" + "g" + "c" + "f" + "f" + "\\" + "g": {"f", "\\", "g", "c", "g", "c", "f", "f", "\\", "g"},
				"f" + "\\" + "g" + "c" + "g" + "c" + "g" + "g" + "\\" + "g": {"f", "\\", "g", "c", "g", "c", "g", "g", "\\", "g"},
				"f" + "\\" + "g" + "c" + "g" + "c" + "h" + "h" + "\\" + "g": {"f", "\\", "g", "c", "g", "c", "h", "h", "\\", "g"},
				"f" + "\\" + "h" + "c" + "g" + "c" + "f" + "f" + "\\" + "h": {"f", "\\", "h", "c", "g", "c", "f", "f", "\\", "h"},
				"f" + "\\" + "h" + "c" + "g" + "c" + "g" + "g" + "\\" + "h": {"f", "\\", "h", "c", "g", "c", "g", "g", "\\", "h"},
				"f" + "\\" + "h" + "c" + "g" + "c" + "h" + "h" + "\\" + "h": {"f", "\\", "h", "c", "g", "c", "h", "h", "\\", "h"},

				"f" + "c" + "f" + "c" + "g" + "c" + "f" + "f" + "c" + "f": {"f", "c", "f", "c", "g", "c", "f", "f", "c", "f"},
				"f" + "c" + "f" + "c" + "g" + "c" + "g" + "g" + "c" + "f": {"f", "c", "f", "c", "g", "c", "g", "g", "c", "f"},
				"f" + "c" + "f" + "c" + "g" + "c" + "h" + "h" + "c" + "f": {"f", "c", "f", "c", "g", "c", "h", "h", "c", "f"},
				"f" + "c" + "g" + "c" + "g" + "c" + "f" + "f" + "c" + "g": {"f", "c", "g", "c", "g", "c", "f", "f", "c", "g"},
				"f" + "c" + "g" + "c" + "g" + "c" + "g" + "g" + "c" + "g": {"f", "c", "g", "c", "g", "c", "g", "g", "c", "g"},
				"f" + "c" + "g" + "c" + "g" + "c" + "h" + "h" + "c" + "g": {"f", "c", "g", "c", "g", "c", "h", "h", "c", "g"},
				"f" + "c" + "h" + "c" + "g" + "c" + "f" + "f" + "c" + "h": {"f", "c", "h", "c", "g", "c", "f", "f", "c", "h"},
				"f" + "c" + "h" + "c" + "g" + "c" + "g" + "g" + "c" + "h": {"f", "c", "h", "c", "g", "c", "g", "g", "c", "h"},
				"f" + "c" + "h" + "c" + "g" + "c" + "h" + "h" + "c" + "h": {"f", "c", "h", "c", "g", "c", "h", "h", "c", "h"},

				"g" + "\\" + "f" + "b" + "g" + "\\" + "f" + "f" + "\\" + "f": {"g", "\\", "f", "b", "g", "\\", "f", "f", "\\", "f"},
				"g" + "\\" + "f" + "b" + "g" + "\\" + "g" + "g" + "\\" + "f": {"g", "\\", "f", "b", "g", "\\", "g", "g", "\\", "f"},
				"g" + "\\" + "f" + "b" + "g" + "\\" + "h" + "h" + "\\" + "f": {"g", "\\", "f", "b", "g", "\\", "h", "h", "\\", "f"},
				"g" + "\\" + "g" + "b" + "g" + "\\" + "f" + "f" + "\\" + "g": {"g", "\\", "g", "b", "g", "\\", "f", "f", "\\", "g"},
				"g" + "\\" + "g" + "b" + "g" + "\\" + "g" + "g" + "\\" + "g": {"g", "\\", "g", "b", "g", "\\", "g", "g", "\\", "g"},
				"g" + "\\" + "g" + "b" + "g" + "\\" + "h" + "h" + "\\" + "g": {"g", "\\", "g", "b", "g", "\\", "h", "h", "\\", "g"},
				"g" + "\\" + "h" + "b" + "g" + "\\" + "f" + "f" + "\\" + "h": {"g", "\\", "h", "b", "g", "\\", "f", "f", "\\", "h"},
				"g" + "\\" + "h" + "b" + "g" + "\\" + "g" + "g" + "\\" + "h": {"g", "\\", "h", "b", "g", "\\", "g", "g", "\\", "h"},
				"g" + "\\" + "h" + "b" + "g" + "\\" + "h" + "h" + "\\" + "h": {"g", "\\", "h", "b", "g", "\\", "h", "h", "\\", "h"},

				"g" + "c" + "f" + "b" + "g" + "\\" + "f" + "f" + "c" + "f": {"g", "c", "f", "b", "g", "\\", "f", "f", "c", "f"},
				"g" + "c" + "f" + "b" + "g" + "\\" + "g" + "g" + "c" + "f": {"g", "c", "f", "b", "g", "\\", "g", "g", "c", "f"},
				"g" + "c" + "f" + "b" + "g" + "\\" + "h" + "h" + "c" + "f": {"g", "c", "f", "b", "g", "\\", "h", "h", "c", "f"},
				"g" + "c" + "g" + "b" + "g" + "\\" + "f" + "f" + "c" + "g": {"g", "c", "g", "b", "g", "\\", "f", "f", "c", "g"},
				"g" + "c" + "g" + "b" + "g" + "\\" + "g" + "g" + "c" + "g": {"g", "c", "g", "b", "g", "\\", "g", "g", "c", "g"},
				"g" + "c" + "g" + "b" + "g" + "\\" + "h" + "h" + "c" + "g": {"g", "c", "g", "b", "g", "\\", "h", "h", "c", "g"},
				"g" + "c" + "h" + "b" + "g" + "\\" + "f" + "f" + "c" + "h": {"g", "c", "h", "b", "g", "\\", "f", "f", "c", "h"},
				"g" + "c" + "h" + "b" + "g" + "\\" + "g" + "g" + "c" + "h": {"g", "c", "h", "b", "g", "\\", "g", "g", "c", "h"},
				"g" + "c" + "h" + "b" + "g" + "\\" + "h" + "h" + "c" + "h": {"g", "c", "h", "b", "g", "\\", "h", "h", "c", "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &model.PushDownAutomaton{
				States:        tt.fields.States,
				StackElements: tt.fields.StackElements,
				Transitions:   tt.fields.Transitions,
				MetaData:      tt.fields.MetaData,
				Rules3:        tt.fields.Rules3,
			}
			if got := p.CreateFourthRule(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFourthRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestPushDownAutomaton_find(t *testing.T) {
	p := model.PushDownAutomaton{
		Transitions: []*model.Transition{
			&model.Transition{"f", "c", "\\", "c", "g"},
			&model.Transition{"g", "b", "\\", "\\", "g"},
			&model.Transition{"g", "c", "c", "\\", "h"},
		},
		MetaData: &model.MetaData{AcceptStates: []string{"h"}},
	}
	p.CollectStates()
	p.CreateFirstRule()
	p.CreateSecondRule()
	p.CreateThirdRule()
	p.CreateFourthRule()
	i := 10
	level := &i
	assert.True(t, p.Find("cbbc", "", level))
}
