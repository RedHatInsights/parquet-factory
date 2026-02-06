// Copyright 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mock

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
)

// RowsPerFileExpectation defines how many rows should be added per file
// It will help to define the expected calls to "AddRow" for each created file
type RowsPerFileExpectation []uint

// PrepareMocks prepare the expectation for the mock writer and returns it
func PrepareMocks(t *testing.T, expectedRows RowsPerFileExpectation) (*MockS3ParquetWriter, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)
	mockWriter := NewMockS3ParquetWriter(mockCtrl)
	mockFile := NewMockS3ParquetFile(mockCtrl)

	// Expect objects
	expectMockFile := mockFile.EXPECT()
	expectMockWriter := mockWriter.EXPECT()

	anyMatcher := gomock.Any()

	for _, numRows := range expectedRows {
		expectMockWriter.Prefix()
		expectMockWriter.GetLastIndexForParquet(anyMatcher, anyMatcher).Return(map[string]int{})
		expectMockWriter.Prefix()
		expectMockWriter.NewFile(anyMatcher, anyMatcher, anyMatcher).
			Return(mockFile, nil)
		for row := uint(0); row < numRows; row++ {
			expectMockFile.AddRow(anyMatcher).Return(nil)
		}
		expectMockFile.CloseFile().Return(nil)
	}

	return mockWriter, mockCtrl
}
