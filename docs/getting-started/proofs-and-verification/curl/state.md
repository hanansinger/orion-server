---
id: state
title: Proof of Existence of a State/Data
---

### State proof query

As mentioned in the Orion description [here](../../../introduction) and detailed [here](../../../architecture-and-design/provenance-data), Orion maintains a separate persisted graph data structure for historical data transitions, so a user can execute queries on those historical changes to understand the lineage of each data item. For more explanations about **provenance queries** and the different views they provide on historical data, see [Provenance queries](../../queries/curl/provenance).

In addition to the provenance graph store, Orion uses an Ethereum-style Merkle-Patricia Trie to provide cryptographically verifiable proofs of all state transitions.
It provides only one single type of proof: that a specific key was associated with a specific value when a specific block was committed to the ledger.

For each block, a root of the Merkle-Particia Trie is stored inside the block header - `state_merkel_tree_root_hash`. And because of the tamper-proof nature of the ledger, the trie root is enough to prove the existence of a specific value at the time the block was committed. 


**Sign json serialized query**
```sh
bin/signer -data '{"user_id":"alice","block_number":5,"db_name":"db2","key":"key1"}' -privatekey=deployment/sample/crypto/alice/alice.key
```

**Signature**
```
MEUCIElKzVqsY/4Yf1gf+3PCU0Su7KF8scdcTkRjZwQjjc0QAiEA/MbcF3XUzdfZkLIPI0jEayrqRwuC4bLLsqtzT5ArObc=
```
**Submit query**
```sh
curl \
     -H "Content-Type: application/json" \
     -H "UserID: alice" \
     -H "Signature: MEUCIElKzVqsY/4Yf1gf+3PCU0Su7KF8scdcTkRjZwQjjc0QAiEA/MbcF3XUzdfZkLIPI0jEayrqRwuC4bLLsqtzT5ArObc=" \
     -X GET -G "http://127.0.0.1:6001/ledger/proof/data/db2/key1?block=5" | jq .
```

**Output**
```json
{
  "response": {
    "header": {
      "node_id": "bdb-node-1"
    },
    "path": [
      {
        "hashes": [
          "AgIKCw4KDAkAAw8PBAkDDAIGDQILCAMHCQ8MAgkGBAsHDwYKAwUKBgkJAA0JCAwICQYBCg4ICwIGDA0ODwI=",
          "y08re7h7i3RY2CTFDDiZB5HQUhTK0pJfB++VRIG204E="
        ]
      },
      {
        "hashes": [
          null,
          null,
          null,
          null,
          "DpE0Npt4BFM3vZsJzmG4PO1xlMYTpoA5U3nK10Vejsc=",
          "0U6UMnKLCuzIMOsosKa1mZjfqk/vU/Lfpxw2nBDfuFQ=",
          null,
          null,
          null,
          "ZHm1aGvKmHJFL5H9gEIGyo3sNK+m+0T5Y4W7RalD3rE=",
          null,
          null,
          null,
          null,
          null,
          null
        ]
      },
      {
        "hashes": [
          "U9PoJB+8o7P9M7EsbVfVxNd2sENGYlXqeviBYUxLvU0=",
          "oj5KCPc3GFzGJPHjfpnFagPpQ5la98LWreApQNc3ssw=",
          "ENQ9r2wqD202CibxAlru6z/xfdHouQeTQ3x+466a4ik=",
          "ni4p/R5zAdxjOje/0vTLhk0nKSDqe9SgLxmXc7Iw+BY=",
          null,
          null,
          "8IqfhFECCgEPd8uOYjNnfNULy4l3oOirqU4hhIQChDE=",
          null,
          null,
          null,
          null,
          "P/6M7PdrnB9aoxEx5MkVtTVCfrQeUoGoFPk/oyey+48=",
          null,
          null,
          "ZLxYv3o1wbGKjAYNlfU5wzvH1VqB6PO+TpbDgrSyhDo=",
          "hgbr1+op4LBUUXNn1paw14L0kuaOyOU5XgtvXbWebl8="
        ]
      }
    ]
  },
  "signature": "MEYCIQCxrMF5rZlsv/4CcICzcvdH/Xbn+C99Mqswuvdy3gLzegIhANla0V7MGrir9c/I5Q+dIzRDjMnS7GJgdcah0p8XgO1U"
}
```

The first element in the path is a ValueNode. For node types, see [here].
