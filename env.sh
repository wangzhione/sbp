sudo sh -c 'cat >> /etc/profile <<EOF

# ===== Go global build options (merged with go env) =====
# Combine system GOFLAGS with user go env configuration
go_env_flags=\$(go env GOFLAGS 2>/dev/null)
if [ -n "\$go_env_flags" ]; then
    export GOFLAGS="\$go_env_flags -trimpath -buildvcs=true"
else
    export GOFLAGS="-trimpath -buildvcs=true"
fi

EOF'

source /etc/profile

go env GOFLAGS