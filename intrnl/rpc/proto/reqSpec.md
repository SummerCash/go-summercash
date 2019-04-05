# JSON Request Specifications

Note: all requests should be formatted in the following format: 

## Common
URL: ```localhost:<port>/twirp/common.Common/```

Each common GeneralRequest takes an optional amount of parameters (each RPC method uses the same request)

```JSON
{
    "input": "BYTE_VALUE_INPUT",
    "s": "STRING_VALUE_INPUT"
}
```

## Chain
URL: ```localhost:<port>/twirp/chain.Chain/```

```JSON
{
    "address": "STRING_ADDRESS_INPUT (e.g. 0x000000...)"
}
```

## Accounts
URL: ```localhost:<port>/twirp/accounts.Accounts/```

```JSON
{
    "address": "STRING_ADDRESS_INPUT (e.g. 0x000000...)",
    "privateKey": "STRING_PRIVATE_KEY_INPUT (e.g. 0x000000...)"
}
```

## Config
URL: ```localhost:<port>/twirp/config.Config/```

```JSON
{
    "genesisPath": "STRING_PATH_TO_GENESIS.JSON"
}
```

## Coordination Chain
URL: ```localhost:<port>/twirp/coordinationChain.CoordinationChain/```

Note: the coordination chain package does not have any requests that require inputs

```JSON
{
}
```

## Crypto
URL: ```localhost:<port>/twirp/crypto.Crypto/```

```JSON
{
    "input": "BYTE_VALUE_INPUT",
    "n": "UINT_INPUT"
}
```

## Transaction
URL: ```localhost:<port>/twirp/transaction.Transaction/```

```JSON
{
    "nonce": "UINT_NONCE_INPUT",
    "address": "STRING_ADDRESS_INPUT (e.g. 0x000000...)",
    "address2": "STRING_ADDRESS_INPUT (e.g. 0x000000...)",
    "amount": "FLOAT64_INPUT",
    "payload": "BYTE_PAYLOAD_INPUT"
}
```

## UPnP
URL: ```localhost:<port>/twirp/upnp.UPnP/```

```JSON
{
    "portNumber": "UINT_32_INPUT"
}
```