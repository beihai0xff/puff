package types

import "testing"

func Test_peer_ValidateIP(t *testing.T) {
	type input struct {
		ip   string
		port int
	}
	tests := []struct {
		input   input
		wantErr bool
	}{
		{
			input:   input{"12.12.12.12", 3306},
			wantErr: false,
		},
		{
			input:   input{"12.12.12.", 3306},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		p := &peer{
			ip:   tt.input.ip,
			port: tt.input.port,
		}
		if err := p.ValidateIP(); (err != nil) != tt.wantErr {
			t.Errorf("ValidateIP() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
