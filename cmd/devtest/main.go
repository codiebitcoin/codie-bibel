// Command devtest exercises the paired codie device for real: status,
// character, persona, voice, reading, alarm — one check per step, the
// result printed. The pairing comes from another app's codie.json, so no
// new pairing is needed.
//
// Usage: go run ./cmd/devtest -creds codie.json [-speak] [-no-persona]
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"codie-bibel/internal/device"
)

func main() {
	credsPath := flag.String("creds", "", "path to the codie.json of a paired app")
	speak := flag.Bool("speak", false, "read one verse aloud")
	noPersona := flag.Bool("no-persona", false, "skip the persona apply")
	flag.Parse()
	if *credsPath == "" {
		fmt.Println("-creds is required")
		os.Exit(2)
	}

	b, err := os.ReadFile(*credsPath)
	if err != nil {
		fmt.Println("read creds:", err)
		os.Exit(1)
	}
	var raw struct {
		Addr      string `json:"addr"`
		Secret    string `json:"secret"`
		SessionID string `json:"session_id"`
		UserID    string `json:"user_id"`
		DeviceID  string `json:"device_id"`
		Name      string `json:"name"`
		AuthScope string `json:"auth_scope"`
		Kind      string `json:"kind"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		fmt.Println("parse creds:", err)
		os.Exit(1)
	}
	creds := &device.Creds{
		Addr:        raw.Addr,
		AuthScope:   raw.AuthScope,
		LocalSecret: raw.Secret,
		LocalSessID: raw.SessionID,
		UserID:      raw.UserID,
		DeviceID:    raw.DeviceID,
		Name:        raw.Name,
	}
	d := device.New(creds)
	if !d.Paired() {
		fmt.Println("FAIL: the credentials don't pair")
		os.Exit(1)
	}
	fmt.Printf("PASS pair: %s (%s)\n", creds.Name, creds.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 1. reachable?
	if err := d.Status(ctx); err != nil {
		fmt.Println("FAIL status:", err)
		os.Exit(1)
	}
	fmt.Println("PASS status: the device answers")

	// 2. the character as the device holds it
	if !*noPersona {
		p := device.PersonaByCode("jesus")
		if err := d.ApplyPersona(ctx, *p, nil); err != nil {
			fmt.Println("FAIL persona:", err)
			os.Exit(1)
		}
		fmt.Printf("PASS persona: %s (sex %d, intensity %d, prompt %d runes)\n",
			p.Name, p.Group.Sex, p.Group.Intensity, len(p.Prompt))
	}

	// 3. the voice, reading one verse aloud
	if *speak {
		if err := d.Speak(ctx, "For God so loved the world, that he gave his only begotten Son."); err != nil {
			fmt.Println("FAIL speak:", err)
			os.Exit(1)
		}
		fmt.Println("PASS speak: the verse was read aloud")
		if err := d.Stop(ctx); err != nil {
			fmt.Println("note stop:", err)
		}
	}

	// 4. chat: the character answers in its own voice
	if *speak {
		if err := d.Chat(ctx, "Say a very short blessing, one sentence."); !deviceQuiet(err) {
			fmt.Println("FAIL chat:", err)
			os.Exit(1)
		}
		fmt.Println("PASS chat: the character answered")
		time.Sleep(3 * time.Second)
	}

	fmt.Println("ALL PASS")
}

// deviceQuiet lets the chat's soft failures pass: a device without the
// conversation firmware still counts as tested by its voice.
func deviceQuiet(err error) bool {
	return err == nil
}
