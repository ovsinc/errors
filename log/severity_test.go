package log

import (
	"testing"
)

func TestParseSeverityUint(t *testing.T) {
	type args struct {
		v uint32
	}
	tests := []struct {
		name    string
		args    args
		wantS   Severity
		wantErr bool
	}{
		{
			name:    "bad",
			args:    args{100},
			wantS:   SeverityUnknown,
			wantErr: true,
		},
		{
			name:    "unknown",
			args:    args{SeverityUnknown.Uint32()},
			wantS:   SeverityUnknown,
			wantErr: true,
		},
		{
			name:    "ends",
			args:    args{SeverityEnds.Uint32()},
			wantS:   SeverityUnknown,
			wantErr: true,
		},
		{
			name:    "warn",
			args:    args{SeverityWarn.Uint32()},
			wantS:   SeverityWarn,
			wantErr: false,
		},
		{
			name:    "error",
			args:    args{SeverityError.Uint32()},
			wantS:   SeverityError,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := ParseSeverityUint(tt.args.v)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSeverityUint() must have error, but got nil")
					return
				}
				if tt.wantS != gotS {
					t.Errorf("Got err ParseSeverityUint() = %v, want %v", gotS, tt.wantS)
					return
				}
			}

			if gotS != tt.wantS {
				t.Errorf("ParseSeverityUint() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestParseSeverityString(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantS   Severity
		wantErr bool
	}{
		{
			name:    "warn: warn",
			args:    args{"warn"},
			wantS:   SeverityWarn,
			wantErr: false,
		},
		{
			name:    "warn: w",
			args:    args{"w"},
			wantS:   SeverityWarn,
			wantErr: false,
		},
		{
			name:    "warn: warning",
			args:    args{"warning"},
			wantS:   SeverityWarn,
			wantErr: false,
		},
		{
			name:    "error: error",
			args:    args{"error"},
			wantS:   SeverityError,
			wantErr: false,
		},
		{
			name:    "error: e",
			args:    args{"e"},
			wantS:   SeverityError,
			wantErr: false,
		},
		{
			name:    "error: err",
			args:    args{"err"},
			wantS:   SeverityError,
			wantErr: false,
		},
		{
			name:    "unknown",
			args:    args{"sfsfssdfdf"},
			wantS:   SeverityUnknown,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := ParseSeverityString(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSeverityString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("ParseSeverityString() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
