package clusterCache

import (
	"testing"

	"github.com/i5heu/bonito-cache/internal/config"
	"github.com/i5heu/bonito-cache/internal/log"
)

func TestCluster_testConfig(t *testing.T) {
	type fields struct {
		Conf config.Config
		Log  log.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Conf: config.Config{
					ClusterActive: true,
					ClusterSeed:   "https://blaaablaa.de",
					ClusterKey:    "paisw089wbefpieb08gwabepifbw0e9fbwpeifbbpaisw089wbefpieb08gwabepifbw0e9fbwpeifbb",
				},
			},
			wantErr: false,
		},
		{
			name: "url invalid",
			fields: fields{
				Conf: config.Config{
					ClusterActive: true,
					ClusterSeed:   "http://blaaablaa.de",
					ClusterKey:    "paisw089wbefpieb08gwabepifbw0e9fbwpeifbb",
				},
			},
			wantErr: true,
		},
		{
			name: "key invalid",
			fields: fields{
				Conf: config.Config{
					ClusterActive: true,
					ClusterSeed:   "https://blaaablaa.de",
					ClusterKey:    "paisw08",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Conf: tt.fields.Conf,
				Log:  tt.fields.Log,
			}
			if err := c.validateConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Cluster.testConfig() name= '%v' error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
