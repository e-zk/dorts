# dorts

Dotfile template engine based around Go's `text/template`.

## configuration directory
By default the config file and template files are stored in `${XDG_CONFIG_HOME}/dorts` (usually `${HOME}/.config/dorts`). The config directory can be overridden by either the `DORTS_DIR` environment variable or the `-c` command-line flag. The priority of the config directory options are, in order:

* `-c <path>` command-line flag (highest priority - overrides everything);
* the `DORTS_DIR` environment variable; if it is not defined, then
* `${XDG_CONFIG_HOME}/dorts`; if `$XDG_CONFIG_HOME` is not defined, then
* `${HOME}/.config/dorts`.

## example setup

Take the following example config and template.

`${DORTS_DIR}/dorts.toml`:
```toml
# common variables specified here are 'global',
# every template will substitute these in.
# they can be overridden per-program.
[common]
background = "#101a1f"
foreground = "#ffffff"
accent = "#b00050"

# program configuration
[cwmrc]
path = "~/.cwmrc"         # output file path
background = "#1f2a2a"    # override global 'background'
gaps = "40 80 40 40"
```

`${DORTS_DIR}/cwmrc.tmpl`:
```console

# gaps
gap {{ .gaps }}

# colors
color menubg  "{{ .background }}"
color font "  "{{ .foreground }}"
color menufg  "{{ .background }}"
color selfont "{{ .accent }}"

# ... rest of config ...
```

After running `dorts`, the config file outputted to `~/.cwmrc` will look like:
```console

# gaps
gap 40 80 40 40

# colors
color menubg  "#1f2a2a"
color font    "#ffffff"
color menufg  "#1f2a2a"
color selfont "#b00050"

# ... rest of config ...
```

## command-line usage
Not implemented yet.

run/execute templates:
```console
$ dorts run
```

