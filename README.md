# KVS

KVS is a commandline tool to store and organize key-value data on your local file system.

- Built on top of the amazing [bbolt](https://github.com/etcd-io/bbolt) db
- Free open source software
- Works on [Linux](https://github.com/lucasepe/kvs/releases/latest), [Mac OSX](https://github.com/lucasepe/kvs/releases/latest), [Windows](https://github.com/lucasepe/kvs/releases/latest)
- Just a single portable binary file

## Store

A _store_ is a single file on your local file system.

KVS save all your key-values data grouped by _buckets_ in a specific _store_.

You can specify the _store_ name using the `--store` (or the short version `-s`) flag.

- each store is located in your `$HOME/.kvs` folder

## Buckets

KVS uses _buckets_ to organize your data. 

You can specify a _bucket_ using the `--bucket` (or the short version `-b`) flag.

- if you are _pushing_ a key-val pair and the bucket does not exists, it will be created
- you cannot nest _buckets_

## Keys

### Slugs

Performing a `push`, `pull` or `del` command, all keys are transformed into _slugs_.

- transliterate Unicode characters into alphanumeric strings

- all punctuation is stripped and whitespace between words are replaced by hyphens

Example: a key named `Hello Wonderful World!` became `hello-wonderful-world`.

Also bucket names are transformed into _slugs_.

## Values 

### Encryption

KVS can encrypt values using the [AES](https://it.wikipedia.org/wiki/Advanced_Encryption_Standard) algorithm in [Galois Counter Mode (GCM)](https://en.wikipedia.org/wiki/Galois/Counter_Mode).

 - the result will be saved as base64 encoded string

If you want to do so, just add the `--encrypt` (or the short version `-e`) flag.

```bash
$ kvs push track-id UA-XXXXXXX-X -s accounts -b google -e
Secret phrase: 
Secret phrase again:
item successfully stored in bucket 'google' with key 'track-id'
```

Pulling the value without decription:

```bash
$ kvs pull track-id -s accounts -b google
zRl1TZZe1JVpfAtY1yFU1g==
```

to decrypt the value you can use the `--decrypt` (or the short version `-d`) flag

```bash
$ kvs pull track-id -s accounts -b google -d
Secret phrase: 
UA-XXXXXXX-X
```

### Binary values

Values ​​can also be binary data (up to 1MB).

## Use cases

- configuration parameters for others local tools and apps
- credentials (using the encryption feature)

```text
$ kvs
 _
| | __  __   __    ___ 
| |/ /  \ \ / /  / __|
|   <    \ V /   \__ \
|_|\_\ey  \_/ al |___/ tore

Usage:
  kvs [command]

Available Commands:
  del         Removes from a store the item with the specified key from a bucket
  help        Help about any command
  list        List all bucket names in a store or all key names for a specific bucket
  pull        Fetch from a store the item with the specified key in a bucket
  push        Update a store adding an item with the specified key in a bucket

Flags:
  -h, --help           help for kvs
  -s, --store string   store name (default "vault")
      --version        version for kvs

Use "kvs [command] --help" for more information about a command.
```

### How to store an item

Example: add a property `user=john.doe@gmail.com` in a bucket called `google` and a store called `accounts`

```bash
$ kvs push --store accounts --bucket google user luca.sepe@gmail.com
item successfully stored in bucket 'google' with key 'user'
```

Example: add a property using shell pipes

```bash
$ pwgen | kvs push --store accounts --bucket google pass
item successfully stored in bucket 'google' with key 'pass'
```

### How to retrieve an item

Example: retrieve the value of the `user` property in the bucket `google`

```bash
$ kvs pull --store accounts --bucket google user
john.doe@gmail.com
```

Example: retrieve the encrypted password and pipe to clipboard

```bash
$ kvs -b aruba pull pass -d | xclip -selection c
Secret phrase: 
```

the decrypted password will be saved to your clipboard - ready to be pasted!

### How to delete an item

```bash
$ kvs del --store accounts --bucket google hello
item with key 'hello' successfully removed from bucket 'google'
```

## TODO

- [ ] encrypt/decrypt secret phrase alternative (using a private key file???)
- [ ] implement an `env` command in order to expose a key-val item as environment variable