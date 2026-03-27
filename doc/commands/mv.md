# mv

```text
mv [--worker-count N] SOURCE TARGET
```

Move `SOURCE` path to `TARGET` path. If executed on a node, (i.e., a path ending with `/`), then move is applied recursively.

Recursive moves run concurrently via a goroutine worker pool. The number of workers can be tuned with `--worker-count` (default: 10).
