#!/usr/bin/env bash
# =============================================================================
# boilerplates/react/scaffold.sh — React + Vite project scaffolder
# gocode v1.0.2
# =============================================================================

scaffold_react() {
    local project_name="$1"
    local project_path="$2"

    if ! command -v node &>/dev/null; then
        log_error "Node.js not found. Install Node first."
        return 1
    fi

    log_step "Scaffolding React + Vite project..."
    cd "$(dirname "$project_path")" || return 1

    spinner_start "Running create-vite..."
    npm create vite@latest "$project_name" -- --template react --yes &>/dev/null
    spinner_stop

    if [ ! -d "$project_path" ]; then
        log_error "Vite scaffold failed."
        return 1
    fi

    cd "$project_path" || return 1

    spinner_start "Installing dependencies..."
    npm install &>/dev/null
    spinner_stop

    log_ok "React + Vite project ready."

    # Create netlify.toml for correct build config
    cat > netlify.toml <<EOF
[build]
  command = "npm run build"
  publish = "dist"

[dev]
  command = "npm run dev"
  port = 5173
EOF

    log_ok "netlify.toml created — Netlify build configured."
}
