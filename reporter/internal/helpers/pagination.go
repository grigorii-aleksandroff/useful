package helpers

func Pagination(page, perPage int64) (int, int) {
	limit := int(perPage)
	if limit < 0 {
		limit = 0
	}

	offset := 0
	if page > 1 && limit > 0 {
		offset = int(page-1) * limit
	}

	return limit, offset
}
