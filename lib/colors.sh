#!/usr/bin/env bash
# =============================================================================
# lib/colors.sh — Random cohesive color palette generator
# gocode v1.0.2
# =============================================================================

# ── Generate a random cohesive HSL palette ────────────────────────────────────
# Strategy: pick one base hue, derive all colors from it using
# fixed saturation/lightness ranges so palette always looks professional
generate_palette() {
    # Random base hue — exclude red-ish range (0-20, 340-360) for professionalism
    local hue
    hue=$(( ( RANDOM % 300 ) + 20 ))

    # Complementary and analogous hues
    local hue_comp=$(( (hue + 180) % 360 ))
    local hue_analog=$(( (hue + 30) % 360 ))

    # Convert HSL to hex
    hsl_to_hex() {
        local h=$1 s=$2 l=$3
        python3 -c "
import colorsys
h, s, l = $h/360, $s/100, $l/100
r, g, b = colorsys.hls_to_rgb(h, l, s)
print('#{:02x}{:02x}{:02x}'.format(int(r*255), int(g*255), int(b*255)))
" 2>/dev/null
    }

    # Generate palette values
    local primary
    primary=$(hsl_to_hex "$hue" 70 55)
    local primary_hover
    primary_hover=$(hsl_to_hex "$hue" 70 45)
    local primary_light
    primary_light=$(hsl_to_hex "$hue" 70 70)
    local bg_main
    bg_main=$(hsl_to_hex "$hue" 20 96)
    local bg_card
    bg_card=$(hsl_to_hex "$hue" 15 99)
    local text_main
    text_main=$(hsl_to_hex "$hue" 20 10)
    local text_secondary
    text_secondary=$(hsl_to_hex "$hue" 15 25)
    local text_muted
    text_muted=$(hsl_to_hex "$hue" 10 45)
    local border
    border=$(hsl_to_hex "$hue" 20 82)
    local success
    success=$(hsl_to_hex 140 60 45)
    local error
    error=$(hsl_to_hex 0 70 55)
    local warning
    warning=$(hsl_to_hex 38 90 55)

    # Dark mode variants — same hue, inverted lightness
    local dark_bg_main
    dark_bg_main=$(hsl_to_hex "$hue" 30 8)
    local dark_bg_card
    dark_bg_card=$(hsl_to_hex "$hue" 25 13)
    local dark_bg_input
    dark_bg_input=$(hsl_to_hex "$hue" 20 20)
    local dark_text_main
    dark_text_main=$(hsl_to_hex "$hue" 15 96)
    local dark_text_secondary
    dark_text_secondary=$(hsl_to_hex "$hue" 12 78)
    local dark_border
    dark_border=$(hsl_to_hex "$hue" 18 28)

    # Export all values
    export PALETTE_PRIMARY="$primary"
    export PALETTE_PRIMARY_HOVER="$primary_hover"
    export PALETTE_PRIMARY_LIGHT="$primary_light"
    export PALETTE_BG_MAIN="$bg_main"
    export PALETTE_BG_CARD="$bg_card"
    export PALETTE_TEXT_MAIN="$text_main"
    export PALETTE_TEXT_SECONDARY="$text_secondary"
    export PALETTE_TEXT_MUTED="$text_muted"
    export PALETTE_BORDER="$border"
    export PALETTE_SUCCESS="$success"
    export PALETTE_ERROR="$error"
    export PALETTE_WARNING="$warning"
    export PALETTE_DARK_BG_MAIN="$dark_bg_main"
    export PALETTE_DARK_BG_CARD="$dark_bg_card"
    export PALETTE_DARK_BG_INPUT="$dark_bg_input"
    export PALETTE_DARK_TEXT_MAIN="$dark_text_main"
    export PALETTE_DARK_TEXT_SECONDARY="$dark_text_secondary"
    export PALETTE_DARK_BORDER="$dark_border"
    export PALETTE_HUE="$hue"

    log_ok "Color palette generated (hue: ${hue}°)"
}

# ── Inject palette into style.css ─────────────────────────────────────────────
inject_palette() {
    local css_file="${1:-style.css}"

    [ ! -f "$css_file" ] && return

    # Replace each CSS variable with generated value
    sed -i "s/--primary:.*$/--primary:           $PALETTE_PRIMARY;/" "$css_file"
    sed -i "s/--primary-hover:.*$/--primary-hover:     $PALETTE_PRIMARY_HOVER;/" "$css_file"
    sed -i "s/--primary-light:.*$/--primary-light:     $PALETTE_PRIMARY_LIGHT;/" "$css_file"
    sed -i "s/--bg-main:.*$/--bg-main:           $PALETTE_BG_MAIN;/" "$css_file"
    sed -i "s/--bg-card:.*$/--bg-card:           $PALETTE_BG_CARD;/" "$css_file"
    sed -i "s/--text-main:.*$/--text-main:         $PALETTE_TEXT_MAIN;/" "$css_file"
    sed -i "s/--text-secondary:.*$/--text-secondary:    $PALETTE_TEXT_SECONDARY;/" "$css_file"
    sed -i "s/--text-muted:.*$/--text-muted:        $PALETTE_TEXT_MUTED;/" "$css_file"
    sed -i "s/--border:.*$/--border:            $PALETTE_BORDER;/" "$css_file"
    sed -i "s/--success:.*$/--success:           $PALETTE_SUCCESS;/" "$css_file"
    sed -i "s/--error:.*$/--error:             $PALETTE_ERROR;/" "$css_file"
    sed -i "s/--warning:.*$/--warning:           $PALETTE_WARNING;/" "$css_file"

    # Dark mode overrides
    sed -i "/\[data-theme=\"dark\"\]/,/^}/ {
        s/--bg-main:.*$/--bg-main:           $PALETTE_DARK_BG_MAIN;/
        s/--bg-card:.*$/--bg-card:           $PALETTE_DARK_BG_CARD;/
        s/--bg-input:.*$/--bg-input:          $PALETTE_DARK_BG_INPUT;/
        s/--text-main:.*$/--text-main:         $PALETTE_DARK_TEXT_MAIN;/
        s/--text-secondary:.*$/--text-secondary:    $PALETTE_DARK_TEXT_SECONDARY;/
        s/--border:.*$/--border:            $PALETTE_DARK_BORDER;/
    }" "$css_file"

    log_ok "Palette injected into style.css"
}
