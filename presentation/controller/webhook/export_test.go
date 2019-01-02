package webhook

import (
	"github.com/duck8823/duci/application/service/executor"
)

type Handler = handler

func (h *Handler) SetExecutor(executor executor.Executor) (reset func()) {
	tmp := h.executor
	h.executor = executor
	return func() {
		h.executor = tmp
	}
}
