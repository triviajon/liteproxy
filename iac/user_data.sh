#!/bin/bash
dnf update -y
dnf install -y podman

# Create a systemd unit file for the proxy
cat <<UNIT > /etc/systemd/system/proxy-processor.service
[Unit]
Description=Proxy Processor Container
After=network-online.target

[Service]
Restart=always
ExecStartPre=-/usr/bin/podman rm -f proxy-processor
ExecStart=/usr/bin/podman run --name proxy-processor \
-p ${PROCESSOR_PORT}:${PROCESSOR_PORT} \
-e PROXY_AUTH_TOKEN=${PROXY_AUTH_TOKEN} \
-e REDIS_HOST=${REDIS_HOST} \
${CONTAINER_IMAGE}
ExecStop=/usr/bin/podman stop proxy-processor

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now proxy-processor
