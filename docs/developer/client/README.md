# Golang DWN Client SDK

## Install

Installing the client is easy!

`go get -u github.com/openreserveio/dwn/go`

## Usage

### Create a Client for a specific DWN Node and optional Protocol

Import the client package

```golang

import (
    "github.com/openreserveio/dwn/go/client"
)

```

Create a new client, passing in a server URL:

```golang

func main() {

    dwnClient := client.CreateDWNClient("https://dwn.alpha.openreserve.io")


}

```

If you're adhering to a [Protocol](https://identity.foundation/decentralized-web-node/spec/#protocols), you can create a client within a context of a server, protocol, and protocol version:

```golang

func main() {

    dwnClient := client.CreateDWNClientForProtocol("https://dwn.alpha.openreserve.io", "https://openreserve.io/protocols/test-protocol.json", "v0.0.1")


}

```

After getting a DWN Client instantiated, you can use it to perform most of the operations in the DIF DWN Spec.


### Set up identities

Nearly all DWN interfaces requires one or more identities to be passed in (as authorDID and recipientDID).  You can create some did:key identities by doing this:

```golang

authorKeypair := client.NewKeypair()
authorIdentity := client.FromKeypair(authorKeypair)

recipientKeypair := client.NewKeypair()
recipientIdentity := client.FromKeypair(recipientKeypair)

```


### Save Data to DWN Node

To save data to a DWN Node, and get its associated Record ID back:

```golang

schemaUrl := "https://openreserve.io/schema/test-schema.json"
data := []byte("This is some plain text data I want to save")

recordId, err := dwnClient.SaveData(schemaUrl, data, "text/plain", &authorIdentity, &recipientIdentity)

```

### Update data already stored on DWN Node

```golang

schemaUrl := "https://openreserve.io/schema/test-schema.json"
data := []byte("This is some plain text data I want to save")

updatedRecordId, err := dwnClient.UpdateDate(schemaUrl, recordId, data, "text/plain", &recipientIdentity)

```

### Query for the data

```golang

schemaUrl := "https://openreserve.io/schema/test-schema.json"
recordId := "baeaabrykujwwizltmnzgs4dun5zeg2lembwxa4tpmnsxg43jn..."
storedData, dataFormat, err := dwnClient.GetData(schemaUrl, recordId, &recipientIdentity)

```