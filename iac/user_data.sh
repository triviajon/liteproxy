#!/bin/bash
set -e

dnf update -y
dnf install -y docker
systemctl enable --now docker

cat <<UNIT > /etc/systemd/system/proxy-processor.service
[Unit]
Description=LiteProxy Processor Container
After=network-online.target docker.service
Wants=network-online.target
Requires=docker.service

[Service]
Restart=always
RestartSec=10
ExecStartPre=-/usr/bin/docker rm -f proxy-processor
ExecStart=/usr/bin/docker run --name proxy-processor \
  --network host \
  -e PROCESSOR_PORT=${PROCESSOR_PORT} \
  -e PROXY_AUTH_TOKEN=${PROXY_AUTH_TOKEN} \
  -e CACHE_SALT=${CACHE_SALT} \
  -e REDIS_HOST=${REDIS_HOST} \
  -e REDIS_PORT=${REDIS_PORT} \
  ${CONTAINER_IMAGE}
ExecStop=/usr/bin/docker stop -t 10 proxy-processor

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now proxy-processor.service