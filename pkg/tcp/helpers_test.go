package tcp

import "testing"

func Test_check_isPortInValidRange(t *testing.T) {
	type fields struct {
		port int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid sample ",
			fields: fields{
				port: 8080,
			},
			want: true,
		},
		{
			name: "minimum",
			fields: fields{
				port: 1,
			},
			want: true,
		},
		{
			name: "maximum",
			fields: fields{
				port: 65535,
			},
			want: true,
		},
		{
			name: "out of range",
			fields: fields{
				port: 66000,
			},
			want: false,
		},
		{
			name: "greater than max",
			fields: fields{
				port: maxPort + 1,
			},
			want: false,
		},
		{
			name: "less than min",
			fields: fields{
				port: minPort - 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcp := &check{
				port: tt.fields.port,
			}
			if got := tcp.isPortInValidRange(); got != tt.want {
				t.Errorf("isPortInValidRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
