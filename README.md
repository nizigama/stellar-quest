# SET OPTIONS

## Summary
Set the home domain of an account

## Description
It's time to dip our toes into the wild and wonderful world of the setOptions operation. There will be several quests that involve this operation but here we'll just focus on the Home Domain field.

The Home Domain field points to the web domain where you host a stellar.toml file. Your stellar.toml file provides a common place where the Internet can find information about your Stellar integration. There's a lot of information you can put into the stellar.toml file, most common being [[CURRENCIES]] info for custom assets issued by your service. In order to be valid this file must be hosted at https://DOMAIN/.well-known/stellar.toml.

The file is a bridge between domain names and the official Stellar accounts utilized by the services available on that domain.

Anyone can look up your stellar.toml file, and it proves that you're the owner of the HTTPS domain linked to a Stellar account and that you claim responsibility for the accounts and assets listed on it.

In this quest, you will create, host, and link to a stellar.toml file from your Quest Keypair using the Set Options operation.