# rm

```text
rm [--worker-count N] PATH
```

Remove `PATH`. If executed on a node, (i.e., a path ending with `/`), then remove is applied recursively.

Recursive removes run concurrently via a goroutine worker pool. The number of workers can be tuned with `--worker-count` (default: 10).
