package pkg

func ResponseError(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "error",
		"message": message,
	}
}

func ResponseSuccess(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
		"data":   data,
	}
}

func ResponseSuccessWithMeta(data interface{}, meta interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
		"data":   data,
		"meta":   meta,
	}
}

func ResponseSuccessWithPagination(data interface{}, meta interface{}, pagination interface{}) map[string]interface{} {

	return map[string]interface{}{
		"status":     "success",
		"data":       data,
		"meta":       meta,
		"pagination": pagination,
	}
}

func ResponseSuccessWithPaginationAndFilter(data interface{}, meta interface{}, pagination interface{}, filter interface{}) map[string]interface{} {

	return map[string]interface{}{
		"status":     "success",
		"data":       data,
		"meta":       meta,
		"pagination": pagination,
		"filter":     filter,
	}
}
