package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#include <stdlib.h>

// クリップボードからHTMLデータを取得する
char* getHTMLFromPasteboard() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSString *html = [pasteboard stringForType:@"public.html"];
    if (html == nil) {
        return NULL;
    }
    return strdup([html UTF8String]);
}

// クリップボードにプレーンテキストのみを設定する
void setTextToPasteboard(const char* plainText) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    [pasteboard clearContents];
    NSString *plainString = [NSString stringWithUTF8String:plainText];
    [pasteboard setString:plainString forType:NSPasteboardTypeString];
}

// クリップボードにカスタムデータとプレーンテキストを設定する
void setRichTextToPasteboard(const void* data, int length, const char* plainText) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    [pasteboard clearContents];

    // 互換性のためにプレーンテキストを設定
    NSString *plainString = [NSString stringWithUTF8String:plainText];
    [pasteboard setString:plainString forType:NSPasteboardTypeString];

    // カスタムバイナリデータを設定
    NSData *customData = [NSData dataWithBytes:data length:length];
    [pasteboard setData:customData forType:@"org.chromium.web-custom-data"];
}
*/
import "C"
import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"unsafe"

	"github.com/fetaro/docks_to_slack_go/src"
)

// readClipboardHTML はクリップボードからHTMLデータを読み込みます
func readClipboardHTML() (string, error) {
	clipboardHTML := C.getHTMLFromPasteboard()
	if clipboardHTML == nil {
		return "", errors.New("クリップボードにHTMLデータがありません")
	}
	defer C.free(unsafe.Pointer(clipboardHTML))
	return C.GoString(clipboardHTML), nil
}

// setClipboardText はクリップボードにプレーンテキストを設定します
func setClipboardText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.setTextToPasteboard(cText)
}

// setClipboardRichText はクリップボードにカスタムデータとプレーンテキストを設定します
func setClipboardRichText(data []byte, plainText string) {
	cPlainText := C.CString(plainText)
	defer C.free(unsafe.Pointer(cPlainText))

	cData := C.CBytes(data)
	defer C.free(unsafe.Pointer(cData))

	C.setRichTextToPasteboard(cData, C.int(len(data)), cPlainText)
}

func createChromiumData(plainText string, textyJSON map[string]interface{}) []byte {
	writer := src.NewPickleWriter()

	// Entry Count (uint32)
	writer.WriteUInt32(2)

	// --- Entry 1: public.utf8-plain-text ---
	writer.WriteString16("public.utf8-plain-text")
	writer.WriteString16(plainText)

	// --- Entry 2: slack/texty ---
	writer.WriteString16("slack/texty")
	// JSON marshal produces compact JSON by default
	jsonBytes, _ := json.Marshal(textyJSON)
	writer.WriteString16(string(jsonBytes))

	payload := writer.GetPayload()

	// Prepend the total size (uint32)
	finalData := make([]byte, 4+len(payload))
	binary.LittleEndian.PutUint32(finalData[0:4], uint32(len(payload)))
	copy(finalData[4:], payload)

	return finalData
}

func main() {
	debug := flag.Bool("debug", false, "デバッグモードを有効にする")
	debugShort := flag.Bool("d", false, "デバッグモードを有効にする（-debugのショートカット）")
	textOnly := flag.Bool("text", false, "プレーンテキスト形式のみをクリップボードにコピーします")
	textOnlyShort := flag.Bool("t", false, "プレーンテキスト形式のみをクリップボードにコピーします（-textのショートカット）")
	flag.Parse()

	isDebug := *debug || *debugShort
	isTextOnly := *textOnly || *textOnlyShort

	if isDebug {
		fmt.Println("クリップボードから読み込んでいます...")
	}

	htmlContent, err := readClipboardHTML()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isDebug {
		fmt.Println("----変換前(html)-----------------")
		fmt.Println(htmlContent)
	}

	generator := src.NewSlackListGenerator()
	result, err := generator.Generate(htmlContent)
	if err != nil {
		fmt.Println("Error generating result:", err)
		return
	}

	if isDebug {
		if isTextOnly {
			fmt.Println("----変換後(text)-----------------")
			fmt.Println(result.PlainText)
			fmt.Println("-----------------")
		} else {
			jsonBytes, _ := json.MarshalIndent(result.TextyJSON, "", "  ")
			fmt.Println("----変換後(slack/texty)-----------------")
			fmt.Println(string(jsonBytes))
			fmt.Println("-----------------")
		}
	}

	if isTextOnly {
		setClipboardText(result.PlainText)
	} else {
		binaryData := createChromiumData(result.PlainText, result.TextyJSON)
		setClipboardRichText(binaryData, result.PlainText)
	}
}
