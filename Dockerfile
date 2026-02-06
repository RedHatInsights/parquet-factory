# Copyright 2021, 2025 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

USER 0

# Print Go info for debugging.
RUN go version
RUN go env

RUN mkdir /app
COPY . /app/
WORKDIR /app

# Build binaries.
RUN make build

# Slim image for running app.
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
COPY --from=builder /app/parquet-factory /app/config.toml /app/

# copy the certificates from builder image
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /etc/pki /etc/pki

USER 1001
WORKDIR /app

CMD ["./parquet-factory"]
