sudo sh -c 'cat >> /etc/profile <<EOF

# ===== Go global build options =====
# Automatically apply to all go commands:
# -trimpath: remove file system paths for reproducible builds
# -buildvcs=true: record VCS info in binaries for traceability
export GOFLAGS="-trimpath -buildvcs=true"

EOF'

source /etc/profile

go env GOFLAGS