# SET OPTIONS - Weights, Thresholds, and Signers

## Summary
Set the master key’s weight, determine the operation threshold, and manage signers

## Description
Imagine you're an adventurer. You and two companions have battled foes across land and sea to collect three keys, hoping they will unlock a vault door guarding a vast treasure. After much tribulation, your squad stands at the door, bewildered, on the brink of despair — the door is sealed with not three locks but five. You cast your eyes down in defeat. Wait, what is that? You see the number “2” etched into the key in your right palm. Curious, you implore your allies to look at their own keys. Ah hah! One has a “2” and the other a “1”. Your brilliant mind starts churning, and you realize that two of the keys fit two of the locks. Shaking, you insert the keys into their locks and turn. At the last spin of the final key, the vault door creaks open, revealing the priceless treasure you've journeyed so far to acquire. Ring Pops! Thousands upon thousands of sparkling, crystalline, tantalizing Ring Pops! What a journey. What a puzzle.

Alright, let's circle it back. The keys in this riveting story work like signatures on the Stellar network. Signatures are needed to authorize transactions. You sign transactions with the master key, which is the private key that corresponds to the account's public key. Some transactions require additional or alternate signatures, which is called multisignature (or just multisig).

In this quest, we will set signature weights and operation thresholds, and we will add additional signers to the Quest Account using the setOptions operation. Then, we'll submit a successful multisignature transaction.

## In english
Basically this quest is all about controlling the power of an account on its assets, xlm balances, and other operations on the stellar network. Imagine having a bank account of company X, all the money transfers need to be signed by at least 3 members of the company. If the 3 members don't agree on the transaction then it will never happen. To do this you create a first transaction where the `bank account`(stellar account) states that it's adding more signers public keys and also that the thresholds(low,med,high) will require the combination of the 3 members' weights. You may set the `bank account`'s weight to something lower than the required sum of threshold to prevent the account from signing its transactions without consent of the other members, or set it higher or equal to that some to keep its power on its own transactions.
The following links might help in understanding the concept better if this didn't help

- [signatures & multi-sig](https://developers.stellar.org/docs/encyclopedia/signatures-multisig)
- [multi-sig on stellar wallets explained](https://publicnode.org/what-is-multisig-on-a-stellar-wallet)