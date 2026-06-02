#!/usr/bin/env bash
set -Eeuo pipefail

APP_NAME="hology-be"
APP_USER="hology"
APP_DIR="/opt/${APP_NAME}"
APP_BINARY="${APP_DIR}/hology-be"
APP_SERVICE="/etc/systemd/system/${APP_NAME}.service"
REPO_URL="https://github.com/BangNopall/hology8-be.git"
REPO_BRANCH="${REPO_BRANCH:-main}"
GO_VERSION="${GO_VERSION:-1.25.8}"

AWS_REGION="${AWS_REGION:-us-east-1}"

RDSHOST="db-projectccs-baru.ca3o3nmpvppt.us-east-1.rds.amazonaws.com"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-DB_baru}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_PASSWORD_SSM_PARAM="${DB_PASSWORD_SSM_PARAM:-}"
RDS_CERT_DIR="/opt/rds"
RDS_CERT_PATH="${RDS_CERT_DIR}/global-bundle.pem"

EFS_ID="${EFS_ID:-}"
EFS_MOUNT_POINT="${EFS_MOUNT_POINT:-/mnt/efs}"

APP_ENV="${APP_ENV:-production}"
APP_PORT="${APP_PORT:-8080}"
API_KEY="${API_KEY:-change-me}"
REDIS_HOST="${REDIS_HOST:-127.0.0.1}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_PASS="${REDIS_PASS:-}"
JWT_SECRET_KEY="${JWT_SECRET_KEY:-change-me}"
JWT_EXP_TIME="${JWT_EXP_TIME:-24}"
JWT_ADMIN_ROLE="${JWT_ADMIN_ROLE:-admin}"
JWT_USER_ROLE="${JWT_USER_ROLE:-user}"

log() {
  printf '[%s] %s\n' "$(date --iso-8601=seconds)" "$*"
}

require_root() {
  if [ "$(id -u)" -ne 0 ]; then
    echo "This script must run as root. Use sudo or EC2 user-data." >&2
    exit 1
  fi
}

install_packages() {
  log "Installing OS packages"
  export DEBIAN_FRONTEND=noninteractive

  apt-get update
  apt-get install -y \
    awscli \
    binutils \
    build-essential \
    ca-certificates \
    clang \
    cmake \
    curl \
    git \
    jq \
    libssl-dev \
    make \
    perl \
    pkg-config \
    postgresql-client \
    rsync
}

install_go() {
  if command -v go >/dev/null 2>&1 && go version | grep -q "go${GO_VERSION}"; then
    log "Go ${GO_VERSION} already installed"
    return
  fi

  log "Installing Go ${GO_VERSION}"
  local arch
  case "$(uname -m)" in
    x86_64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) echo "Unsupported architecture for Go install: $(uname -m)" >&2; exit 1 ;;
  esac

  curl -fsSL -o /tmp/go.tgz "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz"
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tgz
  ln -sf /usr/local/go/bin/go /usr/local/bin/go
  ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
}

install_rust() {
  if [ -x /root/.cargo/bin/cargo ]; then
    log "Rust toolchain already installed"
    return
  fi

  log "Installing Rust stable toolchain"
  curl -fsSL https://sh.rustup.rs | sh -s -- -y --profile minimal
}

install_efs_utils() {
  if command -v mount.efs >/dev/null 2>&1; then
    log "amazon-efs-utils already installed"
    return
  fi

  log "Building and installing amazon-efs-utils"
  export PATH="/root/.cargo/bin:/usr/local/go/bin:${PATH}"
  rm -rf /tmp/efs-utils
  git clone --depth 1 https://github.com/aws/efs-utils /tmp/efs-utils
  cd /tmp/efs-utils
  ./build-deb.sh
  apt-get install -y ./build/amazon-efs-utils*deb
}

mount_efs() {
  if [ -z "${EFS_ID}" ]; then
    log "EFS_ID is empty; skipping EFS mount"
    return
  fi

  log "Mounting EFS ${EFS_ID} at ${EFS_MOUNT_POINT}"
  mkdir -p "${EFS_MOUNT_POINT}"

  if ! grep -q "${EFS_ID}:/ ${EFS_MOUNT_POINT} efs" /etc/fstab; then
    echo "${EFS_ID}:/ ${EFS_MOUNT_POINT} efs _netdev,tls 0 0" >> /etc/fstab
  fi

  mountpoint -q "${EFS_MOUNT_POINT}" || mount -t efs -o tls "${EFS_ID}:/" "${EFS_MOUNT_POINT}"
}

download_rds_bundle() {
  log "Downloading AWS RDS global certificate bundle"
  mkdir -p "${RDS_CERT_DIR}"
  curl -fsSL -o "${RDS_CERT_PATH}" https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem
  chmod 0644 "${RDS_CERT_PATH}"
}

