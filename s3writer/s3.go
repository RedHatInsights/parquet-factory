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
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	s3utils "github.com/RedHatInsights/insights-operator-utils/s3"
)

// GetLastIndexForParquet a map with the last used index for the objects in a given filepath
func (s3Writer *S3Writer) GetLastIndexForParquet(ctx context.Context, folder string) map[string]int {
	retval := map[string]int{}

	output, err := listBucket(ctx, s3Writer, folder)

	if err != nil {
		log.Error().Err(err).Msgf("Unable to retrieve the indexes from S3 bucket")
		return retval
	}

	for _, f := range output {
		log.Debug().Msgf("Filepath: %s", f)
		tablename, index, err := getKeyAndIndex(f)
		if err != nil {
			log.Warn().Msgf("Warning: ignoring %s\n", f)
		} else {
			currentIndex, ok := retval[tablename]

			if !ok || currentIndex < index {
				retval[tablename] = index
			}
		}
	}

	return retval
}

// listBucket all the files inside the Minio bucket inside a folder
func listBucket(ctx context.Context, s3Writer *S3Writer, folder string) ([]string, error) {
	return s3utils.ListBucket(
		ctx,
		s3Writer.S3Client,
		s3Writer.Bucket,
		folder,
		"",
		listMaxKey,
	)
}

func getKeyAndIndex(filepath string) (string, int, error) {
	filename := path.Base(filepath)
	ext := path.Ext(filename)
	filename = filename[0 : len(filename)-len(ext)]
	comps := strings.Split(filename, "-")

	index, err := strconv.Atoi(comps[len(comps)-1])
	if err == nil {
		return comps[0], index, nil
	}
	return "", -1, err
}
