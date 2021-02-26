# add

```text
add [-f|--force] [-y|--confirm] [-n|--dry-run] KEY VALUE PATH
```

Add operation adds or overwrites a single key at path.

By default, it will add the key if it does not already exist. To overwrite the key if it exists, use the `--force` flag.

```bash
> cat secret/path

value = 1
other = thing

> add fizz buzz secret/path
> cat /secret/to

value = 1
fizz = buzz
other = thing
```
