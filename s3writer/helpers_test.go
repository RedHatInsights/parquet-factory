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
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	s3mocks "github.com/RedHatInsights/insights-operator-utils/s3/mocks"
	"github.com/RedHatInsights/parquet-factory/s3writer"
)

// mockS3ClientAdapter adapts MockS3Client to implement the full S3ClientAPI interface
// It provides stub implementations for methods not in MockS3Client
type mockS3ClientAdapter struct {
	*s3mocks.MockS3Client
}

// Implement missing methods from S3ClientAPI with stubs
func (m *mockS3ClientAdapter) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return nil, fmt.Errorf("HeadObject not implemented in mock")
}

func (m *mockS3ClientAdapter) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return nil, fmt.Errorf("GetObject not implemented in mock")
}

func (m *mockS3ClientAdapter) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return nil, fmt.Errorf("PutObject not implemented in mock")
}

func (m *mockS3ClientAdapter) CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	return nil, fmt.Errorf("CreateMultipartUpload not implemented in mock")
}

func (m *mockS3ClientAdapter) UploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	return nil, fmt.Errorf("UploadPart not implemented in mock")
}

func (m *mockS3ClientAdapter) CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	return nil, fmt.Errorf("CompleteMultipartUpload not implemented in mock")
}

func (m *mockS3ClientAdapter) AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	return nil, fmt.Errorf("AbortMultipartUpload not implemented in mock")
}

// newMockS3Writer creates an S3Writer with a mock client for testing
// Note: This mock only fully supports operations from insights-operator-utils (ListObjectsV2, DeleteObjects)
// Tests requiring parquet file operations need integration tests with real S3/Minio
func newMockS3Writer(t *testing.T, mockClient *s3mocks.MockS3Client) s3writer.S3Writer {
	// Wrap the mock client in an adapter that implements the full interface
	adapter := &mockS3ClientAdapter{MockS3Client: mockClient}

	writer := s3writer.S3Writer{
		S3Client: adapter,
		Bucket:   "test_bucket",
	}

	return writer
}
