# JSON Request Specifications

Note: all requests should be formatted in the following format: 

## Common

Each common GeneralRequest takes an optional amount of parameters (each RPC method uses the same request)

```JSON
{
    "input": "BYTE_VALUE_INPUT",
    "s": "STRING_VALUE_INPUT"
}
```

## Chain

```JSON
{
    "address": "STRING_ADDRESS_INPUT (e.g. 0x000000...)"
}
```

## Accounts

```JSON
{
    "address": "STRING_ADDRESS_INPUT (e.g. 0x000000...)",
    "privateKey": "STRING_PRIVATE_KEY_INPUT (e.g. 0x000000...)"
}
```

## Config

```JSON
{
    "genesisPath": "STRING_PATH_TO_GENESIS.JSON"
}
```

## Coordination Chain

Note: the coordination chain package does not have any requests that require inputs

```JSON
{
}
```

## Crypto

```JSON
{
    "input": "BYTE_VALUE_INPUT",
    "n": "UINT_INPUT"
}
```

## Transaction

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

```JSON
{
    "portNumber": "UINT_32_INPUT"
}
```