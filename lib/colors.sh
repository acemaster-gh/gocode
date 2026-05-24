#!/usr/bin/env bash
# =============================================================================
# lib/colors.sh — Random cohesive HSL color palette generator
# gocode v1.0.3 (carried from v1.0.2 — do not simplify)
# =============================================================================

generate_palette() {
    # Random base hue — skip red-ish range for a professional look
    local hue; hue=$(( ( RANDOM % 300 ) + 20 ))
    local hue_comp=$(( (hue + 180) % 360 ))
    local hue_analog=$(( (hue + 30) % 360 ))

    hsl_to_hex() {
        local h=$1 s=$2 l=$3
        python3 -c "
import colorsys
h,s,l = $h/360, $s/100, $l/100
r,g,b = colorsys.hls_to_rgb(h,l,s)
print('#{:02x}{:02x}{:02x}'.format(int(r*255),int(g*255),int(b*255)))
" 2>/dev/null
    }

    export PALETTE_PRIMARY=$(hsl_to_hex "$hue" 70 55)
    export PALETTE_PRIMARY_HOVER=$(hsl_to_hex "$hue" 70 45)
    export PALETTE_PRIMARY_LIGHT=$(hsl_to_hex "$hue" 70 70)
    export PALETTE_BG_MAIN=$(hsl_to_hex "$hue" 20 96)
    export PALETTE_BG_CARD=$(hsl_to_hex "$hue" 15 99)
    export PALETTE_TEXT_MAIN=$(hsl_to_hex "$hue" 20 10)
    export PALETTE_TEXT_SECONDARY=$(hsl_to_hex "$hue" 15 25)
    export PALETTE_TEXT_MUTED=$(hsl_to_hex "$hue" 10 45)
    export PALETTE_BORDER=$(hsl_to_hex "$hue" 20 82)
    export PALETTE_SUCCESS=$(hsl_to_hex 140 60 45)
    export PALETTE_ERROR=$(hsl_to_hex 0 70 55)
    export PALETTE_WARNING=$(hsl_to_hex 38 90 55)
    export PALETTE_DARK_BG_MAIN=$(hsl_to_hex "$hue" 30 8)
    export PALETTE_DARK_BG_CARD=$(hsl_to_hex "$hue" 25 13)
    export PALETTE_DARK_BG_INPUT=$(hsl_to_hex "$hue" 20 20)
    export PALETTE_DARK_TEXT_MAIN=$(hsl_to_hex "$hue" 15 96)
    export PALETTE_DARK_TEXT_SECONDARY=$(hsl_to_hex "$hue" 12 78)
    export PALETTE_DARK_BORDER=$(hsl_to_hex "$hue" 18 28)
    export PALETTE_HUE="$hue"

    log_ok "Color palette generated (hue: ${hue}°)"
}

inject_palette() {
    local css_file="${1:-style.css}"
    [ ! -f "$css_file" ] && return

    sed -i "s/--primary:.*$/--primary:           $PALETTE_PRIMARY;/"             "$css_file"
    sed -i "s/--primary-hover:.*$/--primary-hover:     $PALETTE_PRIMARY_HOVER;/" "$css_file"
    sed -i "s/--primary-light:.*$/--primary-light:     $PALETTE_PRIMARY_LIGHT;/" "$css_file"
    sed -i "s/--bg-main:.*$/--bg-main:           $PALETTE_BG_MAIN;/"             "$css_file"
    sed -i "s/--bg-card:.*$/--bg-card:           $PALETTE_BG_CARD;/"             "$css_file"
    sed -i "s/--text-main:.*$/--text-main:         $PALETTE_TEXT_MAIN;/"         "$css_file"
    sed -i "s/--text-secondary:.*$/--text-secondary:    $PALETTE_TEXT_SECONDARY;/" "$css_file"
    sed -i "s/--text-muted:.*$/--text-muted:        $PALETTE_TEXT_MUTED;/"       "$css_file"
    sed -i "s/--border:.*$/--border:            $PALETTE_BORDER;/"               "$css_file"
    sed -i "s/--success:.*$/--success:           $PALETTE_SUCCESS;/"             "$css_file"
    sed -i "s/--error:.*$/--error:             $PALETTE_ERROR;/"                 "$css_file"
    sed -i "s/--warning:.*$/--warning:           $PALETTE_WARNING;/"             "$css_file"

    sed -i "/\[data-theme=\"dark\"\]/,/^}/ {
        s/--bg-main:.*$/--bg-main:           $PALETTE_DARK_BG_MAIN;/
        s/--bg-card:.*$/--bg-card:           $PALETTE_DARK_BG_CARD;/
        s/--bg-input:.*$/--bg-input:          $PALETTE_DARK_BG_INPUT;/
        s/--text-main:.*$/--text-main:         $PALETTE_DARK_TEXT_MAIN;/
        s/--text-secondary:.*$/--text-secondary:    $PALETTE_DARK_TEXT_SECONDARY;/
        s/--border:.*$/--border:            $PALETTE_DARK_BORDER;/
    }" "$css_file"

    log_ok "Colour palette injected into style.css"
}
