# CLI

## Setup

### go mod

```shell
go mod tidy
```

### AWS Config

Edit `~/.aws/config` as needed.

```ini
[profile myprofile]
source_profile = default
mfa_serial = arn:aws:iam::123456789012:mfa/myname
role_arn = arn:aws:iam::123456789012:role/myrole
role_session_name = myname
region = us-east-1
output = json
```

## Run commands

### `cache` command

```shell
go run main.go cache --profile myprofile
Assume Role MFA token code: 123456
{
  "Account": "123456789012",
  "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
  "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
  "ResultMetadata": {}
}
```

Re-running it will never ask for an MFA token code again.

```shell
go run main.go cache --profile myprofile
{
  "Account": "123456789012",
  "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
  "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
  "ResultMetadata": {}
}
```

#### Share the cache with AWS CLI

```shell
aws sts get-caller-identity --profile myprofile
Enter MFA code for arn:aws:iam::123456789012:mfa/myname:
{
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname"
}
```

```shell
go run main.go cache --profile myprofile
{
  "Account": "123456789012",
  "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
  "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
  "ResultMetadata": {}
}
```

### `nocache` command

An MFA token code will be requested every time you run it.

```shell
go run main.go nocache --profile myprofile
Assume Role MFA token code: 123456
{
  "Account": "123456789012",
  "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
  "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
  "ResultMetadata": {}
}
```
