Output demo/demo.gif

Require vsh

Set Shell "bash"
Set FontSize 18
Set Width 1200
Set Height 550
Set Theme "Dracula"
Set TypingSpeed 40ms
Set Padding 16

# ── Setup (hidden from recording) ──────────────────────────────────────────────
Hide
Type "export VAULT_ADDR=http://localhost:8200 VAULT_TOKEN=root"
Enter
Type "export PATH=$(pwd)/build:$PATH"
Enter
Type "clear"
Enter
Sleep 500ms
Show

# ── Intro comment ──────────────────────────────────────────────────────────────
Type "# vsh – an interactive HashiCorp Vault shell"
Enter
Sleep 1000ms
Type "export VAULT_ADDR=http://localhost:8200 VAULT_TOKEN=root"
Enter
Sleep 1500ms

# ── Launch vsh ────────────────────────────────────────────────────────────────
Type "vsh"
Enter
Sleep 2500ms

# ── Explore top-level paths ───────────────────────────────────────────────────
Type "ls"
Enter
Sleep 1500ms

# ── Navigate into prod ────────────────────────────────────────────────────────
Type "cd secret/prod/"
Enter
Sleep 800ms

Type "ls"
Enter
Sleep 1200ms

# ── Inspect a secret ─────────────────────────────────────────────────────────
Type "cat database"
Enter
Sleep 1500ms

# ── Search across all paths ───────────────────────────────────────────────────
Type "grep -k password /secret/"
Enter
Sleep 2000ms

# ── Copy prod → dr (uses goroutine worker pool) ───────────────────────────────
Type "cp -r /secret/prod /secret/dr"
Enter
Sleep 2500ms

Type "ls /secret/dr/"
Enter
Sleep 1200ms

# ── Preview replacements (dry-run) ────────────────────────────────────────────
Type "replace --dry-run prod dr /secret/dr/"
Enter

# ── End ───────────────────────────────────────────────────────────────────────
Sleep 4s
