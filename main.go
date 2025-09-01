package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/Knetic/govaluate"
    "github.com/diamondburned/gotk4/pkg/gdk/v4"
    "github.com/diamondburned/gotk4/pkg/glib/v2"
    "github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
    app := gtk.NewApplication("com.chrisnharvey.live-calculator", 0)
    app.ConnectActivate(func() {
        win := gtk.NewApplicationWindow(app)
        win.SetTitle("Live Calculator")
        win.SetDefaultSize(400, 50)
        win.SetResizable(false)

        box := gtk.NewBox(gtk.OrientationVertical, 8)
        box.SetMarginTop(12)
        box.SetMarginBottom(12)
        box.SetMarginStart(12)
        box.SetMarginEnd(12)

        input := gtk.NewEntry()
        input.SetPlaceholderText("Type something to calculate (e.g. (4+5)*8 or 2×2)")
        box.Append(input)

        // Create a clickable result boxp
        resultButton := gtk.NewButton()
        resultButton.SetLabel("")
        resultButton.SetHAlign(gtk.AlignFill)
        resultButton.AddCSSClass("result-box")
        resultButton.SetVisible(false) // Hide initially
        
        // Add some styling to make it look like a result box
        cssProvider := gtk.NewCSSProvider()
        cssProvider.LoadFromData(`
            .result-box {
                background-color: #f0f0f0;
                border: 1px solid #ccc;
                border-radius: 4px;
                padding: 8px;
                margin-top: 4px;
                font-family: monospace;
                font-size: 14px;
            }
            .result-box:hover {
                background-color: #e0e0e0;
                cursor: pointer;
            }
        `)
        
        display := gdk.DisplayGetDefault()
        gtk.StyleContextAddProviderForDisplay(display, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
        
        box.Append(resultButton)

        win.SetChild(box)

        input.ConnectChanged(func() {
            txt := input.Text()
            if txt == "" {
                resultButton.SetVisible(false)
                return
            }
            expr, err := govaluate.NewEvaluableExpression(txt)
            if err != nil {
                resultButton.SetVisible(false)
                return
            }
            val, err := expr.Evaluate(nil)
            if err != nil {
                resultButton.SetVisible(false)
                return
            }
            resultText := fmt.Sprintf("= %v", val)
            resultButton.SetLabel(resultText)
            resultButton.SetVisible(true)
        })

        // Handle click to copy to clipboard
        resultButton.ConnectClicked(func() {
            text := resultButton.Label()
            if text != "" {
                clipboard := gdk.DisplayGetDefault().Clipboard()
                clipboard.SetText(strings.TrimPrefix(text, "= "))
                
                // Provide visual feedback
                originalText := resultButton.Label()
                resultButton.SetLabel("Copied!")
                
                // Reset text after a short delay
                glib.TimeoutSecondsAdd(1, func() bool {
                    resultButton.SetLabel(originalText)
                    return false // Don't repeat
                })
            }
        })

        // Close on escape
        keyController := gtk.NewEventControllerKey()
        keyController.ConnectKeyPressed(func(keyval uint, keycode uint, state gdk.ModifierType) bool {
            if keyval == gdk.KEY_Escape {
                app.Quit()
                return true // handled
            }
            return false
        })

        win.AddController(keyController)

        // Close on focus loss
        focusController := gtk.NewEventControllerFocus()
        focusController.ConnectLeave(func() {
            app.Quit()
        })
        win.AddController(focusController)

        win.Show()
    })

    os.Exit(app.Run(nil))
}