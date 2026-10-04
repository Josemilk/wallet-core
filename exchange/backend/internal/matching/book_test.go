package matching

import "testing"

func TestBookMatchesCrossingOrders(t *testing.T) {
    b := NewBook()
    b.Place(Order{ID:"sell-1",UserID:"maker",Symbol:"BTC-USDT",Side:"SELL",Price:100,Quantity:10,Time:1})
    fills := b.Place(Order{ID:"buy-1",UserID:"taker",Symbol:"BTC-USDT",Side:"BUY",Price:105,Quantity:4,Time:2})
    if len(fills) != 1 || fills[0].Quantity != 4 { t.Fatalf("unexpected fills: %#v", fills) }
}

func TestBookDoesNotCross(t *testing.T) {
    b := NewBook()
    b.Place(Order{ID:"sell-1",UserID:"maker",Symbol:"BTC-USDT",Side:"SELL",Price:110,Quantity:10})
    fills := b.Place(Order{ID:"buy-1",UserID:"taker",Symbol:"BTC-USDT",Side:"BUY",Price:105,Quantity:4})
    if len(fills) != 0 { t.Fatalf("expected no fill: %#v", fills) }
}
