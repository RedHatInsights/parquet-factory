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

package s3writer_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
)

var (
	mockFileContent = []byte("some content")
)

func TestPrefix(t *testing.T) {
	mockClient := s3mocks.MockS3Client{}
	sut := newMockS3Writer(t, &mockClient)
	assert.Equal(t, "", sut.Prefix())
}

func TestDeleteFiles(t *testing.T) {
	mockClient := s3mocks.MockS3Client{}
	sut := newMockS3Writer(t, &mockClient)

	t.Run("an empty input shouldn't return an error", func(t *testing.T) {
		err := sut.DeleteFiles([]string{})

		assert.NoError(t, err, "DeleteFiles shouldn't return an error if the list of files to delete is empty")
	})

	t.Run("return an error when there is a problem with S3", func(t *testing.T) {
		mockClient.Err = errors.New("an error")

		err := sut.DeleteFiles([]string{"file_1"})

		assert.Error(t, err, "DeleteFiles should return an error if something happens with S3")
	})

	t.Run("check that delete one file out of two works fine", func(t *testing.T) {
		mockClient.Contents = s3mocks.MockContents{
			"file_1": mockFileContent,
			"file_2": mockFileContent,
		}
		mockClient.Err = nil

		fmt.Println("pre:", mockClient.Contents)
		err := sut.DeleteFiles([]string{"file_1"})
		fmt.Println("post:", mockClient.Contents)

		assert.NoError(t, err, "DeleteFiles shouldn't return an error")

		assert.NotContains(t, mockClient.Contents, "file_1")
		assert.Contains(t, mockClient.Contents, "file_2")
	})
}

func TestGetLastIndexForParquet(t *testing.T) {
	mockClient := s3mocks.MockS3Client{}
	sut := newMockS3Writer(t, &mockClient)

	t.Run("an empty content should return an empty map", func(t *testing.T) {
		res := sut.GetLastIndexForParquet(context.TODO(), "test/path")
		assert.Equal(t, map[string]int{}, res)
	})

	t.Run("return an empty map if there is an error listing the files", func(t *testing.T) {
		mockClient.Err = errors.New("an error")
		mockClient.Contents = s3mocks.MockContents{
			"test/path/cluster_info-0.parquet": mockFileContent,
			"test/path/cluster_info-1.parquet": mockFileContent}
		res := sut.GetLastIndexForParquet(context.TODO(), "test/path")
		assert.Equal(t, map[string]int{}, res)
	})

	t.Run("an empty content should return an empty map", func(t *testing.T) {
		mockClient.Contents = s3mocks.MockContents{
			"test/path/cluster_info-0.parquet": mockFileContent,
			"test/path/cluster_info-1.parquet": mockFileContent}
		mockClient.Err = nil

		res := sut.GetLastIndexForParquet(context.TODO(), "test/path")
		assert.Equal(t, 1, res["cluster_info"])
	})

	t.Run("files with invalid format should be ignored", func(t *testing.T) {
		mockClient.Contents = s3mocks.MockContents{
			"test/path/cluster_info-0.parquet":             mockFileContent,
			"test/path/cluster_info-1.parquet":             mockFileContent,
			"test/path/cluster_info-invalid_index.parquet": mockFileContent}
		mockClient.Err = nil

		res := sut.GetLastIndexForParquet(context.TODO(), "test/path")
		assert.Equal(t, 1, res["cluster_info"])
	})
}
