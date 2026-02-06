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

package utils

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// parquetFilepathWithIndexTemplate is a template for generating filepaths of parquet files
	// prefix/table_name/hourly/date=YYYY-MM-DD/hour=HH/filename-index.parquet
	parquetFilepathWithIndexTemplate = "%v/%v/hourly/date=%d-%02d-%02d/hour=%02d/%v-%d.parquet"

	// for files without index, like file_without_index.parquet
	// prefix/table_name/hourly/date=YYYY-MM-DD/hour=HH/filename.parquet
	parquetFilepathTemplate = "%v/%v/hourly/date=%d-%02d-%02d/hour=%02d/%v.parquet"

	// for passing to GetLastIndexForParquet
	hourPrefixTemplate = "%v/%v/hourly/date=%d-%02d-%02d/hour=%02d/"
)

// minimal struct to parse only the archive path into json
type archivePath struct {
	Path string `json:"path"`
}

// ArchivePathSet is a wrapper around a map[string]struct{}
// to abstract a set. Implements utils.Set interface
type ArchivePathSet struct {
	paths map[string]struct{}
	mux   sync.RWMutex
}

// Set minimal interface for a set
type Set interface {
	Add(string) bool
	Remove(string) bool
	Contains(string) bool
	Size() int
	Range()
}

// Add adds a path to the Set, returns true if the
// set is changed, false otherwise
func (aps *ArchivePathSet) Add(path string) bool {
	aps.mux.Lock()
	defer aps.mux.Unlock()
	if _, ok := aps.paths[path]; ok {
		//already in set
		return false
	}
	aps.paths[path] = struct{}{}
	return true
}

// Remove removes a path from the Set, returns true if the
// set is changed, false otherwise
func (aps *ArchivePathSet) Remove(path string) bool {
	aps.mux.Lock()
	defer aps.mux.Unlock()
	if _, ok := aps.paths[path]; ok {
		delete(aps.paths, path)
		return true
	}
	return false
}

// Contains returns true if the set contains path,
// false otherwise
func (aps *ArchivePathSet) Contains(path string) bool {
	aps.mux.Lock()
	defer aps.mux.Unlock()
	_, ok := aps.paths[path]
	return ok
}

// Size returns the size (length) of the set
func (aps *ArchivePathSet) Size() int {
	aps.mux.Lock()
	defer aps.mux.Unlock()
	return len(aps.paths)
}

// Range allows the user to loop through the set
// example:
//
//	for path := range archivePathSet.Range() {
//			fmt.Println(path)
//	}
//
// calling .Add or .Remove inside the loop is not
// recommended
func (aps *ArchivePathSet) Range() chan string {
	step := make(chan string)
	go func() {
		for path := range aps.paths {
			step <- path
		}

		close(step)
	}()
	return step
}

// NewArchivePathSet returns an initialized,
// empty ArchivePathSet
func NewArchivePathSet() *ArchivePathSet {
	pathsMap := make(map[string]struct{}, 256)
	newAps := ArchivePathSet{
		paths: pathsMap,
	}
	return &newAps
}

// GetHourOnly modifies a given Time by nullifying minutes, seconds, nanoseconds
func GetHourOnly(t time.Time) time.Time {
	// nullify m, s, ns
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// GenerateParquetFilepath generates a full filepath for the file uploaded to Ceph/S3.
// To generate filepath without index postfix, pass index = -1
func GenerateParquetFilepath(timestamp time.Time, prefix, filename string, index int) string {
	var formatted string

	if index != -1 {
		formatted = fmt.Sprintf(parquetFilepathWithIndexTemplate,
			prefix, filename, timestamp.Year(), timestamp.Month(), timestamp.Day(),
			timestamp.Hour(), filename, index,
		)
	} else {
		formatted = fmt.Sprintf(parquetFilepathTemplate,
			prefix, filename, timestamp.Year(), timestamp.Month(), timestamp.Day(),
			timestamp.Hour(), filename,
		)
	}
	return formatted
}

// GenerateHourPrefix generates the full prefix for the current timestamp (hour) without the postfix
// filename to be passed to GetLastIndexForParquet
func GenerateHourPrefix(timestamp time.Time, prefix, filename string) string {
	return fmt.Sprintf(hourPrefixTemplate,
		prefix, filename, timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(),
	)
}

// GetPathFromRawMsg given a message return its "path" key
// or an error if json.Unmarshal goes wrong
func GetPathFromRawMsg(msg []byte) (string, error) {
	var parsed archivePath
	if err := json.Unmarshal(msg, &parsed); err != nil {
		log.Info().Err(err).Msgf("Unable to parse raw message: %v", msg)
		log.Error().Err(err).Msg("Unable to parse raw message")
		return "", err
	}
	return parsed.Path, nil
}
