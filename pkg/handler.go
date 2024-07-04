package pkg

type CommandHandler struct {
	Handler func(IWriter, []Value) bool
	persist bool
}

func (h *CommandHandler) should_persist() bool {
	return h.persist
}

func (h *CommandHandler) call(w IWriter, args []Value) bool {
	return h.Handler(w, args)
}

var defaultHandlers = map[string]CommandHandler{
	// Connection
	"PING": {Handler: pingHandler, persist: false},

	// String
	"SET":    {Handler: SetHandler, persist: true},
	"GET":    {Handler: GetHandler, persist: false},
	"DEL":    {Handler: DelHandler, persist: true},
	"EXISTS": {Handler: ExistsHandler, persist: false},

	// Hash
	"HSET":    {Handler: HSetHandler, persist: true},
	"HGET":    {Handler: HGetHandler, persist: false},
	"HGETALL": {Handler: HGetAllHandler, persist: false},
	"HDEL":    {Handler: HDelHandler, persist: true},
	"HLEN":    {Handler: HLenHandler, persist: false},
	"HKEYS":   {Handler: HKeysHandler, persist: false},
	"HVALS":   {Handler: HValsHandler, persist: false},

	// List
	"LPUSH":  {Handler: LPushHandler, persist: true},
	"RPUSH":  {Handler: RPushHandler, persist: true},
	"LPOP":   {Handler: LPopHandler, persist: true},
	"RPOP":   {Handler: RPopHandler, persist: true},
	"LRANGE": {Handler: LRangeHandler, persist: false},
	"LLEN":   {Handler: LLenHandler, persist: false},
	"LTRIM":  {Handler: LTrimHandler, persist: true},
}

func pingHandler(w IWriter, args []Value) bool {
	resp := "PONG"
	if len(args) > 0 {
		resp = args[0].String()
	}

	w.WriteSimpleString(resp)
	return true
}
