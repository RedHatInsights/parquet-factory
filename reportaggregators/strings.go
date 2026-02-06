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

package reportaggregators

// delete after parquet file generation and upload refactor (DRY)
const (
	// StartGenerateFileStr message when starting to generate a table
	StartGenerateFileStr = "Starting to generate \"%s\" table"
	// UnableGenerateTableStr message when an error generating the table
	UnableGenerateTableStr = "Unable to generate table \"%s\" rows"
	// FileStoredStr message when a file is going to be stored
	FileStoredStr = "File will be stored in %s"
	// UnableCreateFileStr message when it is not possible to create a file
	UnableCreateFileStr = "Unable to create file"
	// UnableSaveFileStr message when it is not possible to save the file
	UnableSaveFileStr = "Unable to save \"%s\" table"
	// UnableCloseFileStr message when it is not possible to close the written file
	UnableCloseFileStr = "Unable to close file"
	// UnableDeleteFileStr message when it is not possible to delete a file
	UnableDeleteFileStr = "Unable to delete file"
	// GenerateFileSuccess message when a new table is generated
	GenerateFileSuccess = "\"%s-%d\" table was generated"
	// UnableSaveRowStr message when a row can not be written
	UnableSaveRowStr = "Unable to save row \"%s\""
)
