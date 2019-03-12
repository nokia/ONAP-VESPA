#!/bin/sh

chown prometheus:prometheus /var/lib/ves-agent/
chmod 755 /var/lib/ves-agent/
systemctl --no-reload preset ves-agent.service