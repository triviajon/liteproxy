#!/bin/bash
set -e

dnf update -y
dnf install -y podman

cat <<UNIT > /etc/systemd/system/proxy-processor.service
[Unit]
Description=LiteProxy Processor Container
After=network-online.target
Wants=network-online.target

[Service]
Restart=always
RestartSec=10
ExecStartPre=-/usr/bin/podman rm -f proxy-processor
ExecStart=/usr/bin/podman run --name proxy-processor \
  --net=host \
  -e PROCESSOR_PORT=${PROCESSOR_PORT} \
  -e PROXY_AUTH_TOKEN=${PROXY_AUTH_TOKEN} \
  -e CACHE_SALT=${CACHE_SALT} \
  -e REDIS_HOST=${REDIS_HOST} \
  -e REDIS_PORT=${REDIS_PORT} \
  ${CONTAINER_IMAGE}
ExecStop=/usr/bin/podman stop -t 10 proxy-processor

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now proxy-processor.service