package router

type VenueQuote struct { Venue, Symbol string; Price, Quantity, Fee int64 }
type Router struct{}
func New() *Router { return &Router{} }
func (r *Router) Best(quotes []VenueQuote) (VenueQuote,bool) { if len(quotes)==0{return VenueQuote{},false}; best:=quotes[0]; for _,q:=range quotes[1:] { if q.Price+q.Fee < best.Price+best.Fee {best=q} }; return best,true }
