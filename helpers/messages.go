package helpers

func Messagess(MessageCode int) string {
	var messages string
	switch MessageCode {
	case 1000:
		messages = "DATA_FOUND"
	case 1001:
		messages = "NOT_FOUND"
	case 1002:
		messages = "Success"
	case 1003:
		messages = "DATA_UPDATED"
	case 1004:
		messages = "DATA_DELETE"
	case 2000:
		messages = "ID not found"
	case 2001:
		messages = "Error In update"
	case 2002:
		messages = "Error in create"
	case 5000:
		messages = "LOGIN"
	case 5001:
		messages = "UNAUTHORIZED"
	case 5002:
		messages = "TOKAN_NOT_FOUND"
	default:
		messages = "okk"
	}
	return messages
}
