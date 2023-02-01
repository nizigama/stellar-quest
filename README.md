# CHANGE TRUST a k a create an asset

## Summary
Trust an asset

## Description
So far in our Stellar Quest Learn journey, we've just worked with Stellar's native token, XLM. In this quest, we'll learn about what needs to happen to start handling custom assets. Custom assets aren't created, they exist in the form of trust on other accounts. Like Santa, they only exist for those who believe (am I banned?). In Stellar, this belief is called a trustline.

Trustlines are an explicit opt-in for an account to hold a particular asset. To hold a specific asset, an account must establish a trustline with the issuing account using the change_trust operation.

In this quest, your challenge is to establish a trustline between the Quest Account and another account for a specific asset using the changeTrust operation.

## In english
Basically there is no such thing as `create an asset` in stellar. For an asset to exist on the stellar network, there needs to be at least one account that trust in it; Hence the trustline operation. An asset has a code and an issuer(keep in mind that the issuer can't hold any balance of that asset), the asset start to exist the moment another account trust in the existence of it by initiating a change trust operation on the asset. In other words, to `create an asset` an account must trust in that asset.