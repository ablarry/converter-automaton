package mapper

import (
	"github.com/ablarry/converter-automaton/pkg/mapper"
	"github.com/ablarry/converter-automaton/pkg/model"
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{"TestReadFile - example1.csv", args{"./example1.csv"}, [][]string{{"g", "c", "c;h", "\\"}, {"f", "c", "\\;g", "c"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapper.ReadFile(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapperToTransition(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want *model.Transition
	}{
		{"TestMapperToTransition - OK", args{[]string{"g", "c", "c;h", "\\"}}, &model.Transition{InitialState: "g", ReadElement: "c", PullElement: "c", FinalState: "h", PushElement: "\\"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapper.MapperToTransition(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapperToTransition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapperToMetadata(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want *model.MetaData
	}{
		{"TestMapperToMetadata - OK", args{[]string{"METADATA", "i"}}, &model.MetaData{AcceptStates: []string{"i"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapper.MapperToMetadata(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapperToMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMapperFileToPA(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.PushDownAutomaton
		wantErr bool
	}{
		{name: "TestMapperFileToPA - example2.pd",
			args: args{"./example2.pd"},
			want: &model.PushDownAutomaton{
				Transitions: []*model.Transition{
					{
						InitialState: "g",
						ReadElement:  "c",
						PullElement:  "c",
						FinalState:   "h",
						PushElement:  "\\",
					},
					{
						InitialState: "f",
						ReadElement:  "c",
						PullElement:  "\\",
						FinalState:   "g",
						PushElement:  "c",
					},
				},
				MetaData: &model.MetaData{AcceptStates: []string{"h"}},
			},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.MapperFileToPA(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapperFileToPA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapperFileToPA() got = %v, want %v", got, tt.want)
			}
		})
	}
}
