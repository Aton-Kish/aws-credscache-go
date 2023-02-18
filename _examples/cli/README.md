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

### `sdkv2`

#### `sdkv2 cache`

```shell
go run main.go sdkv2 cache --profile myprofile
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
go run main.go sdkv2 cache --profile myprofile
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
    "ResultMetadata": {}
}
```

##### Share the cache with AWS CLI

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
go run main.go sdkv2 cache --profile myprofile
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
    "ResultMetadata": {}
}
```

#### `sdkv2 nocache`

An MFA token code will be requested every time you run it.

```shell
go run main.go sdkv2 nocache --profile myprofile
Assume Role MFA token code: 123456
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
    "ResultMetadata": {}
}
```

### `sdkv1`

#### `sdkv1 cache`

```shell
go run main.go sdkv1 cache --profile myprofile
Assume Role MFA token code: 123456
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
}
```

Re-running it will never ask for an MFA token code again.

```shell
go run main.go sdkv1 cache --profile myprofile
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
}
```

##### Share the cache with AWS CLI

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
go run main.go sdkv1 cache --profile myprofile
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
}
```

#### `sdkv1 nocache`

An MFA token code will be requested every time you run it.

```shell
go run main.go sdkv1 nocache --profile myprofile
Assume Role MFA token code: 123456
{
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:assumed-role/myrole/mysessionname",
    "UserId": "AXXXXXXXXXXXXXXXXXXXX:mysessionname",
}
```
