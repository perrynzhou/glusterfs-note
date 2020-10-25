/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2020 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package notification

import (
	"encoding/xml"
	"testing"
)

func TestEqualEventTypeList(t *testing.T) {
	type args struct {
		a []EventType
		b []EventType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same order",
			args: args{
				a: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				b: []EventType{ObjectCreatedAll, ObjectAccessedAll},
			},
			want: true,
		},
		{
			name: "different order",
			args: args{
				a: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				b: []EventType{ObjectAccessedAll, ObjectCreatedAll},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				a: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				b: []EventType{ObjectRemovedAll, ObjectAccessedAll},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualEventTypeList(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("EqualEventTypeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualFilterRuleList(t *testing.T) {
	type args struct {
		a []FilterRule
		b []FilterRule
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same order",
			args: args{
				a: []FilterRule{{Name: "prefix", Value: "prefix1"}, {Name: "suffix", Value: "suffix1"}},
				b: []FilterRule{{Name: "prefix", Value: "prefix1"}, {Name: "suffix", Value: "suffix1"}},
			},
			want: true,
		},
		{
			name: "different order",
			args: args{
				a: []FilterRule{{Name: "prefix", Value: "prefix1"}, {Name: "suffix", Value: "suffix1"}},
				b: []FilterRule{{Name: "suffix", Value: "suffix1"}, {Name: "prefix", Value: "prefix1"}},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				a: []FilterRule{{Name: "prefix", Value: "prefix1"}, {Name: "suffix", Value: "suffix1"}},
				b: []FilterRule{{Name: "prefix", Value: "prefix2"}, {Name: "suffix", Value: "suffix1"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualFilterRuleList(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("EqualFilterRuleList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Equal(t *testing.T) {
	type fields struct {
		ID     string
		Arn    Arn
		Events []EventType
		Filter *Filter
	}
	type args struct {
		events []EventType
		prefix string
		suffix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "same order",
			fields: fields{
				Arn:    NewArn("minio", "sqs", "", "1", "postgresql"),
				Events: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				Filter: &Filter{
					S3Key: S3Key{
						FilterRules: []FilterRule{{Name: "prefix", Value: "prefix1"}, {Name: "suffix", Value: "suffix1"}},
					},
				},
			},
			args: args{
				events: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				prefix: "prefix1",
				suffix: "suffix1",
			},
			want: true,
		},
		{
			name: "different order",
			fields: fields{
				Arn:    NewArn("minio", "sqs", "", "1", "postgresql"),
				Events: []EventType{ObjectAccessedAll, ObjectCreatedAll},
				Filter: &Filter{
					S3Key: S3Key{
						FilterRules: []FilterRule{{Name: "suffix", Value: "suffix1"}, {Name: "prefix", Value: "prefix1"}},
					},
				},
			},
			args: args{
				events: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				prefix: "prefix1",
				suffix: "suffix1",
			},
			want: true,
		},
		{
			name: "not equal",
			fields: fields{
				Arn:    NewArn("minio", "sqs", "", "1", "postgresql"),
				Events: []EventType{ObjectAccessedAll},
				Filter: &Filter{
					S3Key: S3Key{
						FilterRules: []FilterRule{{Name: "suffix", Value: "suffix1"}, {Name: "prefix", Value: "prefix1"}},
					},
				},
			},
			args: args{
				events: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				prefix: "prefix1",
				suffix: "suffix1",
			},
			want: false,
		},
		{
			name: "different arn",
			fields: fields{
				Events: []EventType{ObjectAccessedAll},
				Filter: &Filter{
					S3Key: S3Key{
						FilterRules: []FilterRule{{Name: "suffix", Value: "suffix1"}, {Name: "prefix", Value: "prefix1"}},
					},
				},
			},
			args: args{
				events: []EventType{ObjectCreatedAll, ObjectAccessedAll},
				prefix: "prefix1",
				suffix: "suffix1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &Config{
				ID:     tt.fields.ID,
				Arn:    tt.fields.Arn,
				Events: tt.fields.Events,
				Filter: tt.fields.Filter,
			}
			if got := nc.Equal(tt.args.events, tt.args.prefix, tt.args.suffix); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_RemoveQueueByArnEventsPrefixSuffix(t *testing.T) {
	type fields struct {
		XMLName       xml.Name
		LambdaConfigs []LambdaConfig
		TopicConfigs  []TopicConfig
		QueueConfigs  []QueueConfig
	}
	type args struct {
		arn    Arn
		events []EventType
		prefix string
		suffix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Queue Configuration Removed with events, prefix",
			fields: fields{
				XMLName:       xml.Name{},
				LambdaConfigs: nil,
				TopicConfigs:  nil,
				QueueConfigs: []QueueConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sqs",
								Region:    "",
								AccountID: "1",
								Resource:  "postgresql",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Queue: "arn:minio:sqs::1:postgresql",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sqs",
					Region:    "",
					AccountID: "1",
					Resource:  "postgresql",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "",
			},
			wantErr: false,
		},
		{
			name: "Queue Configuration Removed with events, prefix, suffix",
			fields: fields{
				XMLName:       xml.Name{},
				LambdaConfigs: nil,
				TopicConfigs:  nil,
				QueueConfigs: []QueueConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sqs",
								Region:    "",
								AccountID: "1",
								Resource:  "postgresql",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
										{
											Name:  "suffix",
											Value: "y",
										},
									},
								},
							},
						},
						Queue: "arn:minio:sqs::1:postgresql",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sqs",
					Region:    "",
					AccountID: "1",
					Resource:  "postgresql",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "y",
			},
			wantErr: false,
		},
		{
			name: "Error Returned Queue Configuration Not Removed",
			fields: fields{
				XMLName:       xml.Name{},
				LambdaConfigs: nil,
				TopicConfigs:  nil,
				QueueConfigs: []QueueConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sqs",
								Region:    "",
								AccountID: "1",
								Resource:  "postgresql",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Queue: "arn:minio:sqs::1:postgresql",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sqs",
					Region:    "",
					AccountID: "1",
					Resource:  "postgresql",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "",
				suffix: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Configuration{
				XMLName:       tt.fields.XMLName,
				LambdaConfigs: tt.fields.LambdaConfigs,
				TopicConfigs:  tt.fields.TopicConfigs,
				QueueConfigs:  tt.fields.QueueConfigs,
			}
			if err := b.RemoveQueueByArnEventsPrefixSuffix(tt.args.arn, tt.args.events, tt.args.prefix, tt.args.suffix); (err != nil) != tt.wantErr {
				t.Errorf("RemoveQueueByArnEventsPrefixSuffix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfiguration_RemoveLambdaByArnEventsPrefixSuffix(t *testing.T) {
	type fields struct {
		XMLName       xml.Name
		LambdaConfigs []LambdaConfig
		TopicConfigs  []TopicConfig
		QueueConfigs  []QueueConfig
	}
	type args struct {
		arn    Arn
		events []EventType
		prefix string
		suffix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Lambda Configuration Removed with events, prefix",
			fields: fields{
				XMLName:      xml.Name{},
				QueueConfigs: nil,
				TopicConfigs: nil,
				LambdaConfigs: []LambdaConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "lambda",
								Region:    "",
								AccountID: "1",
								Resource:  "provider",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Lambda: "arn:minio:lambda::1:provider",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "lambda",
					Region:    "",
					AccountID: "1",
					Resource:  "provider",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "",
			},
			wantErr: false,
		},
		{
			name: "Lambda Configuration Removed with events, prefix, suffix",
			fields: fields{
				XMLName:      xml.Name{},
				QueueConfigs: nil,
				TopicConfigs: nil,
				LambdaConfigs: []LambdaConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "lambda",
								Region:    "",
								AccountID: "1",
								Resource:  "provider",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
										{
											Name:  "suffix",
											Value: "y",
										},
									},
								},
							},
						},
						Lambda: "arn:minio:lambda::1:provider",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "lambda",
					Region:    "",
					AccountID: "1",
					Resource:  "provider",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "y",
			},
			wantErr: false,
		},
		{
			name: "Error Returned Lambda Configuration Not Removed",
			fields: fields{
				XMLName:      xml.Name{},
				QueueConfigs: nil,
				TopicConfigs: nil,
				LambdaConfigs: []LambdaConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "lambda",
								Region:    "",
								AccountID: "1",
								Resource:  "provider",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Lambda: "arn:minio:lambda::1:provider",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "lambda",
					Region:    "",
					AccountID: "1",
					Resource:  "provider",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "",
				suffix: "",
			},
			wantErr: true,
		},
		{
			name: "Error Returned Invalid ARN",
			fields: fields{
				XMLName:      xml.Name{},
				QueueConfigs: nil,
				TopicConfigs: nil,
				LambdaConfigs: []LambdaConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "lambda",
								Region:    "",
								AccountID: "1",
								Resource:  "provider",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Lambda: "arn:minio:lambda::1:provider",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",

					Service:   "lambda",
					Region:    "",
					AccountID: "2",
					Resource:  "provider",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "",
				suffix: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Configuration{
				XMLName:       tt.fields.XMLName,
				LambdaConfigs: tt.fields.LambdaConfigs,
				TopicConfigs:  tt.fields.TopicConfigs,
				QueueConfigs:  tt.fields.QueueConfigs,
			}
			if err := b.RemoveLambdaByArnEventsPrefixSuffix(tt.args.arn, tt.args.events, tt.args.prefix, tt.args.suffix); (err != nil) != tt.wantErr {
				t.Errorf("RemoveLambdaByArnEventsPrefixSuffix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfiguration_RemoveTopicByArnEventsPrefixSuffix(t *testing.T) {
	type fields struct {
		XMLName       xml.Name
		LambdaConfigs []LambdaConfig
		TopicConfigs  []TopicConfig
		QueueConfigs  []QueueConfig
	}
	type args struct {
		arn    Arn
		events []EventType
		prefix string
		suffix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Topic Configuration Removed with events, prefix",
			fields: fields{
				XMLName:       xml.Name{},
				QueueConfigs:  nil,
				LambdaConfigs: nil,
				TopicConfigs: []TopicConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sns",
								Region:    "",
								AccountID: "1",
								Resource:  "kafka",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Topic: "arn:minio:sns::1:kafka",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sns",
					Region:    "",
					AccountID: "1",
					Resource:  "kafka",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "",
			},
			wantErr: false,
		},
		{
			name: "Topic Configuration Removed with events, prefix, suffix",
			fields: fields{
				XMLName:       xml.Name{},
				QueueConfigs:  nil,
				LambdaConfigs: nil,
				TopicConfigs: []TopicConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sns",
								Region:    "",
								AccountID: "1",
								Resource:  "kafka",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
										{
											Name:  "suffix",
											Value: "y",
										},
									},
								},
							},
						},
						Topic: "arn:minio:sns::1:kafka",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sns",
					Region:    "",
					AccountID: "1",
					Resource:  "kafka",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "x",
				suffix: "y",
			},
			wantErr: false,
		},
		{
			name: "Error Returned Topic Configuration Not Removed",
			fields: fields{
				XMLName:       xml.Name{},
				QueueConfigs:  nil,
				LambdaConfigs: nil,
				TopicConfigs: []TopicConfig{
					{
						Config: Config{
							ID: "",
							Arn: Arn{
								Partition: "minio",
								Service:   "sns",
								Region:    "",
								AccountID: "1",
								Resource:  "kafka",
							},
							Events: []EventType{
								ObjectAccessedAll,
							},
							Filter: &Filter{
								S3Key: S3Key{
									FilterRules: []FilterRule{
										{
											Name:  "prefix",
											Value: "x",
										},
									},
								},
							},
						},
						Topic: "arn:minio:sns::1:kafka",
					},
				},
			},
			args: args{
				arn: Arn{
					Partition: "minio",
					Service:   "sns",
					Region:    "",
					AccountID: "1",
					Resource:  "kafka",
				},
				events: []EventType{
					ObjectAccessedAll,
				},
				prefix: "",
				suffix: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Configuration{
				XMLName:       tt.fields.XMLName,
				LambdaConfigs: tt.fields.LambdaConfigs,
				TopicConfigs:  tt.fields.TopicConfigs,
				QueueConfigs:  tt.fields.QueueConfigs,
			}
			if err := b.RemoveTopicByArnEventsPrefixSuffix(tt.args.arn, tt.args.events, tt.args.prefix, tt.args.suffix); (err != nil) != tt.wantErr {
				t.Errorf("RemoveTopicByArnEventsPrefixSuffix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
