package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Define an interface for testing API with mock DB
type Store interface{
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}


// This struct provides the functionality to execute individual queries as well as transactions.
type SQLStore struct {
	*Queries // The Queries struct defined in db.go provides the functionality to execute individual queries. Composition over inheritance
	db       *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Add function to generate database transaction
// fn is callback function
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()

}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Account   `json:"from_account"`
	ToAccount   Account   `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		// fmt.Println(txName, "create entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// update account balances.
		// Just getting the account and updating using the crud functions is not the correct way. The transactions running in different goroutines may fetch the accoutn without the balance being updated, leading to discrepencies
		// Have to use isolation by adding 'FOR UPDATE' in sql query. But this leads to deadlock because of foreign key constraint.
		// Have to tell postgres that pk will not be changed. Use 'NO KEY'
		// fmt.Println(txName, "get account 1")
		// Now there still is a possibility of deadlock when two concurrent transactions are updating the balance of the same account.
		// This happens when TX1 updates balance of account 1 first, and TX2 updates balance of account 2 first.
		// Can fix this by executing queries in consistent order. In this case, account with a smaller id should always be updated first.
		if arg.FromAccountID < arg.ToAccountID{
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else{
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64,
) (account1 Account, account2 Account, err error){
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: account1ID,
		Amount: amount1,
	})

	if err != nil{
		return 
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: account2ID,
		Amount: amount2,
	})

	if err != nil{
		return 
	}

	return
}
