package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

const (
	accountPublicKey  string = "GBVY4IXZHVT7HRBSOWDVTQBNJMT7WM6UJV4LK6NXCNIPKT76MI3WV4NZ"
	accountSecretKey  string = "SBB5UOQXVS2GDXG4ACKKVXDR25IF4G5RROEQPZZSKIYVRZUAUCRSXSHX"
	networkPassphrase string = network.TestNetworkPassphrase
)

var (
	networkClient *horizonclient.Client = horizonclient.DefaultTestNetClient
	domainName    string                = "untitled-weigpj7k4sff.runkit.sh"
)

func getMainAccountDetails() (horizon.Account, error) {

	kp, _ := keypair.ParseFull(accountSecretKey)

	req := horizonclient.AccountRequest{
		AccountID: kp.Address(),
	}

	return networkClient.AccountDetail(req)
}

func getMainAccountKeypair() *keypair.Full {

	kp, _ := keypair.ParseFull(accountSecretKey)

	return kp
}

func createSignedTransaction(account horizon.Account, signingAccountKeypair *keypair.Full) (*txnbuild.Transaction, error) {

	setOptionOperation := txnbuild.SetOptions{
		HomeDomain: &domainName,
	}

	transaction, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&setOptionOperation,
		},
		BaseFee: txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(),
		},
	})

	if err != nil {
		return nil, err
	}

	return transaction.Sign(networkPassphrase, signingAccountKeypair)
}

func main() {

	account, err := getMainAccountDetails()

	if err != nil {
		log.Println("error getting account details", err)
		return
	}

	signer := getMainAccountKeypair()

	transaction, err := createSignedTransaction(account, signer)

	if err != nil {
		log.Println("error creaing signed transaction", err)
		return
	}

	_, err = networkClient.SubmitTransaction(transaction)

	if err != nil {
		log.Println("error submitting the transaction", err)
		return
	}

	fmt.Println("Successfully submitted the set options operation")
}
