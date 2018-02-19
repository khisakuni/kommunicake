package gmail

type messageProcessor struct {
	historyID uint64
	userID    int
	db        *database.DB
}

func processMessage(srv *gmail.Service, messageID string) error {
	messageRes, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return err
	}

	// LOGIC
	// IF messateRes.Payload.Body.Data is present
	//     Parse this and treat this as message
	// ELSE
	//     Use messageRes.Payload.Parts as message
	//     (Prefer mime HTML)
	fmt.Printf(">>> DATA >>> %s\n", messageRes.Payload.Body.Data)
	fmt.Printf("PARTS COUNT: %d\n", len(messageRes.Payload.Parts))
	for _, part := range messageRes.Payload.Parts {
		fmt.Printf(">>>> mime: %s\n", part.MimeType)
		decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return err
		}
		fmt.Printf("message: %s\n", string(decoded))

	}

	return nil
}
