# dorts

Dotfile template engine based around Go's `text/template`.

## example setup

dorts.toml:
```toml
# common variables specified here are 'global'
# every program's template will substitute these
# they can be overridden per-program
[common]
background = "#101a1f"
foreground = "#ffffff"
accent = "#b00050"

# program configuration
[cwmrc]
path = "~/.myprogram.conf" # output file path
background = "#1f2a2a"     # override global 'background'
```

cwmrc.tmpl:
```console

# gaps
gap 40 80 40 40

# colors
color menubg  "{{ .background }}"
color font "  "{{ .foreground }}"
color menufg  "{{ .background }}"
color selfont "{{ .accent }}"

$ ... rest of config ...

```

## command-line usage

run/execute templates:
```console
$ dorts run
```

