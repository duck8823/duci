package job

import "github.com/duck8823/duci/application/service/job"

type Handler = handler

func (h *Handler) SetService(service job_service.Service) (reset func()) {
	tmp := h.service
	h.service = service
	return func() {
		h.service = tmp
	}
}
