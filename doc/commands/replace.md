# replace

```text
grep [-e|--regexp] [-k|--keys] [-v|--values] [-S|--shallow] [-n|--dry-run] [-y|--confirm] [-o|--output inline|diff] PATTERN REPLACE [PATH]
```

`replace` works similarly to `grep`, but has the ability to mutate data inside Vault. By default, confirmation is required before writing data. You may skip confirmation by using the `-y`/`--confirm` flags. Conversely, you may use the `-n`/`--dry-run` flags to skip both confirmation and any writes. Changes that would be made are presented in red (delete) and green (add) coloring.
To treat the search string as a regular-expression, add `-e` or `--regexp` to the command.
If PATH is not specified, the currently active path is taken.

If you do not desire a deep recursive search, you can use the `-S` or `--shallow` flag to limit searching to only the leaf nodes of the given path.

This command has two output formats available via the `--output` flag:

- `inline`: A colorized inline format where deletions are in red background text and additions are in green background text. This is the default.
- `diff`: A non-colorized format that prints changes in two lines prefixed with a `-` for before and `+` for after replacement. This is more useful for copying and pasting the result.
