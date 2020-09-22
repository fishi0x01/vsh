## v0.6.3 (September 22, 2020)

BUG FIXES:

* Properly handle ambiguous files without read permission [#46](https://github.com/fishi0x01/vsh/pull/46) - Thank you for detailed report [agaudreault-jive](https://github.com/agaudreault-jive)

## v0.6.2 (June 18, 2020)

ENHANCEMENTS:

* Allow separation of command params with multiple whitespaces ([#40](https://github.com/fishi0x01/vsh/pull/40)) - Thank you [vikin91](https://github.com/vikin91)
* Unify verbose behavior ([#41](https://github.com/fishi0x01/vsh/pull/41)) - Thank you [vikin91](https://github.com/vikin91)

BUG FIXES:

* Properly handle ambiguous files ([#37](https://github.com/fishi0x01/vsh/pull/37))

## v0.6.1 (June 15, 2020)

BUG FIXES:

* Properly handle suffix '/' in source path ([#39](https://github.com/fishi0x01/vsh/pull/39))

## v0.6.0 (June 13, 2020)

ENHANCEMENTS:

* add `append` command ([#30](https://github.com/fishi0x01/vsh/issues/30)) - Thank you [vikin91](https://github.com/vikin91)

BUG FIXES:

* Remove file/dir ambiguity for `rm` ([#29](https://github.com/fishi0x01/vsh/issues/29))

## v0.5.0 (April 5, 2020)

ENHANCEMENTS:

* add `grep` command ([#25](https://github.com/fishi0x01/vsh/issues/25))
* `ls` with new line instead of single line ([#27](https://github.com/fishi0x01/vsh/issues/27))

BUG FIXES:

* remove `//` from paths ([#26](https://github.com/fishi0x01/vsh/issues/26))
* fix broken tests

## v0.4.1 (March 21, 2020)

ENHANCEMENTS:

* performance: cache `List()` queries ([#23](https://github.com/fishi0x01/vsh/issues/23))

## v0.4.0 (March 10, 2020)

ENHANCEMENTS:

* use TokenHelper mechanism ([#20](https://github.com/fishi0x01/vsh/issues/20))

## v0.3.1 (October 31, 2019)

BUG FIXES:

* fix top-level path panic ([#17](https://github.com/fishi0x01/vsh/issues/17))

## v0.3.0 (October 31, 2019)

ENHANCEMENTS:

* token list permission on sys/mounts is not mandatory

## v0.2.0 (October 20, 2019)

ENHANCEMENTS:

* use `~/.vault-token` as fallback if `VAULT_TOKEN` is not set ([#12](https://github.com/fishi0x01/vsh/issues/12))

BUG FIXES:

* error handling to catch bad input ([#13](https://github.com/fishi0x01/vsh/issues/13))

## v0.1.1 (October 8, 2019)

BUG FIXES:

* more sanity checks on user input to avoid crashes ([#10](https://github.com/fishi0x01/vsh/issues/10))

## v0.1.0 (October 7, 2019)

Initial release
