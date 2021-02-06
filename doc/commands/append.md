# append

```text
append [-f|--force] [-s|--skip] [-r|--rename] SOURCE TARGET
```

Append operation reads secrets from `SOURCE` and merges it to `TARGET`.
The `TARGET` will be created with a placeholder value if it does not exists.
Both `SOURCE` and `TARGET` must be leaves (path cannot end with `/`).

By default, `append` does not overwrite secrets if the `TARGET` already contains a key.
The default behavior can be explicitly set using flag: `-s` or `--skip`. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append --skip /secret/from /secret/to

> cat /secret/to

fruit=pear
vegetable=tomato
tree=oak
```

Setting flag `-f` or `--force` will cause the conflicting keys from the `<to-secret>` to be overwritten with keys from the `<from-secret`>. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append -f /secret/from /secret/to

> cat /secret/to

fruit=apple
vegetable=tomato
tree=oak
```

Setting flag `-r` or `--rename` will cause the conflicting keys from the `<to-secret>` to be kept as they are. Instead the keys from the `<from-secret`> will be stored under a renamed key. Example:

```bash
> cat /secret/from

fruit=apple
vegetable=tomato

> cat /secret/to

fruit=pear
tree=oak

> append -r /secret/from /secret/to

> cat /secret/to

fruit=pear
fruit_1=apple
vegetable=tomato
tree=oak
```
