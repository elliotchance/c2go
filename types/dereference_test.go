package types

import "testing"
import "fmt"

func TestGetDereferenceType(t *testing.T) {
	type args struct {
		cType string
	}
	tests := []struct {
		args    args
		want    string
		wantErr bool
	}{
		{args{"char [8]"}, "char", false},
		{args{"char**"}, "char*", false},
		{args{"char *[8]"}, "char *", false},
		{args{"char *[8][7]"}, "char *[7]", false},
		{args{"char *[8][7][6]"}, "char *[7][6]", false},
		{args{"char **[8]"}, "char **", false},
		{args{"char ***[8]"}, "char ***", false},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			got, err := GetDereferenceType(tt.args.cType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDereferenceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDereferenceType() = %v, want %v", got, tt.want)
			}
		})
	}
}
