package gotoimage

import (
	"context"
	"fmt"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func Render(sourcefile string) ([]byte,error){

	//it will return if source file is not exist
	if ! Existfile(sourcefile){

		return nil,errors.New("source file is not exist")
	}


	// Disable chrome headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx1, cancel1 := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel1()

	// create a timeout ctx for chrome ctx
	ctx, cancel2 := context.WithTimeout(ctx1, time.Minute)
	defer cancel2()

	//imgbuf is used for store the screenshot of "https://play.golang.org/" execution page
	var imgbuf []byte

	selcode:=`//*[@id="code"]`
	golangplayurl:=`https://play.golang.org/`

	err := chromedp.Run(ctx,

		chromedp.Navigate(golangplayurl),
		chromedp.WaitVisible("div.linedtextarea"),
		//Clear the default code area
		chromedp.Clear(selcode,chromedp.BySearch),
		chromedp.SendKeys(selcode,generategocode(sourcefile),chromedp.BySearch),
		chromedp.Click("fmt",chromedp.ByID),
		chromedp.Click("run",chromedp.ByID),
		chromedp.WaitVisible("span.system"),

		screenshot("body",&imgbuf,chromedp.BySearch),
		//wait 2 seconds
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		log.Println(err)
		return nil,err
	}

	if err:=Saveimage(filepath.Base(sourcefile)+".png",imgbuf);err!=nil{
		log.Println(err)
	}

	return imgbuf,err

}

//read source code from file
func generategocode(filename string) string  {


	buffile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	return string(buffile)

}

func Existfile(filename string) bool {

	_,err:=os.Stat(filename)
	if err!=nil{

		return false
	}
	return true

}

//save the result image as png
func Saveimage(name string,buf []byte) error  {
	if buf==nil{
		panic("buf cannot be nil")
	}
	err:=ioutil.WriteFile(name,buf,644)
	if err!=nil{
		log.Println("write file err ",err)
	}
	return err
}

//screenshot for play.golang.org executing result
func screenshot(sel interface{}, picbuf *[]byte, opts ...chromedp.QueryOption) chromedp.QueryAction {
	if picbuf == nil {
		panic("picbuf cannot be nil")
	}

	return chromedp.QueryAfter(sel, func(ctx context.Context,id runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		// get layout metrics
		_, _, contentSize, _,_,_,err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// get box model
		box, err := dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
		if err != nil {
			return err
		}
		if len(box.Margin) != 8 {
			return chromedp.ErrInvalidBoxModel
		}

		// take screenshot of the box
		buf, err := page.CaptureScreenshot().
			WithFormat(page.CaptureScreenshotFormatPng).
			WithClip(&page.Viewport{
				X:      math.Round(box.Margin[0]),
				Y:      math.Round(box.Margin[1]),
				Width:  math.Round(box.Margin[4] - box.Margin[0]),
				Height: math.Round(box.Margin[5] - box.Margin[1]),
				Scale:  1.0,
			}).Do(ctx)
		if err != nil {
			return err
		}

		*picbuf = buf
		return nil
	}, append(opts, chromedp.NodeVisible)...)
}
