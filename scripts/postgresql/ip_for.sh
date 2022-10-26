# Copyright 2022 gab
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Usage: ip_for.sh postgresql_server
CONTAINER_ID=$(docker ps | grep ${1} | awk '{print $1}')
IP=$(docker inspect $CONTAINER_ID | python3 -c 'import json,sys;obj=json.load(sys.stdin);print(obj[0]["NetworkSettings"]["IPAddress"])')
echo $IP