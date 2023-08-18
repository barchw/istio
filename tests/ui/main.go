// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	_ "embed"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

//go:embed kubeconfig.yaml
var kubeconfig []byte

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)

	allocator, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocator,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var buf []byte
	// navigate to a page, wait for an element, click
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(`http://localhost:3001`),
		chromedp.Click(`.fd-button.fd-margin-end--tiny.fd-margin-begin--tiny`, chromedp.NodeVisible),
		chromedp.Nodes(`div.view-line span`, &nodes),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(ctx,
		chromedp.MouseClickNode(nodes[0]),
		input.InsertText(string(kubeconfig)),
		chromedp.Click("//button/span[text()='Next step']"),
		chromedp.Click(`//button[@class="fd-button fd-button--emphasized fd-button--compact"]/span[text()='Connect cluster']`),
	)
	if err != nil {
		log.Fatal(err)
	}

	/*err = chromedp.Run(ctx,
		chromedp.Click(`//span[text()='Namespaces']`, chromedp.NodeVisible),
		chromedp.Click(`//a[text()='kyma-system']`, chromedp.NodeVisible),
		chromedp.Click(`//span[text()='Kyma']`, chromedp.NodeVisible),
		chromedp.Click(`//a/span[text()='Istio']`, chromedp.NodeVisible),
		chromedp.Click(`//a[text()='default']`, chromedp.NodeVisible),
		chromedp.Click(`//button/span[text()='Edit']`, chromedp.NodeVisible),
	)*/

	err = chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:3001/cluster/k3d-kyma/namespaces/kyma-system/istios/default"),
		chromedp.Click(`//button/span[text()='Edit']`, chromedp.NodeVisible),
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
}
