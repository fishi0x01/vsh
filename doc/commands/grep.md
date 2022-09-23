# grep

```text
grep [-e|--regexp] [-k|--keys] [-v|--values] [-S|--shallow] SEARCH [PATH]
```

`grep` recursively searches the given `SEARCH` substring in key and value pairs of given `PATH`. To treat the search string as a regular-expression, add `-e` or `--regexp` to the command. By default, both keys and values will be searched. If you would like to limit the search, you may add `-k` or `--keys` to the end of the command to search only a path's keys, or `-v` or `--values` to search only a path's values.
If PATH is not specified, the currently active path is taken.
If you do not desire a deep recursive search, you can use the `-S` or `--shallow` flag to limit searching to only the leaf nodes of the given path.
If you are looking for copies or just trying to find the path to a certain string, this command might come in handy.
