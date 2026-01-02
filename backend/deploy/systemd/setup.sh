#!/bin/bash
# Server setup script for DB Lab Backend
# Run this script as root on the target server

set -e

echo "=== DB Lab Backend Server Setup ==="

# 1. Create dedicated user for the service
echo "[1/6] Creating service user..."
if ! id "dblab" &>/dev/null; then
    useradd -r -s /sbin/nologin -d /data/db-lab -M dblab
    echo "User 'dblab' created"
else
    echo "User 'dblab' already exists"
fi

# 2. Create directory structure
echo "[2/6] Creating directory structure..."
mkdir -p /data/db-lab
chown dblab:dblab /data/db-lab
chmod 755 /data/db-lab

# 3. Setup sudo permissions for cicd user
echo "[3/6] Configuring sudo permissions for cicd user..."
cat > /etc/sudoers.d/cicd-dblab << 'EOF'
# Allow cicd user to deploy db-lab service
cicd ALL=(ALL) NOPASSWD: /bin/mv /tmp/db-lab-deploy/* /data/db-lab/*
cicd ALL=(ALL) NOPASSWD: /bin/chown dblab\:dblab /data/db-lab/*
cicd ALL=(ALL) NOPASSWD: /bin/chmod 755 /data/db-lab/server
cicd ALL=(ALL) NOPASSWD: /bin/systemctl restart db-lab
cicd ALL=(ALL) NOPASSWD: /bin/systemctl status db-lab
EOF
chmod 440 /etc/sudoers.d/cicd-dblab
echo "Sudo permissions configured"

# 4. Install systemd service
echo "[4/6] Installing systemd service..."
cat > /etc/systemd/system/db-lab.service << 'EOF'
[Unit]
Description=DB Lab Backend Server
After=network.target postgresql.service

[Service]
Type=simple
User=dblab
Group=dblab
WorkingDirectory=/data/db-lab
ExecStart=/data/db-lab/server -conf /data/db-lab/prod.yml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/data/db-lab
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

# 5. Reload systemd and enable service
echo "[5/6] Enabling service..."
systemctl daemon-reload
systemctl enable db-lab

# 6. Print status
echo "[6/6] Setup complete!"
echo ""
echo "=== Next Steps ==="
echo "1. Copy the server binary to /data/db-lab/server"
echo "2. Copy prod.yml to /data/db-lab/prod.yml"
echo "3. Set permissions: chown dblab:dblab /data/db-lab/*"
echo "4. Start service: systemctl start db-lab"
echo "5. Check status: systemctl status db-lab"
echo "6. View logs: journalctl -u db-lab -f"
echo ""
echo "=== Service Management Commands ==="
echo "Start:   systemctl start db-lab"
echo "Stop:    systemctl stop db-lab"
echo "Restart: systemctl restart db-lab"
echo "Status:  systemctl status db-lab"
echo "Logs:    journalctl -u db-lab -f"