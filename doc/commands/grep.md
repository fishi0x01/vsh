# grep

```text
grep [-e|--regexp] [-k|--keys] [-v|--values] SEARCH PATH
```

`grep` recursively searches the given `SEARCH` substring in key and value pairs of given `PATH`. To treat the search string as a regular-expression, add `-e` or `--regexp` to the end of the command. By default, both keys and values will be searched. If you would like to limit the search, you may add `-k` or `--keys` to the end of the command to search only a path's keys, or `-v` or `--values` to search only a path's values.
 If you are looking for copies or just trying to find the path to a certain string, this command might come in handy.
