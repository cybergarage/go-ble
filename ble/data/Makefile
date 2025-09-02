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

COMPANY_IDENTIFIERS_YAML="company_identifiers.yaml"
COMPANY_IDENTIFIERS_URL="https://bitbucket.org/bluetooth-SIG/public/raw/main/assigned_numbers/company_identifiers/"

${COMPANY_IDENTIFIERS_YAML}:
	@wget -q -O ${COMPANY_IDENTIFIERS_YAML} ${COMPANY_IDENTIFIERS_URL}${COMPANY_IDENTIFIERS_YAML} 

download: ${COMPANY_IDENTIFIERS_YAML} ${AD_TYPES_YAML}
	@echo "Downloaded"