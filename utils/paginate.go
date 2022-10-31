package utils

// TotalPage returns the total number of pages.
func TotalPage(total int64, pageSize int64) (totalPage int64) {
	if total%pageSize == 0 {
		totalPage = total / pageSize
	} else {
		totalPage = total/pageSize + 1
	}

	if totalPage == 0 {
		totalPage = 1
	}

	return
}

// CurrentPage returns the current page.
func CurrentPage(page int64, totalPages int64) int64 {
	if page <= 0 || totalPages <= 0 {
		return 1
	}
	if page > totalPages {
		return totalPages
	}

	return page
}
