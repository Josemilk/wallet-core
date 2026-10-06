package ledger

import "testing"

func TestJournalMustBalancePerAsset(t *testing.T) {
    l := New()
    if err := l.Post(Journal{
        ID: "j1", IdempotencyKey: "k1",
        Entries: []Entry{
            {Account:"user", Asset:"BTC", Amount:-10},
            {Account:"custody", Asset:"BTC", Amount:10},
        },
    }); err != nil { t.Fatal(err) }

    if err := l.Post(Journal{
        ID: "j2", IdempotencyKey: "k2",
        Entries: []Entry{
            {Account:"user", Asset:"BTC", Amount:-10},
            {Account:"custody", Asset:"ETH", Amount:10},
        },
    }); err == nil { t.Fatal("cross-asset journal must not balance") }
}

func TestJournalIdempotency(t *testing.T) {
    l := New()
    j := Journal{ID:"j1",IdempotencyKey:"same",Entries:[]Entry{
        {Account:"user",Asset:"BTC",Amount:-1},
        {Account:"custody",Asset:"BTC",Amount:1},
    }}
    if err:=l.Post(j); err!=nil {t.Fatal(err)}
    if err:=l.Post(j); err!=nil {t.Fatal(err)}
}
