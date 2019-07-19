# AWS KMS

## Get a KMS key with sufficient privileges provisioned

First thing your need to do is being able to get programmatic access on top the AWS APIs and get a KMS key provisioned.

`s5` should support all the authentication mechanisms supported by the [aws-go-sdk](https://aws.amazon.com/sdk-for-go/)

Assuming you got a regular access and secret keypair

```bash
~$ export AWS_ACCESS_KEY_ID=<access_key_id>
~$ export AWS_SECRET_ACCESS_KEY=<secret_access_key>
```

You can create your key pair and correct IAM permissions from your preferred method (CLI, API, Terraform, or else)

```bash
~$ aws kms create-key
```

### IAM policy

Here is the required statements in order for s5 to work properly with your KMS key

```json
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": [
      "kms:Encrypt",
      "kms:Decrypt"
    ],
    "Resource": [
      "arn:aws:kms:*:111111111111:key/mykey"
    ]
  }
}
```

If you want to manage permissions a bit more granuarily

- `kms:Encrypt` is required by the `cipher` function
- `kms:Decrypt` is required by both the `decipher` and `render` functions

## Configure

For ease of use, you can export your key as an environment variable

```bash
~$ export S5_AWS_KMS_KEY_ARN="arn:aws:kms:*:111111111111:key/mykey"
```

Otherwise you will need to use the `--kms-key-arn` flag on the `cipher` function. As the key ID is embedded in the ciphertext, we do not need to specify it for deciphering functions.

```bash
# Key ARN is required
~$ s5 cipher aws --kms-key-arn "arn:aws:kms:*:111111111111:key/mykey" "foo"

# Key ARN is not required
~$ s5 decipher aws "{{ s5:xxx }}"
~$ s5 render aws example.txt
```


## Usage

```bash
# Set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and S5_aws_TRANSIT_KEY variables
~$ export AWS_ACCESS_KEY_ID=<access_key_id>
~$ export AWS_SECRET_ACCESS_KEY=<secret_access_key>
~$ export S5_AWS_KMS_KEY_ARN="arn:aws:kms:*:111111111111:key/mykey"

# Cipher text
~$ s5 cipher aws foo
{{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}

# Decipher text
~$ s5 decipher aws "{{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}"
bar

# Store it anywhere in your files
~$ cat > example.yml <<EOF
---
var1: {{ s5:NKOAQ4chbVbkx8zxdxcG1hCmEbjszJ5Bx40gPtrtLQ== }}
var2: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}
{{ s5:rOY/SepnPdoWJs+uo+JdPEyzS3TnXcLDKbVjRsp5a0zadVU9s7M= }}: {{ s5:smx5lDxYeRUL9pwq7LyAT0YfmfGFDAfHPQJ5faTsw4HCAdosdLc= }}
EOF

# Render!
~$ s5 render aws example.yml
---
var1: foo
var2: bar
secret_key: secret_value

# s5 can also read from stdin
~$ echo "foo" | s5 cipher aws | s5 decipher aws
foo
~$ echo "foo: {{ s5:8GRhfpD7HttVPri0vDfADmcW8hhhzAdav/CfwagwJw== }}" | s5 render aws
foo: bar
```
