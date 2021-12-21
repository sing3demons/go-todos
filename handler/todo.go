package handler

type TodoHandler interface{}

type todoHandler struct{}

func NewtodoHandler() TodoHandler {
	return &todoHandler{}
}

// func (h *todoHandler) Find() {}
