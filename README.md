ctar
===

`ctar` is a CLI tool to archive a directory into a `tar.gz` file with the option to specify a maximum size of the files to be included in the archived.

When the size (`-s`) option is specified `ctar` will start archiving files in the specified directory from the newest to the oldest until the maximum size is reached. The size refers to the size of the files before compression.

The tool is useful to package cache files for CI when other tools don't provide an option to limit the cache leading the cache file to grow unbound in size with every build. One such tool is go where cache files are kept 5 days without possibility of changing.

Example usage:

```bash
ctar -s 10MB archive.tar.gz dir_to_archive
```

or via docker image with

```bash
docker run --rm -v /some_dir_path:/workspace dedalusj/ctar /root/ctar -s 10MB /workspace/archive.tar.gz /workspace/dir_to_archive
```
which will create the archive file at path `/some_dir_path/archive.tar.gz` on the host machine.