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

package s3writer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	sourceS3 "github.com/xitongsys/parquet-go-source/s3v2"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

const (
	// The new version of the parquet-go-source library requires an ACL parameter.
	// Since we didn't use it before and relied on the default ACL setting in AWS, we are leaving
	// the value as an empty string.
	// Jira link https://issues.redhat.com/browse/CCXDEV-15240
	// Docs about ACL can be found here https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html
	ACL = ""
)

// S3File handle parquet file
type S3File struct {
	file   source.ParquetFile
	writer *writer.ParquetWriter
}

// AddRow add row to current parquet file
func (file *S3File) AddRow(row interface{}) error {
	return file.writer.Write(row)
}

// CloseFile close file
func (file *S3File) CloseFile() error {
	if err := file.writer.WriteStop(); err != nil {
		return err
	}
	err := file.file.Close()
	if err != nil {
		return err
	}

	return nil
}

// NewFile create new parquet file instance
func (s3Writer *S3Writer) NewFile(ctx context.Context, path string, schema interface{}) (S3ParquetFile, error) {
	pfw, err := newS3FileWriterWithS3Writer(ctx, s3Writer, path)
	if err != nil {
		return nil, err
	}
	pw, err := writer.NewParquetWriter(pfw, schema, 4)
	if err != nil {
		return nil, err
	}
	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	file := &S3File{
		file:   pfw,
		writer: pw,
	}

	return file, nil
}

func newS3FileWriterWithS3Writer(ctx context.Context, s3Writer *S3Writer, path string) (source.ParquetFile, error) {
	// s3v2 package signature: NewS3FileWriterWithClient(ctx, client, bucket, key, uploaderOptions, putObjectOptions)
	// The ACL is set via putObjectOptions function
	putObjectOptions := func(input *s3.PutObjectInput) {
		if ACL != "" {
			input.ACL = types.ObjectCannedACL(ACL)
		}
	}
	return sourceS3.NewS3FileWriterWithClient(
		ctx, s3Writer.S3Client, s3Writer.Bucket, path, []func(*manager.Uploader){}, putObjectOptions)
}
