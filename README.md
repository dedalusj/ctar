ctar
===

`ctar` is a CLI tool to archive a directory into a `tar.gz` file with the option to specify a maximum size of the files to be archived.

When the size option is specified `ctar` will start archiving files in the specified directory from the newest to the newest until the maximum size is reached.

The tool is useful to package cache files for CI when other tools cannot guarantee a maximum size for the archive and you want to avoid an explosion of the cache size that would slow the CI build. One such example is Go that caches its build files for up to 5 days with no way of limiting the cache size.

Example usage:

```bash
ctar --size 10MB dir_to_archive archive.tar.gz
```