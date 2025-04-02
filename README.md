# JSON tools: fmt, diff

## Commands

### cmd/jsonfmt

Normal JSON pretty print, except if an array does not contain other arrays or objects, then it is printed on a single line.

### cmd/jsondiff

Diff JSON logically. Object order is ignored. Array order is compared.

## Use

**jsonfmt**

```sh
go run github.com/kardianos/json/cmd/jsonfmt@v1.0.2 -w file.json
```

**jsondiff**

```sh
go run  github.com/kardianos/json/cmd/jsonfmt@v1.0.2 file1.json file2.json
```