load_secrets() {
  if [ -n "${DB_PASSWORD}" ]; then
    return
  fi

  if [ -n "${DB_PASSWORD_SSM_PARAM}" ]; then
    log "Loading DB password from SSM Parameter Store"
    DB_PASSWORD="$(aws ssm get-parameter \
      --name "${DB_PASSWORD_SSM_PARAM}" \
      --with-decryption \
      --region "${AWS_REGION}" \
      --query 'Parameter.Value' \
      --output text)"
  fi

  if [ -z "${DB_PASSWORD}" ]; then
    echo "DB_PASSWORD is required. Set DB_PASSWORD or DB_PASSWORD_SSM_PARAM in user-data." >&2
    exit 1
  fi
}

test_rds_connection() {
  log "Testing PostgreSQL RDS connection"
  PGPASSWORD="${DB_PASSWORD}" psql \
    "host=${RDSHOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} sslmode=verify-full sslrootcert=${RDS_CERT_PATH}" \
    -c "select 1;"
}

create_app_user() {
  if ! id "${APP_USER}" >/dev/null 2>&1; then
    useradd --system --home "${APP_DIR}" --shell /usr/sbin/nologin "${APP_USER}"
  fi
}

deploy_backend() {
  log "Deploying backend from ${REPO_URL}"
  mkdir -p "${APP_DIR}"

  if [ -d "${APP_DIR}/.git" ]; then
    git -C "${APP_DIR}" fetch origin "${REPO_BRANCH}"
    git -C "${APP_DIR}" checkout "${REPO_BRANCH}"
    git -C "${APP_DIR}" reset --hard "origin/${REPO_BRANCH}"
  else
    rm -rf "${APP_DIR:?}/"*
    git clone --branch "${REPO_BRANCH}" --depth 1 "${REPO_URL}" "${APP_DIR}"
  fi

  cat > "${APP_DIR}/.env" <<EOF
APP_ENV=${APP_ENV}
APP_PORT=${APP_PORT}
API_KEY=${API_KEY}

DB_HOST=${RDSHOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASS=${DB_PASSWORD}
DB_NAME=${DB_NAME}
DB_SSLMODE=verify-full
DB_SSLROOTCERT=${RDS_CERT_PATH}

GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID:-}
GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET:-}

GOMAIL_HOST=${GOMAIL_HOST:-}
GOMAIL_PORT=${GOMAIL_PORT:-}
GOMAIL_USERNAME=${GOMAIL_USERNAME:-}
GOMAIL_PASSWORD=${GOMAIL_PASSWORD:-}

FIREBASE_BUCKET=${FIREBASE_BUCKET:-}
FIREBASE_CREDENTIALS_PATH=${FIREBASE_CREDENTIALS_PATH:-}

REDIS_HOST=${REDIS_HOST}
REDIS_PORT=${REDIS_PORT}
REDIS_PASS=${REDIS_PASS}

JWT_SECRET_KEY=${JWT_SECRET_KEY}
JWT_EXP_TIME=${JWT_EXP_TIME}
JWT_ADMIN_ROLE=${JWT_ADMIN_ROLE}
JWT_USER_ROLE=${JWT_USER_ROLE}

AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-}
AWS_REGION=${AWS_REGION}
AWS_BUCKET=${AWS_BUCKET:-}
EOF

  cd "${APP_DIR}"
  go mod download
  CGO_ENABLED=0 GOOS=linux go build -o "${APP_BINARY}" ./cmd/app

  chown -R "${APP_USER}:${APP_USER}" "${APP_DIR}"
  chmod 0600 "${APP_DIR}/.env"
  chmod 0755 "${APP_BINARY}"
}

install_systemd_service() {
  log "Installing systemd service ${APP_NAME}"
  cat > "${APP_SERVICE}" <<EOF
[Unit]
Description=Hology 8 Backend
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${APP_USER}
Group=${APP_USER}
WorkingDirectory=${APP_DIR}
ExecStart=${APP_BINARY}
Restart=always
RestartSec=5
Environment=TZ=Asia/Jakarta

[Install]
WantedBy=multi-user.target
EOF

  systemctl daemon-reload
  systemctl enable "${APP_NAME}"
  systemctl restart "${APP_NAME}"
}

main() {
  require_root
  install_packages
  install_go
  install_rust
  install_efs_utils
  mount_efs
  download_rds_bundle
  load_secrets
  test_rds_connection
  create_app_user
  deploy_backend
  install_systemd_service
  systemctl --no-pager --full status "${APP_NAME}" || true
}

main "$@"
