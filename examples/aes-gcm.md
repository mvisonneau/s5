# AES/GCM

## Generate a key

There are probably many ways of generating random hexadecimal keys, here are a couple examples using `openssl` or `hexdump`

```bash
# 128B
~$ openssl rand -hex 16
0b9f3051e89ba4d5018fc7fb96fe5b44

# 256B
~$ openssl rand -hex 32
e52914b7dabf62cfe2c481ce0fff7ef6ac3034577454917eda80bead8c4e373c

# 512B
~$ openssl rand -hex 64
c7b39ab8eb67db22f24bbf2fe18226994583f73be0dedc895e50d12b607f101cc9f9a58126f78148880aefa7ed406c6ccffe3e959c5cb79a83b735b0b02f8c48

# Or using hexdump (eg: 128B)
~$ hexdump -n 16 -e '4/4 "%08X" 1 "\n"' /dev/urandom
66A35214AC00216C4D03E5FF747F936E
```

## Configure

For ease of use, you can export your key as an environment variable

```bash
~$ export S5_AES_KEY=0b9f3051e89ba4d5018fc7fb96fe5b44
```

Otherwise you will need to use the `--key` flag on each of your commands

```bash
~$ s5 cipher aes --key 0b9f3051e89ba4d5018fc7fb96fe5b44 "foo"
```

## Usage

```bash
# Set S5_AES_KEY variable
~$ export S5_AES_KEY=0b9f3051e89ba4d5018fc7fb96fe5b44

# Cipher text
~$ s5 cipher aes foo
{{ s5:YTdlOTQ2M2VhNzE1MGQ3NzlkYTRkZGRhOTM1MjEzMDBkOTNjNzY6ODhhNzI2NGUzZTllZjgwYTAyNWVhOWRm }}

# Decipher text
~$ s5 decipher aes "{{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}"
bar

# Store it anywhere in your files
~$ cat > example.yml <<EOF
---
var1: {{ s5:YTdlOTQ2M2VhNzE1MGQ3NzlkYTRkZGRhOTM1MjEzMDBkOTNjNzY6ODhhNzI2NGUzZTllZjgwYTAyNWVhOWRm }}
var2: {{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}
{{ s5:YjVhNTI4YTcwNzk1MGIyNWQ5MjhmN2FjMzFjZTRlZTllMWQwYWExYzY1ODQ3M2U2MDQyZTpmZDg5YjQ0MThmMjVkYTg5YjUyYjIxYjU= }}: {{ s5:NjMzNGEwYTUzNTIwMWJiMWNmNWMwNmFkM2EyZmExODI2YzI1NWFhMDg3OWU1NzI0NGM3NjNlY2Q6YzY2NTEyN2JmMjEwZjFlMDI1OGQwMmRk }}
EOF

# Render!
~$ s5 render aes example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher aes | s5 decipher aes
foo
~$ echo "foo: {{ s5:ZjMxN2JhYTBjNWE1OWQ5N2Q3MzBhMWEwOGIxYzBkNDQ0NDljMjY6ZTY5MDY1YzU3YTU1ZjViMzhmZDg3MTNj }}" | s5 render aes
foo: bar
```