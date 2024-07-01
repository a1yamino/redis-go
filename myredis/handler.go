package myredis

type CommandHandler struct {
	Handler func(*Conn, []Value) bool
	persist bool
}

func (h *CommandHandler) should_persist() bool {
	return h.persist
}

func (h *CommandHandler) call(conn *Conn, args []Value) bool {
	return h.Handler(conn, args)
}

var defaultHandlers = map[string]CommandHandler{
	"PING":    {Handler: pingHandler, persist: false},
	"SET":     {Handler: SetHandler, persist: true},
	"GET":     {Handler: GetHandler, persist: false},
	"DEL":     {Handler: DelHandler, persist: true},
	"EXISTS":  {Handler: ExistsHandler, persist: false},
	"HSET":    {Handler: HSetHandler, persist: true},
	"HGET":    {Handler: HGetHandler, persist: false},
	"HGETALL": {Handler: HGetAllHandler, persist: false},
	"HDEL":    {Handler: HDelHandler, persist: true},
	"HLEN":    {Handler: HLenHandler, persist: false},
	"HKEYS":   {Handler: HKeysHandler, persist: false},
	"HVALS":   {Handler: HValsHandler, persist: false},
}

func pingHandler(conn *Conn, args []Value) bool {
	resp := "PONG"
	if len(args) > 0 {
		resp = args[0].String()
	}

	conn.Writer.WriteSimpleString(resp)
	return true
}
