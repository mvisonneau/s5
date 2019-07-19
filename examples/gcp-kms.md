# GCP KMS

## Get a KMS key with sufficient privileges provisioned

First thing your need to do is being able to get programmatic access on top the gcp APIs and get a KMS key provisioned.

`s5` should support all the authentication mechanisms supported by the [google-cloud-go](https://github.com/googleapis/google-cloud-go) libraries.

If you have `gcloud` available

```bash
~$ gcloud auth login
```

You can create your key pair and correct IAM permissions from your preferred method (CLI, API, Terraform, or else)

```bash
# First you need a keyring
~$ gcloud kms keyrings create mykeyring \
      --location=global

# Then add a new key to your keyring
~$ gcloud kms keys create mykey \
      --location=global \
      --keyring=mykeyring \
      --purpose=encryption
```

### IAM roles

- `Cloud KMS CryptoKey Encrypter` is required by the `cipher` function
- `Cloud KMS CryptoKey Decrypter` is required by both the `decipher` and `render` functions

## Configure

For ease of use, you can export your key as an environment variable

```bash
~$ export S5_GCP_KMS_KEY_NAME="projects/myproject/locations/global/keyRings/mykeyring/cryptoKeys/mykey"
```

Otherwise you will need to use the `--kms-key-name` flag on each of your commands

```bash
~$ s5 cipher gcp --kms-key-name "projects/myproject/locations/global/keyRings/mykeyring/cryptoKeys/mykey" foo
```

## Usage

```bash
# Get GCP credentials and set the S5_GCP_KMS_KEY_NAME
~$ gcloud auth login
~$ export S5_GCP_KMS_KEY_NAME="arn:gcp:kms:*:111111111111:key/mykey"

# Cipher text
~$ s5 cipher gcp foo
{{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}

# Decipher text
~$ s5 decipher gcp "{{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}"
bar

# Store it anywhere in your files
~$ cat > example.yml <<EOF
---
var1: {{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}
var2: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}
{{ s5:rOY/SepnPdoWJs+uo+JdPEyzS3TnXcLDKbVjRsp5a0zadVU9s7M= }}: {{ s5:smx5lDxYeRUL9pwq7LyAT0YfmfGFDAfHPQJ5faTsw4HCAdosdLc= }}
EOF

# Render!
~$ s5 render gcp example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher gcp | s5 decipher gcp
foo
~$ echo "foo: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}" | s5 render gcp
foo: bar
```
