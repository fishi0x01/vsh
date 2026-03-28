# cp

```text
cp [-r] SOURCE TARGET
```

Copy `SOURCE` path to `TARGET` path. If executed on a node, (i.e., a path ending with `/`), then copy is applied recursively.

Recursive copies run concurrently via a goroutine worker pool. The number of workers can be tuned with the `--worker-count` flag when starting `vsh` (default: 10).
