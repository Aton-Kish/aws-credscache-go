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

#### `sdkv2 nocache`

An MFA token code will be requested every time you run it.
This is a very common problem with CLIs using the AWS SDK v2.

![`sdkv2 nocache`](./images/gif/sdkv2_nocache.gif)

```mermaid
sequenceDiagram
    autonumber

    actor U as User
    participant CLI as CLI
    participant CC as Credentials Cache<br>(In-Memory Cache)
    participant ARP as Assume Role Provider
    participant AWS as AWS

    U ->>+ CLI: execute CLI

    CLI ->>+ CC: retrieve credentials
    alt exist valid credentials cache in memory
        CC ->> CC: load credentials cache from memory
    else
        note over U, AWS: This flow will always be executed the first time, as there is no cache in memory.
        CC ->>+ ARP: retrieve credentials
        ARP ->>+ CLI: ask MFA token code
        CLI ->>+ U: ask MFA token code
        U -->>- CLI: MFA token code
        CLI -->>- ARP: MFA token code
        ARP ->>+ AWS: call sts:AssumeRole API with MFA token code
        AWS -->>- ARP: credentials
        ARP -->>- CC: credentials
        CC ->> CC: store credentials cache in memory
    end

    CC -->>- CLI: credentials
    CLI ->>+ AWS: call AWS API with credentials
    AWS -->>- CLI: response
    CLI -->>- U: output
```

#### `sdkv2 cache`

A credentials cache can be shared between processes.

![`sdkv2 cache`](./images/gif/sdkv2_cache.gif)

It can also be shared with the AWS CLI.

![`sdkv2 cache` shared with AWS CLI](./images/gif/sdkv2_cache_awscli.gif)

```mermaid
sequenceDiagram
    autonumber

    actor U as User
    participant CLI as CLI
    participant CC as Credentials Cache<br>(In-Memory Cache)
    participant FC as File Cache Provider
    participant ARP as Assume Role Provider
    participant AWS as AWS


    U ->>+ CLI: execute CLI

    CLI ->>+ CC: retrieve credentials
    alt exist valid credentials cache in memory
        CC ->> CC: load credentials cache from memory
    else
        note over U, AWS: This flow will always be executed the first time, as there is no cache in memory.

        CC ->>+ FC: retrieve credentials

        alt exist valid credentials cache in file system
            FC ->> FC: load credentials cache from file system
        else
            FC ->>+ ARP: retrieve credentials
            ARP ->>+ CLI: ask MFA token code
            CLI ->>+ U: ask MFA token code
            U -->>- CLI: MFA token code
            CLI -->>- ARP: MFA token code
            ARP ->>+ AWS: call sts:AssumeRole API with MFA token code
            AWS -->>- ARP: credentials
            ARP -->>- FC: credentials
            FC ->> FC: store credentials cache in file system
        end

        FC ->>- CC: credentials
            CC ->> CC: store credentials cache in memory
    end

    CC -->>- CLI: credentials
    CLI ->>+ AWS: call AWS API with credentials
    AWS -->>- CLI: response
    CLI -->>- U: output
```

### `sdkv1`

The concept is the same as SDK v2.

#### `sdkv1 nocache`

![`sdkv1 nocache`](./images/gif/sdkv1_nocache.gif)

#### `sdkv1 cache`

![`sdkv1 cache`](./images/gif/sdkv1_cache.gif)

![`sdkv1 cache` shared with AWS CLI](./images/gif/sdkv1_cache_awscli.gif)
