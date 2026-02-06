// Copyright 2021 Red Hat, Inc
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

package rulereportaggregator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	gomock "github.com/golang/mock/gomock"
	"github.com/RedHatInsights/parquet-factory/metrics"
	"github.com/RedHatInsights/parquet-factory/reportaggregators/rulereportaggregator"
	"github.com/RedHatInsights/parquet-factory/s3writer/mock"
	"github.com/RedHatInsights/parquet-factory/testdata"
)

func TestNew(t *testing.T) {
	sut := rulereportaggregator.NewRulesReportAggregator()
	assert.Equal(t, 0, len(sut.ReceivedReports), "The RulesReportAggregator is not empty!")
}

func TestHandle(t *testing.T) {
	sut := rulereportaggregator.NewRulesReportAggregator()
	err := sut.Handle(testdata.RuleHitReport)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sut.ReceivedReports))
}

func TestHandleBadData(t *testing.T) {
	sut := rulereportaggregator.NewRulesReportAggregator()
	err := sut.Handle([]byte("Hello world"))
	assert.Error(t, err)
	assert.Equal(t, 0, len(sut.ReceivedReports))
}

// TestWriteResults checks that the files are generated as expected
func TestWriteResults(t *testing.T) {
	expectedFileWritten := 2

	sut := rulereportaggregator.NewRulesReportAggregator()
	err := sut.Handle(testdata.RuleHitReport)
	assert.NoError(t, err)

	mockWriter, controller := mock.PrepareMocks(
		t,
		[]uint{1, 1},
	)
	defer controller.Finish()

	// Init metrics to avoid errors
	err = metrics.InitMetrics("testEnv")
	assert.NoError(t, err)

	actualFileWritten, err := sut.WriteResults(mockWriter)
	assert.NoError(t, err)
	assert.Equal(t, expectedFileWritten, actualFileWritten)
}

func TestWriteResultsError(t *testing.T) {
	sut := rulereportaggregator.NewRulesReportAggregator()
	err := sut.Handle(testdata.RuleHitReport)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	mockWriter := mock.NewMockS3ParquetWriter(mockCtrl)
	defer mockCtrl.Finish()
	anyMatcher := gomock.Any()

	mockWriter.EXPECT().
		NewFile(anyMatcher, anyMatcher, anyMatcher).
		Return(nil, errors.New("test new file error"))
	mockWriter.EXPECT().DeleteFiles(anyMatcher).Times(1)
	mockWriter.EXPECT().Prefix().AnyTimes()
	mockWriter.EXPECT().GetLastIndexForParquet(anyMatcher, anyMatcher).AnyTimes()

	// Init metrics to avoid errors
	err = metrics.InitMetrics("testEnv")
	assert.NoError(t, err)

	_, err = sut.WriteResults(mockWriter)
	assert.Error(t, err)
}
