# 1Password

Any secret that `s5` accepts on the command line (currently the AES key) can be
replaced by a [1Password secret reference](https://developer.1password.com/docs/cli/secret-reference-syntax/),
so that it never has to be stored in your shell profile, your environment or on
disk.

`s5` resolves those references by invoking the [1Password CLI](https://developer.1password.com/docs/cli/get-started/)
(`op`), which therefore needs to be installed, available in your `PATH` and
signed in. On a workstation this usually means enabling the desktop app
integration (**1Password → Settings → Developer → Integrate with 1Password
CLI**).

## Store the key in 1Password

Any item type works, `op read` returns whichever field you point it at. Keep the
field **concealed** (`[password]`) rather than plain text: it gets masked in the
UI and keeps a value history, so rotating the key does not lose the previous one
while older ciphered files still refer to it.

Using the `API Credential` category, which is meant for this kind of secret

```bash
~$ op item create \
  --category="API Credential" \
  --vault=Private \
  --title=s5 \
  "credential[password]=$(openssl rand -hex 32)"
```

Or the `Password` category, which lets you name the field yourself and holds
several keys on a single item

```bash
~$ op item create \
  --category=password \
  --vault=Private \
  --title=s5 \
  "aes-key[password]=$(openssl rand -hex 32)"
```

In both cases the key itself never lands in your shell history, only the
`openssl` command does. You can then get the reference of the field with

```bash
~$ op item get s5 --vault=Private --format json | jq -r '.fields[] | select(.value != null) | .reference'
op://Private/s5/credential
```

## Usage

Use the reference wherever you used to pass the key itself

```bash
~$ s5 cipher aes --key "op://Private/s5/credential" foo
{{s5:YTdlOTQ2M2VhNzE1MGQ3NzlkYTRkZGRhOTM1MjEzMDBkOTNjNzY6ODhhNzI2NGUzZTllZjgwYTAyNWVhOWRm}}
```

Or, more conveniently, export the reference instead of the key. Contrary to the
key, the reference is not sensitive and can safely live in your `~/.zshrc`,
`~/.bashrc` or in a committed `.envrc`

```bash
~$ export S5_AES_KEY="op://Private/s5/credential"
~$ echo "foo" | s5 cipher aes | s5 decipher aes
foo
```

The value is only fetched when a command actually needs it, and 1Password
unlocks it the way it is configured to: through Touch ID / the desktop app when
you are on a workstation, and using `OP_SERVICE_ACCOUNT_TOKEN` when running
unattended, eg: in CI

```bash
~$ export OP_SERVICE_ACCOUNT_TOKEN="ops_..."
~$ export S5_AES_KEY="op://CI/s5/credential"
~$ s5 render aes example.yml
```

## Reference syntax

```
op://<vault>/<item>/<field>
op://<vault>/<item>/<section>/<field>
```

Vault, item and field names may contain spaces, in which case the reference
needs to be quoted. Values that do not start with `op://` are used as-is, so
this is entirely opt-in and does not change any existing setup.
