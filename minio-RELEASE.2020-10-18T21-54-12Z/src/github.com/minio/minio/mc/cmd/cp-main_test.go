/*
 * MinIO Client (C) 2019 MinIO, Inc.
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

package cmd

import (
	"reflect"
	"testing"
)

func TestParseMetaData(t *testing.T) {
	metaDataCases := []struct {
		input  string
		output map[string]string
		err    error
		status bool
	}{
		// success scenario using ; as delimiter
		{"key1=value1;key2=value2", map[string]string{"Key1": "value1", "Key2": "value2"}, nil, true},
		// success scenario using ; as delimiter
		{"key1=m1=m2,m3=m4;key2=value2", map[string]string{"Key1": "m1=m2,m3=m4", "Key2": "value2"}, nil, true},
		// success scenario using = more than once
		{"Cache-Control=max-age=90000,min-fresh=9000;key1=value1;key2=value2", map[string]string{"Cache-Control": "max-age=90000,min-fresh=9000", "Key1": "value1", "Key2": "value2"}, nil, true},
		// using different delimiter, other than '=' between key value
		{"key1:value1;key2:value2", nil, ErrInvalidMetadata, false},
		// using no delimiter
		{"key1:value1:key2:value2", nil, ErrInvalidMetadata, false},
		//success: use value in quotes
		{"Content-Disposition='form-data; name=\"description\"'", map[string]string{"Content-Disposition": "form-data; name=\"description\""}, nil, true},
		//success: use value in double quotes
		{"Content-Disposition=\"form-data; name='description'\"", map[string]string{"Content-Disposition": "form-data; name='description'"}, nil, true},
		//fail: unterminated quote
		{"Content-Disposition='form-data; name=\"description\"", nil, ErrInvalidMetadata, false},
		//fail: unterminated double quote
		{"Content-Disposition=\"form-data; name='description'", nil, ErrInvalidMetadata, false},
		//success: use value and key in quotes
		{"\"Content-Disposition\"='form-data; name=\"description\"'", map[string]string{"Content-Disposition": "form-data; name=\"description\""}, nil, true},
		//success: use value and key in quotes
		{"\"Content=Disposition;Other key part=this is also key data\"='form-data; name=\"description\"'", map[string]string{"Content=Disposition;Other key part=this is also key data": "form-data; name=\"description\""}, nil, true},
	}

	for idx, testCase := range metaDataCases {
		metaDatamap, errMeta := getMetaDataEntry(testCase.input)
		if testCase.status == true {
			if errMeta != nil {
				t.Fatalf("Test %d: generated error not matching, expected = `%s`, found = `%s`", idx+1, testCase.err, errMeta)
			}
			if !reflect.DeepEqual(metaDatamap, testCase.output) {
				t.Fatalf("Test %d: generated Map not matching, expected = `%s`, found = `%s`", idx+1, testCase.input, metaDatamap)
			}
		}

		if testCase.status == false {
			if !reflect.DeepEqual(metaDatamap, testCase.output) {
				t.Fatalf("Test %d: generated Map not matching, expected = `%s`, found = `%s`", idx+1, testCase.input, metaDatamap)
			}
			if errMeta.Cause.Error() != testCase.err.Error() {
				t.Fatalf("Test %d: generated error not matching, expected = `%s`, found = `%s`", idx+1, testCase.err, errMeta)
			}
		}
	}
}
