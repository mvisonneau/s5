# PGP

## Generate a keypair

The first thing that comes to mind is to use the `gpg` command line tool. There are probably many other ways though!

Beware that one of the current limitation is that your private key pair **must not be password encrypted**.

```bash
# Define the characteristics of your key
~$ cat > ~/foo.def <<EOF
Key-Type: 1
Key-Length: 4096
Name-Real: Foo Bar
Name-Email: foo@bar.com
Expire-Date: 0
EOF

# Generate it, when prompted do not set a passphrase
~$ gpg --batch --gen-key ~/foo.def
gpg: key BEB39418F246D0DB marked as ultimately trusted

# Export both public and private keys from your keystore
~$ gpg --export --armor foo@bar.com > ~/public.pem
~$ gpg --export-secret-key --armor foo@bar.com > ~/private.pem
```

## Configure

For ease of use, you can export your key as an environment variable

```bash
~$ export S5_PGP_PUBLIC_KEY_PATH=~/public.pem
~$ export S5_PGP_PRIVATE_KEY_PATH=~/private.pem
```

Otherwise you will need to use the `--public-key` and/or `--private-key` flags on each of your commands

```bash
~$ s5 cipher pgp --public-key ~/public.pem "foo"
~$ s5 decipher pgp --public-key ~/public.pem --private-key ~/private.pem "{{ s5:xxx }}"
~$ s5 render pgp --public-key ~/public.pem --private-key ~/private.pem example.txt
```

## Usage

```bash
# Set S5_PGP_PUBLIC_KEY_PATH and S5_PGP_PRIVATE_KEY_PATH variables
~$ export S5_PGP_PUBLIC_KEY_PATH=~/public.pem
~$ export S5_PGP_PRIVATE_KEY_PATH=~/private.pem

# Cipher text
~$ s5 cipher pgp foo
{{ s5:YTdlOTQ2M2VhNzE1MGQ3NzlkYTRkZGRhOTM1MjEzMDBkOTNjNzY6ODhhNzI2NGUzZTllZjgwYTAyNWVhOWRm }}

# Decipher text
~$ s5 decipher pgp "{{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}"
bar

# Store it anywhere in your files
~$ cat > example.yml <<EOF
---
var1: {{ s5:YTdlOTQ2M2VhNzE1MGQ3NzlkYTRkZGRhOTM1MjEzMDBkOTNjNzY6ODhhNzI2NGUzZTllZjgwYTAyNWVhOWRm }}
var2: {{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}
{{ s5:YjVhNTI4YTcwNzk1MGIyNWQ5MjhmN2FjMzFjZTRlZTllMWQwYWExYzY1ODQ3M2U2MDQyZTpmZDg5YjQ0MThmMjVkYTg5YjUyYjIxYjU= }}: {{ s5:NjMzNGEwYTUzNTIwMWJiMWNmNWMwNmFkM2EyZmExODI2YzI1NWFhMDg3OWU1NzI0NGM3NjNlY2Q6YzY2NTEyN2JmMjEwZjFlMDI1OGQwMmRk }}
EOF

# Render!
~$ s5 render pgp example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher pgp | s5 decipher pgp
foo
~$ echo "foo: {{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}" | s5 render pgp
foo: bar
```