# Vault

## Have a working Vault installation with sufficient permissions

Vault can run on most common platforms and is quite easy to get started with, specially for testing purposes. Have a look onto [Hashicorp's documentation](https://www.vaultproject.io/docs/install/index.html)

If you have docker running on your machine you can also do the following which I find very handy

```bash
~$ docker run -it --rm -p 8200:8200 vault:latest
```

Search or grep for `Root Token` in the logs in order to be able to authenticate against it:

```bash
~$ export VAULT_ADDR='http://127.0.0.1:8200'
~$ export VAULT_TOKEN=s.5Jc2jilLwYcTvHiKyTY0UNVu
~$ vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.1.3
Cluster Name    vault-cluster-485da001
Cluster ID      1b8fc2cd-e09c-073e-279f-6f5cef08584b
HA Enabled      false
```

You are all set!

## Configure a transit key

```bash
# Enable the transit backend
~$ vault secrets enable transit
Success! Enabled the transit secrets engine at: transit/

# Generate a key
~$ vault write -f transit/keys/foo
Success! Data written to: transit/keys/foo
```

For more detailed information, you can refer to [Hashicorp's documentation](https://www.vaultproject.io/docs/secrets/transit/index.html).

## Configure

In order to connect onto Vault, `s5` will use the same mechanisms that the `vault` agent uses. The easiest is to export those couple environment variables

```bash
~$ export VAULT_ADDR='http://127.0.0.1:8200'
~$ export VAULT_TOKEN=s.5Jc2jilLwYcTvHiKyTY0UNVu
```

`s5` will also lookup into `~/.vault-token` for the token if the file exists.

The only remaining thing is to specify the name of the key, by default `s5` will use the `default` one. You can specify which key to use with the following variable

```bash
~$ export S5_VAULT_TRANSIT_KEY=foo
```

Otherwise you will need to use the `--transit-key` flag on each of your commands

```bash
~$ s5 cipher vault --transit-key foo "bar"
```

## Usage

```bash
# Set VAULT_ADDR, VAULT_TOKEN and S5_VAULT_TRANSIT_KEY variables
~$ export VAULT_ADDR='http://127.0.0.1:8200'
~$ export VAULT_TOKEN=s.5Jc2jilLwYcTvHiKyTY0UNVu
~$ export S5_VAULT_TRANSIT_KEY=foo

# Cipher text
~$ s5 cipher vault foo
{{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}

# Decipher text
~$ s5 decipher vault "{{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}"
bar

# Store it anywhere in your files
~$ cat > example.yml <<EOF
---
var1: {{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}
var2: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}
{{ s5:rOY/SepnPdoWJs+uo+JdPEyzS3TnXcLDKbVjRsp5a0zadVU9s7M= }}: {{ s5:smx5lDxYeRUL9pwq7LyAT0YfmfGFDAfHPQJ5faTsw4HCAdosdLc= }}
EOF

# Render!
~$ s5 render vault example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher vault | s5 decipher vault
foo
~$ echo "foo: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}" | s5 render vault
foo: bar
```
