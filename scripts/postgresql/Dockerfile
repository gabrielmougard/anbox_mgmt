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

FROM xela7/cocoamap-base

ENV LC_ALL en_US.UTF-8 
ENV	DEBIAN_FRONTEND noninteractive

RUN wget -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add - \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" >> /etc/apt/sources.list

RUN apt-get update && apt-get install -y \
	postgresql-9.3 \
	postgresql-contrib-9.3 \
	postgresql-9.3-postgis \
	libpq-dev \
	ssl-cert

#Update the servers default ssl certs (not sure if this is really required)
RUN make-ssl-cert generate-default-snakeoil --force-overwrite
RUN cp /etc/ssl/certs/ssl-cert-snakeoil.pem  /home/docker/ssl-cert-snakeoil.pem \ 
 && cp /etc/ssl/private/ssl-cert-snakeoil.key /home/docker/ssl-cert-snakeoil.key \
 && chown postgres.postgres /home/docker/ssl-cert-snakeoil.pem \
 && chown postgres.postgres /home/docker/ssl-cert-snakeoil.key

ADD run.server.sh /home/docker/run.server.sh
ADD run.sql.sh /home/docker/run.sql.sh
RUN chmod +x /home/docker/run.*

ADD postgresql.conf /etc/postgresql/9.3/main/postgresql.conf
ADD pg_hba.conf /etc/postgresql/9.3/main/pg_hba.conf

CMD ["/home/docker/run.sh"]