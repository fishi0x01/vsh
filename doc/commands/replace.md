# replace

`replace` works similarly to `grep`, but has the ability to mutate data inside Vault. By default, confirmation is required before writing data. You may skip confirmation by using the `-y`/`--confirm` flags. Conversely, you may use the `-n`/`--dry-run` flags to skip both confirmation and any writes. Changes that would be made are presented in red (delete) and green (add) coloring.
