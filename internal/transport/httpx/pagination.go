package httpx

func PaginationMeta(requestID string, page, pageSize int, total int64) Meta {
	return Meta{
		RequestID: requestID,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
	}
}
