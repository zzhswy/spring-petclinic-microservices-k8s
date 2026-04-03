package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
)

func executeBrowserJourney(targetURL string) error {
	// Create headless Chromium execution allocator directly from the Alpine binary
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("incognito", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	// Browser context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Prevent zombie runners by enforcing a hard 30s timeout per journey
	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	// Seed observability by simulating various web actions
	routes := []string{
		"/",
		"/vets.html",
		"/owners/find",
	}
	route := routes[rand.Intn(len(routes))]

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL+route),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
	)
	
	if err != nil {
		return fmt.Errorf("initial navigation failed for %s: %w", route, err)
	}

	// Dynamic interaction logic based on the user journey assigned
	if route == "/owners/find" {
		err = chromedp.Run(ctx,
			// Trigger a raw DOM element click on the "Find Owner" submit button
			chromedp.Click(`button[type="submit"]`, chromedp.NodeVisible),
			// Await for Kubernetes to pass the db query and render HTML back to the browser buffer
			chromedp.WaitVisible(`table`, chromedp.ByQuery),
		)
		if err != nil {
			return fmt.Errorf("owner search dom interaction failed: %w", err)
		}
	}

	return nil
}
