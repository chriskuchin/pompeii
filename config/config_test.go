package config

import (
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "valid1",
			args: args{
				path: "testdata/valid1.yml",
			},
			want: &Config{
				Services: map[string]*ServiceConfig{
					"service1": {
						ListenerARN: "listener-rule-arn-service1",
						Canary: &PoolConfig{
							TargetGroupARN: "tg-arn-canary-service1",
							Service:        "service1-canary",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "tg-arn-primary-service1",
							Service:        "service1",
						},
					},
				},
				Workflows: WorkflowConfig{
					"validate": []*Action{
						{
							Type:   ValidatePool,
							Target: "task",
						},
						{
							Type:   ValidatePool,
							Target: "prompt",
						},
						{
							Type:   TrafficShift,
							Target: "primary",
							Ratio:  100,
						},
						{
							Type:   UpdatePool,
							Target: "canary",
							Count:  1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigFromFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetRegion(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "root_region",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {},
					"service2": {},
				},
			},
			args: args{
				service: "service1",
			},
			want: "root_region",
		},
		{
			name: "service1_region",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region: "service1_region",
					},
					"service2": {},
				},
			},
			args: args{
				service: "service1",
			},
			want: "service1_region",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.GetRegion(tt.args.service); got != tt.want {
				t.Errorf("Config.GetRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetListenerRuleARN(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "service1_listener_arn",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
					},
				},
			},
			args: args{
				service: "service1",
			},
			want: "service1-listener-rule-arn",
		},
		{
			name: "service1_listener_arn",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
					},
				},
			},
			args: args{
				service: "service2",
			},
			want: "service2-listener-rule-arn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.GetListenerRuleARN(tt.args.service); got != tt.want {
				t.Errorf("Config.GetListenerRuleARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetCanaryTargetGroupARN(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "service1_canary_tg_arn",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service1",
			},
			want: "service1-canary-tg-arn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.GetCanaryTargetGroupARN(tt.args.service); got != tt.want {
				t.Errorf("Config.GetCanaryTargetGroupARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetPrimaryTargetGroupARN(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "service1_primary_tg_arn",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service1-primary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service2-primary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service1",
			},
			want: "service1-primary-tg-arn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.GetPrimaryTargetGroupARN(tt.args.service); got != tt.want {
				t.Errorf("Config.GetPrimaryTargetGroupARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsCanaryTargetGroup(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
		arn     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "service1_canary_true",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service1-primary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service2-primary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service1",
				arn:     "service1-canary-tg-arn",
			},
			want: true,
		},
		{
			name: "service2_canary_false",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service1-primary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service2-primary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service2",
				arn:     "service1-canary-tg-arn",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.IsCanaryTargetGroup(tt.args.service, tt.args.arn); got != tt.want {
				t.Errorf("Config.IsCanaryTargetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsPrimaryTargetGroup(t *testing.T) {
	type fields struct {
		ClusterARN string
		Region     string
		Services   map[string]*ServiceConfig
	}
	type args struct {
		service string
		arn     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "service1_primary_true",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service1-primary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service2-primary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service1",
				arn:     "service1-primary-tg-arn",
			},
			want: true,
		},
		{
			name: "service2_primary_false",
			fields: fields{
				Region: "root_region",
				Services: map[string]*ServiceConfig{
					"service1": {
						Region:      "service1_region",
						ListenerARN: "service1-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service1-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service1-primary-tg-arn",
						},
					},
					"service2": {
						ListenerARN: "service2-listener-rule-arn",
						Canary: &PoolConfig{
							TargetGroupARN: "service2-canary-tg-arn",
						},
						Primary: &PoolConfig{
							TargetGroupARN: "service2-primary-tg-arn",
						},
					},
				},
			},
			args: args{
				service: "service2",
				arn:     "service1-primary-tg-arn",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClusterARN: tt.fields.ClusterARN,
				Region:     tt.fields.Region,
				Services:   tt.fields.Services,
			}
			if got := c.IsPrimaryTargetGroup(tt.args.service, tt.args.arn); got != tt.want {
				t.Errorf("Config.IsPrimaryTargetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
