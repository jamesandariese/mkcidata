# `mkcidata`

Make a CIDATA ISO (without figuring out what package has mkisofs today).

## Usage

```bash
./mkcidata cidata.iso user-data meta-data
```

The user-data and meta-data files may be in any directory.  The path will
be stripped prior to adding to the ISO.

## Installation

```bash
go install github.com/jamesandariese/mkcidata@latest
```

## License

See [LICENSE](./LICENSE) file.
