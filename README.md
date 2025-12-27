# MailBreeze Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/MailBreeze/mailbreeze-go.svg)](https://pkg.go.dev/github.com/MailBreeze/mailbreeze-go)
[![CI](https://github.com/MailBreeze/mailbreeze-go/actions/workflows/ci.yml/badge.svg)](https://github.com/MailBreeze/mailbreeze-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/MailBreeze/mailbreeze-go)](https://goreportcard.com/report/github.com/MailBreeze/mailbreeze-go)

The official Go SDK for the MailBreeze email platform.

## Installation

```bash
go get github.com/MailBreeze/mailbreeze-go
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/MailBreeze/mailbreeze-go"
)

func main() {
    client := mailbreeze.NewClient("sk_live_xxx")

    email, err := client.Emails.Send(context.Background(), &mailbreeze.SendEmailParams{
        From:    "hello@yourdomain.com",
        To:      []string{"user@example.com"},
        Subject: "Welcome!",
        HTML:    "<h1>Welcome to MailBreeze!</h1>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Email sent: %s\n", email.ID)
}
```

## Configuration

```go
// Custom configuration
client := mailbreeze.NewClient("sk_live_xxx",
    mailbreeze.WithTimeout(60*time.Second),
    mailbreeze.WithMaxRetries(5),
    mailbreeze.WithBaseURL("https://custom.api.com"),
)
```

## Resources

### Emails

```go
// Send an email
email, err := client.Emails.Send(ctx, &mailbreeze.SendEmailParams{
    From:    "hello@yourdomain.com",
    To:      []string{"user@example.com"},
    Subject: "Hello",
    HTML:    "<p>Hello, world!</p>",
})

// Send with idempotency key
email, err := client.Emails.Send(ctx, params, mailbreeze.WithIdempotencyKey("unique-key"))

// List emails
emails, err := client.Emails.List(ctx, &mailbreeze.ListEmailsParams{
    Status: mailbreeze.EmailStatusDelivered,
    Page:   1,
    Limit:  20,
})

// Get email by ID
email, err := client.Emails.Get(ctx, "email_123")

// Get email stats
stats, err := client.Emails.Stats(ctx)
```

### Lists

```go
// Create a list
list, err := client.Lists.Create(ctx, &mailbreeze.CreateListParams{
    Name:        "Newsletter Subscribers",
    Description: "Main newsletter list",
})

// List all lists
lists, err := client.Lists.List(ctx, nil)

// Get list by ID
list, err := client.Lists.Get(ctx, "list_123")

// Update list
list, err := client.Lists.Update(ctx, "list_123", &mailbreeze.UpdateListParams{
    Name: "Updated Name",
})

// Delete list
err := client.Lists.Delete(ctx, "list_123")

// Get list stats
stats, err := client.Lists.Stats(ctx, "list_123")
```

### Contacts

Contacts are scoped to a list:

```go
contacts := client.Contacts("list_123")

// Create contact
contact, err := contacts.Create(ctx, &mailbreeze.CreateContactParams{
    Email:     "user@example.com",
    FirstName: "John",
    LastName:  "Doe",
})

// List contacts
result, err := contacts.List(ctx, &mailbreeze.ListContactsParams{
    Status: mailbreeze.ContactStatusActive,
    Page:   1,
})

// Get contact
contact, err := contacts.Get(ctx, "contact_123")

// Update contact
contact, err := contacts.Update(ctx, "contact_123", &mailbreeze.UpdateContactParams{
    FirstName: "Jane",
})

// Delete contact
err := contacts.Delete(ctx, "contact_123")

// Suppress contact
contact, err := contacts.Suppress(ctx, "contact_123")
```

### Email Verification

```go
// Verify single email
result, err := client.Verification.Verify(ctx, &mailbreeze.VerifyEmailParams{
    Email: "test@example.com",
})
if result.IsValid {
    fmt.Println("Email is valid!")
}

// Batch verification
batch, err := client.Verification.Batch(ctx, []string{
    "user1@example.com",
    "user2@example.com",
})

// Get batch status
batch, err := client.Verification.Get(ctx, "ver_123")

// Get verification stats
stats, err := client.Verification.Stats(ctx)
fmt.Printf("Total Valid: %d, Valid %%: %.1f\n", stats.TotalValid, stats.ValidPercentage)
```

### Attachments

```go
// Create upload URL
upload, err := client.Attachments.CreateUpload(ctx, &mailbreeze.CreateUploadParams{
    Filename:    "document.pdf",
    ContentType: "application/pdf",
    Size:        12345,
})

// Upload file to upload.UploadURL using your preferred HTTP client

// Confirm upload
attachment, err := client.Attachments.Confirm(ctx, upload.AttachmentID)

// Use attachment in email
email, err := client.Emails.Send(ctx, &mailbreeze.SendEmailParams{
    From:          "hello@yourdomain.com",
    To:            []string{"user@example.com"},
    Subject:       "Document attached",
    HTML:          "<p>Please see attached.</p>",
    AttachmentIDs: []string{attachment.ID},
})
```

## Error Handling

```go
email, err := client.Emails.Get(ctx, "nonexistent")
if err != nil {
    if mailbreeze.IsNotFoundError(err) {
        fmt.Println("Email not found")
    } else if mailbreeze.IsAuthenticationError(err) {
        fmt.Println("Invalid API key")
    } else if mailbreeze.IsRateLimitError(err) {
        retryAfter := mailbreeze.GetRetryAfter(err)
        fmt.Printf("Rate limited, retry after %d seconds\n", retryAfter)
    } else if mailbreeze.IsValidationError(err) {
        fmt.Println("Validation error:", err)
    } else if mailbreeze.IsServerError(err) {
        fmt.Println("Server error, please retry")
    }
}
```

## Automatic Retries

The SDK automatically retries on:
- 429 Too Many Requests (with Retry-After header support)
- 5xx Server Errors

Configure retries:
```go
client := mailbreeze.NewClient("sk_live_xxx",
    mailbreeze.WithMaxRetries(5), // Default is 3
)
```

## License

MIT
