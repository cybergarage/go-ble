# Copyright (C) 2021 The uecho-go Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

all: download

.PHONY: download

AD_TYPES_YAML="ad_types.yaml"
AD_TYPES_URL="https://bitbucket.org/bluetooth-SIG/public/raw/main/assigned_numbers/core/"

${AD_TYPES_YAML}:
	@wget -q -O ${AD_TYPES_YAML} ${AD_TYPES_URL}${AD_TYPES_YAML}
	@git commit ${AD_TYPES_YAML} -m "feat(database): update ${AD_TYPES_YAML}" || true

COMPANY_IDENTIFIERS_YAML="company_identifiers.yaml"
COMPANY_IDENTIFIERS_URL="https://bitbucket.org/bluetooth-SIG/public/raw/main/assigned_numbers/company_identifiers/"

${COMPANY_IDENTIFIERS_YAML}:
	@wget -q -O ${COMPANY_IDENTIFIERS_YAML} ${COMPANY_IDENTIFIERS_URL}${COMPANY_IDENTIFIERS_YAML}
	@git commit ${COMPANY_IDENTIFIERS_YAML} -m "feat(database): update ${COMPANY_IDENTIFIERS_YAML}" || true

SERVICE_UUIDS_URL="https://bitbucket.org/bluetooth-SIG/public/raw/main/assigned_numbers/uuids/"
SERVICE_UUIDS_YAML="service_uuids.yaml"
SDO_UUIDS_YAML="sdo_uuids.yaml"

${SERVICE_UUIDS_YAML}:
	@wget -q -O ${SERVICE_UUIDS_YAML} ${SERVICE_UUIDS_URL}${SERVICE_UUIDS_YAML}
	@wget -q -O ${SDO_UUIDS_YAML} ${SERVICE_UUIDS_URL}${SDO_UUIDS_YAML}
	@git commit ${SERVICE_UUIDS_YAML} -m "feat(database): update ${SERVICE_UUIDS_YAML}" || true
	@git commit ${SDO_UUIDS_YAML} -m "feat(database): update ${SDO_UUIDS_YAML}" || true

download: ${COMPANY_IDENTIFIERS_YAML} ${AD_TYPES_YAML} ${SERVICE_UUIDS_YAML}
