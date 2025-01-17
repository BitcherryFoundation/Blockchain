# crypto

crypto is the cryptographic package adapted for HashRs's uses

## Importing it
To get the interfaces,
`import "github.com/hashrs/blockchain/core/consensus/dpos-pbft/crypto"`

For any specific algorithm, use its specific module e.g.
`import "github.com/hashrs/blockchain/core/consensus/dpos-pbft/crypto/ed25519"`

If you want to decode bytes into one of the types, but don't care about the specific algorithm, use
`import "github.com/hashrs/blockchain/core/consensus/dpos-pbft/crypto/amino"`

## Binary encoding

For Binary encoding, please refer to the [HashRs encoding specification](https://github.com/hashrs/blockchain/core/consensus/dpos-pbft/blob/master/docs/spec/blockchain/encoding.md).

## JSON Encoding

crypto `.Bytes()` uses Amino:binary encoding, but Amino:JSON is also supported.

```go
Example Amino:JSON encodings:

ed25519.PrivKeyEd25519     - {"type":"hashrs/PrivKeyEd25519","value":"EVkqJO/jIXp3rkASXfh9YnyToYXRXhBr6g9cQVxPFnQBP/5povV4HTjvsy530kybxKHwEi85iU8YL0qQhSYVoQ=="}
ed25519.PubKeyEd25519      - {"type":"hashrs/PubKeyEd25519","value":"AT/+aaL1eB0477Mud9JMm8Sh8BIvOYlPGC9KkIUmFaE="}
sr25519.PrivKeySr25519   - {"type":"hashrs/PrivKeySr25519","value":"xtYVH8UCIqfrY8FIFc0QEpAEBShSG4NT0zlEOVSZ2w4="}
sr25519.PubKeySr25519    - {"type":"hashrs/PubKeySr25519","value":"8sKBLKQ/OoXMcAJVxBqz1U7TyxRFQ5cmliuHy4MrF0s="}
crypto.PrivKeySecp256k1   - {"type":"hashrs/PrivKeySecp256k1","value":"zx4Pnh67N+g2V+5vZbQzEyRerX9c4ccNZOVzM9RvJ0Y="}
crypto.PubKeySecp256k1    - {"type":"hashrs/PubKeySecp256k1","value":"A8lPKJXcNl5VHt1FK8a244K9EJuS4WX1hFBnwisi0IJx"}
```
