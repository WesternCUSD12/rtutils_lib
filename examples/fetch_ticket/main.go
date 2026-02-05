package main

import (
	"context"
	"fmt"
	"os"
	"rtutils_lib"
	"strings"
	"text/tabwriter"
)

func main() {
	client := rtutils_lib.NewClient(os.Getenv("RT_BASE_URL"), os.Getenv("RT_TOKEN"))
	ctx := context.Background()
	
	ticket, err := client.Tickets.Get(ctx, "249")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", ticket.ID)
	fmt.Fprintf(w, "Subject:\t%s\n", ticket.Subject)
	fmt.Fprintf(w, "Status:\t%s\n", ticket.Status)
	fmt.Fprintf(w, "Queue:\t%s\n", ticket.Queue)
	fmt.Fprintf(w, "Owner:\t%s\n", ticket.Owner)
	if len(ticket.Requestor) > 0 {
		fmt.Fprintf(w, "Requestor:\t%s\n", strings.Join(ticket.Requestor, ", "))
	}
	fmt.Fprintf(w, "Created:\t%s\n", ticket.Created)
	w.Flush()

	history, err := client.Tickets.GetHistory(ctx, "249")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching history: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stdout, "\n\nComments and History:\n")
	fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("=", 80))

	for i, tx := range history {
		fullTx, err := client.Tickets.GetTransaction(ctx, tx.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching transaction %s: %v\n", tx.ID, err)
			continue
		}
		if fullTx == nil {
			continue
		}

		attachmentIDs := fullTx.GetAttachmentIDs()
		if len(attachmentIDs) == 0 {
			continue
		}

		for _, attID := range attachmentIDs {
			att, err := client.Tickets.GetAttachment(ctx, attID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching attachment %s: %v\n", attID, err)
				continue
			}

			if att.Content == "" {
				continue
			}

			decoded, err := att.DecodedContent()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding attachment %s: %v\n", attID, err)
				continue
			}

			if att.ContentType != "text/plain" {
				continue
			}

			creator := fullTx.Creator
			if creator == "" {
				creator = "System"
			}
			fmt.Fprintf(os.Stdout, "\n[%d] %s at %s\n", i+1, creator, fullTx.Created)
			fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("-", 80))
			fmt.Fprintf(os.Stdout, "%s\n", decoded)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s\n", strings.Repeat("=", 80))
}
